package api

import (
	"net/http"
)



// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)


	// Our stuff
	rt.router.POST("/login", rt.doLogin)
	rt.router.PUT("/users/:username/username", rt.setMyUserName)
	rt.router.GET("/users/:username/profile", rt.getUserProfile)
	rt.router.POST("/users/:username/profile/photos", rt.uploadPhoto)
	rt.router.PUT("/users/:username/profile/following/:user2", rt.followUser)
	return rt.router
}
