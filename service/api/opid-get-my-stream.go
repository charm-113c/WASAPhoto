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

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/users/:username/Stream", rt.getMyStream)
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authorization check
	user := strings.TrimPrefix(ps.ByName("username"), "username=")
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

	// stream consists of images from followed users, organized in reverse chronological order
	// for efficiency and simplicity, create a single db op to get wanted images
	stream, err := rt.db.GetFollowedPhotos(uData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving followed users' photos: ", err)
		return
	}

	// stream is still in chronological order, so we reorganize it
	sort.Slice(stream, func(i, j int) bool {
		return stream[i].UploadDate.After(stream[j].UploadDate)
	})
	// then we prepare to output
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(stream)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error marshalling followed users' photos: ", err)
		return
	}
	// and we output
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error sending user's stream: ", err)
		return
	}
	
}