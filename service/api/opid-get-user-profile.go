package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"errors"
	"database/sql"

	"github.com/R-Andom13/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
)

/*
IMPORTANT:
To enforce authorization, Client will be required to send requesting user's username
in header field "username".
*/

func (rt *_router ) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// after authorization check, see if current user is in target user's blacklist
	// in which case return 'banned' message.
	user := strings.TrimPrefix(r.Header.Get("requesting-user"), "requesting-user ")
	uData, err := rt.db.GetUserData(user)
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
	
	// let's focus on target now
	target_user := strings.TrimPrefix(ps.ByName("username"), "username=")
	// check user2 existence
	ok, err := rt.db.UserInDB(target_user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking user existence: ", err)
		return
	}
	if !ok {
		http.Error(w, "Searched user does not exist", http.StatusBadRequest)
		return
	}
	user1Data, err := rt.db.GetUserData(user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	targetUserData, err := rt.db.GetUserData(target_user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}

	// check if target has banned searching user
	banned, err := rt.db.HasBanned(targetUserData.UserID, user1Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		fmt.Fprintf(w, "%s has blacklisted you", target_user)
		return
	}

	// else retrieve their data
	w.Header().Set("Content-Type", "application/json")
	// we already have general user data in targeUserData, so we need the photos
	targetPhotos, err := rt.db.GetPhotos(targetUserData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user photos: ", err)
		return
	}
	// sorting photos in reverse chrono order
	sort.Slice(targetPhotos, func(i, j int) bool {
		return targetPhotos[i].UploadDate.After(targetPhotos[j].UploadDate)
	})

	// then we need target's followers and following
	followers, err := rt.db.GetFollowers(targetUserData.UserID) 
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user's followers: ", err)
		return 
	}
	following, err := rt.db.GetFollowing(targetUserData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user's following: ", err)
		return
	}

	// now that we have all we need, we json and marshal
	targetProfile := &database.UserProfile {
			Username: target_user,
			Photos: targetPhotos,
			Nphotos: targetUserData.Nphotos,
			Followers: followers,
			Following: following,
		}
	fmt.Fprint(w, targetProfile)
	out, err := json.Marshal(targetProfile)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error marshalling user profile: ", err)
		return 
	}

	// and finally send out
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error outputing user profile: ", err)
		return 
	}
}
