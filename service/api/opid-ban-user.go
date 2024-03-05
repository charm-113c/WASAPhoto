package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

// /users/{username}/blacklist/{user2}
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// auth check + data retrieval
	user := strings.TrimPrefix(ps.ByName("username"), "username=")
	user2 := strings.TrimPrefix(ps.ByName("user2"), "user2=")
	user1Data, err := rt.db.GetUserData(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data")
		return
	} 
	user2Data, err := rt.db.GetUserData(user2)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data")
		return
	}
	err = validateToken(r, user1Data.UserID, rt.seckey)
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

	// check user2's existence
	ok, err := rt.db.UserInDB(user2)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking user existence: ", err)
		return
	}
	if !ok {
		http.Error(w, "Searched user does not exist", http.StatusBadRequest)
		return
	}

	// we revoke user2's following of user1 if it exists
	err = rt.db.UnfollowUser(user2Data.UserID, user1Data.UserID) // done first as op is idempotent and doesn't error on sqlErrNoRows
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error adding blacklist pair to DB: ", err)
		return
	}

	// check if user is already banned, in which case do nothing
	banned, err := rt.db.HasBanned(user1Data.UserID, user2Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User is already blacklisted")
		return
	}

	// else, add to blacklist
	err = rt.db.BanUser(user1Data.UserID, user2Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error adding blacklist pair to DB: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s blacklisted", user2)
}