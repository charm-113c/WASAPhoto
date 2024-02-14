package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Update username while checking for auth and username validity
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-type", "text/plain")

	// check authorization
	currname := ps.ByName("username")
	currname = strings.TrimPrefix(currname, "username=")
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	fmt.Fprintln(w, "Current username: ", currname)

	valid, err := rt.db.TokenIsValid(currname, token)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking token validity: ", err)
		return 
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Operation unauthorised, header Authorization missing or invalid")
		return 
	}

	// check new name validity
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
	newname := string(body)
	if !usernameIsValid(newname) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid username: name must be alphanumeric string of 3-30 characters")
		return
	}
	
	// check if new name already present in DB
	inDB, err := rt.db.UserInDB(newname)
	if err != nil{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		fmt.Println("Error searching username in DB: ", err)
		return
	}
	if inDB {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Chosen username (%s) already exists, please choose a different username", newname)
		return
	}

	// else update username
	err = rt.db.SetNewName(currname, newname)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error updating username in DB: ", err)
		return
	}
	fmt.Fprint(w, "New name set, your new username is: ", newname)
}
