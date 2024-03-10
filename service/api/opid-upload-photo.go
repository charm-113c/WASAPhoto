package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"errors"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

// Allows a user to post a photo to their profile
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-type", "text/plain")
	// Check authorization
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	uData, err := rt.db.GetUserData(username)
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
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid"){
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Operation unauthorised, identifier missing or invalid")
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error performing authorization check: ", err)
		}
		return
	} 

	// get number of images uploaded
	n := r.Header.Get("n-images")
	if n == "" {
		http.Error(w, "Number of images not specified in header field 'n-images'", http.StatusBadRequest)
	}
	n_img, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error converting n-img to int: ", err)
		return
	}

	// read multipart/form-data request body
	err = r.ParseMultipartForm(64 << 20) // max size of 64 MB
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error parsing form data: ", err)
		return 
	}

	// Collect metadata
	upDate := r.FormValue("UploadDate")
	imgDesc := r.FormValue("Description")

	// we allow multiple images to be uploaded: process each
	// in case of errors:
	problemFiles := 0 
	for i := 1; i <= n_img; i++ {
		// note: failure to upload one image != abort entire operation
		// we instead try to upload the next images if they exist
		// and warn user of the failed attempts
		i_file:= strconv.Itoa(i)
		fieldname := "UploadedImage" + i_file  
		// open img file
		imgFile, handler, err := r.FormFile(fieldname)
		if err != nil {
			problemFiles++
			errMess := fmt.Sprintf("Something went wrong with file %s so it was ignored", handler.Filename)
			http.Error(w, errMess, http.StatusInternalServerError)
			log.Println("Error retrieving uploaded image file: ", err)
			imgFile.Close() // force close file
			continue
		}
		defer imgFile.Close()
		imgExtension := handler.Header.Get("Content-Type")
		if !strings.Contains(imgExtension, "image") {
			problemFiles++
			fmt.Fprintf(w, "File: %s is not an image and was ignored\n", handler.Filename)
			imgFile.Close()
			continue
		}
		// read file
		imgData, err := io.ReadAll(imgFile)
		if err != nil {
			problemFiles++
			errMess := fmt.Sprintf("Something went wrong with file %s so it was ignored", handler.Filename)
			http.Error(w, errMess, http.StatusInternalServerError)
			log.Println("Error reading image file: ", err)
			imgFile.Close()
			continue 
		} 
		// Get user data (with updated nphotos)
		uData, err = rt.db.GetUserData(username)
		if err != nil {
			problemFiles++
			errMess := fmt.Sprintf("Something went wrong with file %s so it was ignored", handler.Filename)
			http.Error(w, errMess, http.StatusInternalServerError)
			log.Println("Error finding user ID: ", err)
			imgFile.Close()
			continue
		}
		// upload with necessary data
		err = rt.db.UploadImage(uData.Nphotos, uData.UserID, imgData, imgDesc, upDate, imgExtension)
		if err != nil {
			problemFiles++
			errMess := fmt.Sprintf("Something went wrong with file %s so it was ignored", handler.Filename)
			http.Error(w, errMess, http.StatusInternalServerError)
			log.Println("Error inserting image into DB: ", err)
			imgFile.Close()
			continue
		}
	}
	if problemFiles > 0 {
		fmt.Fprint(w, "Some file(s) couldn't be uploaded, you may want to try uploading them again")
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "Image(s) uploaded succesfully")
}
