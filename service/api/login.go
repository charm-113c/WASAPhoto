package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/gofrs/uuid"
)

func (rt *_router) doLogin (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Logs in a user given a username. If username isn't in DB, we create a new user
	w.Header().Set("Content-type", "text-plain")
	// create nil value for ID
	// NilID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000") // Here we make use of the ostrich algorithm

	// read request header to check that content-type is text/plain
	if r.Header.Get("Content-type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Content type not valid, want 'text/plain'")
		return
	}
	
	username := ps.ByName("username")

	// check that username is valid
	if ! usernameIsValid(username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid username, must be alphanumeric of 3-30 characters")
		return
	}

	// check if present in DB
	_, err := rt.db.FindByName(username)
	if err != nil {
		// if not, create user -> insert in DB
		err = rt.db.AddUser(username) // adds user with id=0 and nphotos=0
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// create bearer auth identifier
	id, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	// update DB
	err = rt.db.SetID(id, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	// return id in response
	fmt.Fprintf(w, "User identifier: %s", id)
}