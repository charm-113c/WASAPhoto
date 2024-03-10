package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// A playground to test stuff

func (rt *_router) test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// s, e := rt.db.GetFollowers("018dadaf-2d28-74ef-bb23-356317c4185a")
	// if e != nil {
	// 	http.Error(w, "GetFollowers test error", http.StatusInternalServerError)
	// 	log.Println("Error:", e)
	// 	return 
	// }
	// log.Println(s)

	// p, e := rt.db.GetPhotos("018dbd8f-d290-7edd-97df-442730797ddd")
	// if e != nil {
	// 	http.Error(w, "GetPhotos test error", http.StatusInternalServerError)
	// 	log.Println(e)
	// 	return
	// }
	// w.Header().Set("Content-type", "application/json")
	// fmt.Fprint(w, p)
	fmt.Fprint(w, rt.seckey)
	
}