package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/photos/:photoID/likes/:username", rt.unlikePhoto)
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate liking user
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	likingUserData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, likingUserData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Operation unauthorised, identifier missing or invalid")
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
	uploader := r.Header.Get("uploader")
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

	w.WriteHeader(http.StatusNoContent)
}
