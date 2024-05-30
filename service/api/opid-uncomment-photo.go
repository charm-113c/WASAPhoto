package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate requesting user
	user := r.Header.Get("requesting-user")
	user1Data, err := rt.db.GetUserData(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("uncommentPhoto() -> rt.db.GetUserData(user) -> Error getting user data: ", err)
		return
	}
	err = validateToken(r, user1Data.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("uncommentPhoto() -> validateToken() -> Error performing authorization check: ", err)
		}
		return
	}

	// get uploader's ID
	uploader := strings.TrimPrefix(ps.ByName("username"), "username=")
	uploaderData, err := rt.db.GetUserData(uploader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Searched user does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("uncommentPhoto() -> rt.db.GetUserData(uploader) -> Error getting user data: ", err)
		return
	}
	// get photoID
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("uncommentPhoto() -> strconv.Atoi(pID) -> Atoi conversion error: ", err)
		return
	}
	// same for commentID
	cID := strings.TrimPrefix(ps.ByName("commentID"), "commentID=")
	commentID, err := strconv.Atoi(cID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("uncommentPhoto() -> strconv.Atoi(cID) -> Atoi conversion error: ", err)
		return
	}

	// proceed to uncomment
	err = rt.db.UncommentPhoto(uploaderData.UserID, photoID, commentID, user1Data.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("uncommentPhoto() -> rt.db.UncommentPhoto -> Error uncommenting photo in DB: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
