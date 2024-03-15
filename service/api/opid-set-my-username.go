package api

import (
	"database/sql"
	"errors"
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
	currname := strings.TrimPrefix(ps.ByName("username"), "username=")
	uData, err := rt.db.GetUserData(currname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, uData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error performing authorization check: ", err)
		}
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
		http.Error(w, "Invalid username: name must be alphanumeric string of 3-30 characters", http.StatusBadRequest)
		return
	}

	// check if new name already present in DB
	inDB, err := rt.db.UserInDB(newname)
	if err != nil {
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
	w.WriteHeader(http.StatusNoContent)
	// fmt.Fprintln(w, "Current username: ", currname)
	// fmt.Fprint(w, "New name set, your new username is: ", newname)
}
