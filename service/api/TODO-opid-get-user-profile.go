package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router ) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// returns searched user's profile to user
	// Check for authorization
	// Read searched user's name -> compile profile
	// 	- Get photos, nphotos, followers and following
	// 	- marshall into json, send as response
	// Pause. Let's deal with posting images first. Changing stuff here will be a pain if we don't do this first
}