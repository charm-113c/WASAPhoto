package api

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Allows a user to post a photo to their profile
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check authorization
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	uData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, uData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
			log.Println("Error performing authorization check: ", err)
		}
		return
	}

	// read multipart/form-data request body
	err = r.ParseMultipartForm(32 << 20) // max size of 32 MB
	if err != nil {
		http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
		log.Println("Error parsing form data: ", err)
		return
	}

	// Collect metadata
	upDate := r.FormValue("UploadDate")
	imgDesc := r.FormValue("Description")

	// get image file
	imgFile, handler, err := r.FormFile("UploadedImage")
	if err != nil {
		http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
		log.Println("Error getting form file: ", err)
		return
	}
	defer imgFile.Close()
	imgExtension := handler.Header.Get("Content-Type")
	// note: frontend ensures we get an image type
	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
		log.Println("Error reading img file:", err)
		return
	}
	// we're now ready to upload (photoID = user's current TotNphotos -> also counts deleted photos)
	err = rt.db.UploadImage(uData.TotNphotos, uData.UserID, imgData, imgDesc, upDate, imgExtension)
	if err != nil {
		http.Error(w, "Something went wrong while uploading the photo", http.StatusInternalServerError)
		log.Println("Error uploading image to DB: ", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
