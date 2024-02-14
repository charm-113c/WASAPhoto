package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/R-Andom13/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authorization check
	username := ps.ByName("username")
	token := r.Header.Get("Authorization")
	target_user := ps.ByName("user2")
	username = strings.TrimPrefix(username, "username=")
	token = strings.TrimPrefix(token, "Bearer ")
	target_user = strings.TrimPrefix(target_user, "user2=")
	valid, err := rt.db.TokenIsValid(username, token)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking token validity: ", err)
		return
	}
	if !valid {
		http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusBadRequest)
		return
	}
	// user2 existence check
	exist, err := rt.db.UserInDB(target_user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking user's presence in DB: ", err)
		return
	}	
	if !exist {
		http.Error(w, "Searched user does not exist", http.StatusNotFound)
		return
	}
	// get user1's ID 
	var userData database.UserData
	userData, err = rt.db.GetUserData(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving user data: ", err)
		return
	}
	// update DB
	err = rt.db.FollowingUser(userData.UserID, target_user)
	if err != nil {
		// if follow pair already exists, do nothing
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// log.Println("Hey there buddy")
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Someting went wrong", http.StatusInternalServerError)
		log.Println("Error adding follow pair in DB: ", err)
		return 
	}
	w.WriteHeader(http.StatusOK)
}