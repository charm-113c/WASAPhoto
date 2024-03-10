package api

import (
	"fmt"
	"io"
	"net/http"
	"log"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/golang-jwt/jwt/v5"
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
		log.Println("Error reading body: ", err)
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
		log.Println("Error searching username in DB: ", err)
		return
	}
	if !inDB {
		// create user -> insert in DB
		uid, err := uuid.NewV7()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error creating userID: ", err)
			return 
		}
		// add user with nphotos=0
		err = rt.db.AddUser(username, uid.String()) 
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error adding user to DB: ", err)
			return
		}
	}
	// in any case, get userID to generate token
	uData, err := rt.db.GetUserData(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	// create JWT token for bearer auth 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userID":uData.UserID})
	signedToken, err := token.SignedString([]byte(rt.seckey))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error signing JWT token: ", err)
		return
	}
	// return token in response
	fmt.Fprintf(w, "Bearer authentication token: %s", signedToken)

}