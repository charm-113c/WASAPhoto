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

// rt.router.DELETE("/users/:username/profile/photos/:photoID", rt.deletePhoto)
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// auth user
	username := strings.TrimPrefix(ps.ByName("username"), "username=")
	userData, err := rt.db.GetUserData(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			http.Error(w, "Provided username is invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error getting user data: ", err)
		return
	}
	err = validateToken(r, userData.UserID, rt.seckey)
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

	// get photoID and delete row from DB
	// note: all downstream operations (eliminating likes and comments) are done in the DBop
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Atoi conversion error: ", err)
		return
	}
	err = rt.db.DeletePhoto(userData.UserID, photoID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error deleting photo from DB: ", err)
		return
	} 

	w.WriteHeader(http.StatusNoContent)
}