package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"errors"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// authorization check
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	uData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, uData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized"){
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Authorization check failed: %s", err)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error performing authorization check: ", err)
		}
		return
	}

	// user2 existence check
	target_user := strings.TrimPrefix(ps.ByName("user2"), "user2=")
	exist, err := rt.db.UserInDB(target_user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking user's presence in DB: ", err)
		return
	}	
	if !exist {
		http.Error(w, "Searched user does not exist", http.StatusNotFound)
		return
	}
	if username == target_user {
		http.Error(w, "... Are you that lonely...?", http.StatusBadRequest)
		return
	}

	// get user1 and 2's ID 
	user1Data, err := rt.db.GetUserData(username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving user data: ", err)
		return
	}
	user2Data, err := rt.db.GetUserData(target_user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving user data: ", err)
		return
	}

	// next check that user2 hasn't banned user1
	banned, err := rt.db.HasBanned(user2Data.UserID, user1Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "You are blacklisted by searched user")
		return
	}

	// update DB
	err = rt.db.FollowUser(user1Data.UserID, user2Data.UserID)
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
	fmt.Fprint(w, "All good")
}
