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

// Liking a photo involves addding a row to Like table
// and incrementing the number of likes in Photo's db
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate liking user
	likeuser := strings.TrimPrefix(ps.ByName("username"), "username=")
	likeUserData, err := rt.db.GetUserData(likeuser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Username invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, likeUserData.UserID, rt.seckey)
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

	// then uploader's data in order to like photo
	up := strings.TrimPrefix(ps.ByName("uploader"), "uploader=")
	uploaderData, err := rt.db.GetUserData(up)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID) 
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi conversion error: ", err)
		return
	}

	err = rt.db.LikePhoto(uploaderData.UserID, photoID, likeUserData.UserID)
	if err != nil {
		// if tuple already exists, do nothing
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error inserting likePhoto tuple: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}