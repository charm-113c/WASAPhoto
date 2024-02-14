package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/R-Andom13/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
)

// Allows a user to post a photo to their profile
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-type", "text/plain")
	// Check authorization -> read request body -> obtain data/metadata -> add to db
	// Check authorization
	username := ps.ByName("username")
	username = strings.TrimPrefix(username, "username=")
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	valid, err := rt.db.TokenIsValid(username, token)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error checking token validity: ", err)
		return
	}
	if !valid {
		http.Error(w, "Operation unauthorised, header Authorization missing or invalid", http.StatusUnauthorized)
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
	// Collect/create metadata
	upDate := r.FormValue("UploadDate")
	imgDesc := r.FormValue("Description")
	var uData database.UserData

	// we allow multiple images uploaded: process each
	for i := 1; i <= n_img; i++ {
		i_file:= strconv.Itoa(i)
		fieldname := "UploadedImage" + i_file  
		// open img file
		imgFile, handler, err := r.FormFile(fieldname)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error retrieving uploaded image file: ", err)
			return
		}
		// read file
		imgData, err := io.ReadAll(imgFile)
		if err != nil {
			http.Error(w, "Error reading the image", http.StatusInternalServerError)
			log.Println("Error reading image file: ", err)
			return 
		} 
		imgExtension := handler.Header.Get("Content-Type")
		// Get user data (with updated nphotos)
		uData, err = rt.db.GetUserData(username)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Error finding user ID: ", err)
			return
		}
		// upload with necessary data
		err = rt.db.UploadImage(uData.Nphotos, uData.UserID, imgData, imgDesc, upDate, imgExtension)
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			log.Println("Error inserting image into DB: ", err)
			return
		}
		// finally, we increase user's nphotos
		rt.db.UpdateNphotos(username, true)
		imgFile.Close()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Image(s) uploaded succesfully")
}
