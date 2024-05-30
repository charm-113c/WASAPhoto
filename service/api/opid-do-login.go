package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Logs in a user given a username. If username isn't in DB, we create a new user
	w.Header().Set("Content-Type", "text/plain")

	// read request header to check that content-type is text/plain
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type invalid, want 'text/plain'", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> io.ReadAll() -> Error reading body: ", err)
		return
	}
	username := string(body)

	// check that username is valid
	if !usernameIsValid(username) {
		http.Error(w, "Invalid username: name must be alphanumeric string of 3-30 characters", http.StatusBadRequest)
		return
	}

	// check if present in DB
	inDB, err := rt.db.UserInDB(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> rt.db.UserInDB() -> Error searching username in DB: ", err)
		return
	}
	if !inDB {
		// create user -> insert in DB
		uid, err := uuid.NewV7()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("doLogin() -> uuid.NewV7() -> Error creating userID: ", err)
			return
		}
		// add user with nphotos=0
		err = rt.db.AddUser(username, uid.String())
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("doLogin() -> rt.db.AddUser() -> Error adding user to DB: ", err)
			return
		}
	}
	// in any case, get userID to generate token
	uData, err := rt.db.GetUserData(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> rt.db.GetUserData() -> Error getting user data: ", err)
		return
	}
	// create JWT token for bearer auth
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userID": uData.UserID})
	signedToken, err := token.SignedString([]byte(rt.seckey))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> token.SignedString() -> Error signing JWT token: ", err)
		return
	}
	// marshall and output the token
	out, err := json.Marshal(signedToken)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> json.Marshall() -> Error marshalling token: ", err)
		return
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("doLogin() -> w.Write() -> Error sending token: ", err)
		return
	}
}
