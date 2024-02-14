package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Logs in a user given a username. If username isn't in DB, we create a new user
	w.Header().Set("Content-type", "text/plain")

	// read request header to check that content-type is text/plain
	if r.Header.Get("Content-type") != "text/plain" {
		http.Error(w, "Content-type invalid, want 'text/plain'", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		fmt.Println("Error reading body: ", err)
		return
	}
	username := string(body)
	fmt.Fprintf(w, "Your username is: %s \n", username)

	// check that username is valid
	if !usernameIsValid(username) {
		http.Error(w, "Invalid username: name must be alphanumeric string of 3-30 characters", http.StatusBadRequest)
		return
	}

	// check if present in DB
	inDB, err := rt.db.UserInDB(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		fmt.Println("Error searching username in DB: ", err)
		return
	}
	if !inDB {
		// create user -> insert in DB
		uid, err := uuid.NewV7()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			fmt.Println("Error creating userID: ", err)
			return 
		}
		// add user with bearerAuthid='00000000-0000-0000-0000-000000000000' and nphotos=0
		err = rt.db.AddUser(username, uid.String()) 
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			fmt.Println("Error adding user to DB: ", err)
			return
		}
	}

	// create bearer auth identifier
	bearauthid, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		fmt.Println("Error creating bearerAuthID: ", err)
		return 
	}

	// update DB
	err = rt.db.SetID(bearauthid.String(), username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		fmt.Println("Error updating DB: ", err)
		return 
	}
	// return id in response
	fmt.Fprintf(w, "Bearer authentication token: %s", bearauthid)

	// ###########################
	// doLogin should also launch getMyStream at the end
	// ###########################
}