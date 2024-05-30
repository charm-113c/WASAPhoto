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

// Liking a photo involves addding a row to Like table
// and incrementing the number of likes in Photo's db
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// authenticate liking user
	likeuser := strings.TrimPrefix(ps.ByName("liker"), "liker=")
	likeUserData, err := rt.db.GetUserData(likeuser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Username invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("likePhoto() -> rt.db.GetUserData(likeuser) -> Error getting user data: ", err)
		return
	}
	err = validateToken(r, likeUserData.UserID, rt.seckey)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "token signature is invalid") {
			http.Error(w, "Operation unauthorised, identifier missing or invalid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("likePhoto() -> validateToken() -> Error performing authorization check: ", err)
		}
		return
	}

	// then get uploader's data in order to like photo
	up := strings.TrimPrefix(ps.ByName("username"), "username=")
	uploaderData, err := rt.db.GetUserData(up)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("likePhoto() -> rt.db.GetUserData(up) -> Error getting user data: ", err)
		return
	}
	pID := strings.TrimPrefix(ps.ByName("photoID"), "photoID=")
	photoID, err := strconv.Atoi(pID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("likePhoto() -> strconv.Atoi() -> Atoi conversion error: ", err)
		return
	}

	// first, however, check that U2 hasn't banned U1
	// (all U1 needs to interact with the photo is the ID, which is public and easy to get)
	banned, err := rt.db.HasBanned(uploaderData.UserID, likeUserData.UserID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("likePhoto() -> rt.db.HasBanned() -> Error checking blacklist pair in DB: ", err)
		return
	}
	if banned {
		http.Error(w, "You are blacklisted by the photo uploader", http.StatusUnauthorized)
		return
	}

	err = rt.db.LikePhoto(uploaderData.UserID, photoID, likeUserData.UserID)
	if err != nil {
		// idempotency: if tuple already exists, do nothing
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("likePhoto() -> rt.db.LikePhoto() -> Error inserting likePhoto tuple: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
