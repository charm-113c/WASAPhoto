package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Our stuff
	rt.router.POST("/login", rt.doLogin)
	rt.router.PUT("/users/:username/username", rt.setMyUserName)
	rt.router.GET("/users/:username/profile", rt.getUserProfile)
	rt.router.POST("/users/:username/profile/photos", rt.uploadPhoto)
	rt.router.PUT("/users/:username/profile/following/:user2", rt.followUser)
	rt.router.GET("/users/:username/stream", rt.getMyStream)
	rt.router.PUT("/users/:username/blacklist/:user2", rt.banUser)
	rt.router.PUT("/photos/:photoID/likes/:username", rt.likePhoto)
	rt.router.POST("/photos/:photoID/comments", rt.commentPhoto)
	rt.router.DELETE("/users/:username/profile/photos/:photoID", rt.deletePhoto)
	rt.router.DELETE("/users/:username/profile/following/:user2", rt.unfollowUser)
	rt.router.DELETE("/users/:username/blacklist/:user2", rt.unbanUser)
	rt.router.DELETE("/photos/:photoID/likes/:username", rt.unlikePhoto)
	rt.router.DELETE("/photos/:photoID/comments/:commentID", rt.uncommentPhoto)

	return rt.router
}
