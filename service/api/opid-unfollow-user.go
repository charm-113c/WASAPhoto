package api

import (
	"log"
	"net/http"
	"strings"
	"fmt"
	"errors"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/users/:username/profile/following/:user2", rt.unfollowUser)
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate check
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	user1Data, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
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

	// get user2's data
	targetUser := strings.TrimPrefix(ps.ByName("user2"), "user2=")
	user2Data, err := rt.db.GetUserData(targetUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}

	// unfollow is a delete: idempotent operation
	err = rt.db.UnfollowUser(user1Data.UserID, user2Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error deleting follow pair: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}