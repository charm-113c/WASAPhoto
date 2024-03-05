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

// rt.router.DELETE("/photos/:uploader/:photoID/comments/:commentID", rt.uncommentPhoto)
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate requesting user
	user := r.Header.Get("requesting-user")
	user1Data, err := rt.db.GetUserData(user)
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

	// get uploader's ID 
	uploader := strings.TrimPrefix(ps.ByName("uploader"), "uploader=")
	uploaderData, err := rt.db.GetUserData(uploader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	// get photoID
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi conversion error: ", err)
		return
	}
	// same for commentID
	cID := strings.TrimPrefix(ps.ByName("commentID"), "commentID=")
	commentID, err := strconv.Atoi(cID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi conversion error: ", err)
		return
	}

	// proceed to uncomment
	err = rt.db.UncommentPhoto(uploaderData.UserID, photoID, commentID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error uncommenting photo in DB: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}