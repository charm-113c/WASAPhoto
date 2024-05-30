package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authorization check
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	uData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("followUser() -> rt.db.GetUserData(username) -> Error getting user data: ", err)
		return
	}
	err = validateToken(r, uData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("followUser() -> validateToken() -> Error performing authorization check: ", err)
		}
		return
	}

	// user2 existence check
	target_user := strings.TrimPrefix(ps.ByName("user2"), "user2=")
	exist, err := rt.db.UserInDB(target_user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("followUser() -> rt.db.UserInDB() -> Error checking user's presence in DB: ", err)
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

	// get user2's ID
	user2Data, err := rt.db.GetUserData(target_user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("followUser() -> rt.db.GetUserData(target_user) -> Error retrieving user data: ", err)
		return
	}

	// next check that user2 hasn't banned user1
	banned, err := rt.db.HasBanned(user2Data.UserID, uData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("followUser() -> rt.db.HasBanned() -> Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("You are blacklisted by searched user"))
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("followUser() -> w.Write() -> Error checking blacklist pair in DB: ", err)
			return
		}
		return
	}

	// update DB
	err = rt.db.FollowUser(uData.UserID, user2Data.UserID)
	if err != nil {
		// if follow pair already exists, do nothing
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// log.Println("Hey there buddy")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "Someting went wrong", http.StatusInternalServerError)
		log.Println("followUser() -> rt.db.FollowUser() -> Error adding follow pair in DB: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
