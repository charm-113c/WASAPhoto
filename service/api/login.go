package api

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin (w http.ResponseWriter, r *http.Request, ps httprouter.Params) (id int, e error) {
	// Logs in a user given a username. If username isn't in DB, we create a new user

	// read header to check that content-type is text/plain
	if r.Header.Get("Content-type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return -1, errors.New("content type not valid, want 'text/plain'")
	}
	
	body, err := io.ReadAll(r.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading body: ", err)
		return -1, err
	}
	username := string(body)

	// check if present in DB
	_, err = rt.db.FindByName(username)
	if err != nil {
		// if not, create user -> insert in DB
		err = rt.db.AddUser(username) // adds user with id=0 and nphotos=0
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error adding user to DB: ", err)
			return -1, err
		}
	}

	// return bearer auth identifier
	id = rand.Int() 
	

}