package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
	"errors"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

// rt.router.POST("/photos/:uploader/:photoID/comments", rt.commentPhoto)
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate requester
	commenter := r.Header.Get("commenter-username")
	commenterData, err := rt.db.GetUserData(commenter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving user data: ", err)
		return
	}
	err = validateToken(r, commenterData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid"){
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Operation unauthorised, identifier missing or invalid")
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error performing authorization check: ", err)
		}
		return
	}

	// retrieve metadata from URI
	uploader := strings.TrimPrefix(ps.ByName("uploader"), "uploader=")
	uploaderData, err := rt.db.GetUserData(uploader)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving user data: ", err)
		return
	}
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi error conversion: ", err)
		return
	}
	// retrieve metadata from header
	uploadDate := r.Header.Get("upload-date")
	// get commentID == number of comments on photo + 1
	photoData, err := rt.db.GetPhotoData(uploaderData.UserID, photoID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving photo data: ", err)
		return
	}

	// first, however, check that U2 hasn't banned U1 
	// (all U1 needs to interact with the photo is the ID, which is public and easy to get)
	banned, err := rt.db.HasBanned(uploaderData.UserID, commenterData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		http.Error(w, "You are blacklisted by the photo uploader", http.StatusUnauthorized)
		return
	}

	// then, get the comment itself 
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-type invalid, want 'text/plain'", http.StatusBadRequest)
		return
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}
	comment := string(raw)
	// check comment length to see if it's not too long
	switch {
	case utf8.RuneCountInString(comment) > 500:
		http.Error(w, "Comment must not exceed 500 characters", http.StatusRequestEntityTooLarge)
		return
	case utf8.RuneCountInString(comment) == 0:
		http.Error(w, "Comment is empty", http.StatusBadRequest)
		return
	}

	// finally, upload to DB + update ncomments
	err = rt.db.UploadComment(comment, photoData.Comments, commenterData.UserID, photoID, uploaderData.UserID, uploadDate)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error inserting comment: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}