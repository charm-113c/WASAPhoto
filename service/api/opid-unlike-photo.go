package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"errors"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/photos/:uploader/:photoID/likes/:username", rt.unlikePhoto)
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate liking user
	username := strings.TrimPrefix(ps.ByName("username"),"username=")
	likingUserData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, likingUserData.UserID, rt.seckey)
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

	// identify the photo 
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi conversion error: ", err)
		return
	}
	uploader := strings.TrimPrefix(ps.ByName("uploader"), "uploader=")
	uploaderData, err := rt.db.GetUserData(uploader)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}

	// and remove Like triple from DB table
	err = rt.db.UnlikePhoto(uploaderData.UserID, photoID, likingUserData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error unliking photo: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}