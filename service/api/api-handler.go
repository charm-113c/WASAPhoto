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
	rt.router.POST("/users/:username/photos", rt.uploadPhoto)
	rt.router.PUT("/users/:username/following/:user2", rt.followUser)
	rt.router.GET("/users/:username/stream", rt.getMyStream)
	rt.router.PUT("/users/:username/blacklist/:user2", rt.banUser)
	rt.router.PUT("/users/:username/photos/:photoID/likes/:liker", rt.likePhoto)
	rt.router.POST("/users/:username/photos/:photoID/comments", rt.commentPhoto)
	rt.router.DELETE("/users/:username/photos/:photoID", rt.deletePhoto)
	rt.router.DELETE("/users/:username/following/:user2", rt.unfollowUser)
	rt.router.DELETE("/users/:username/blacklist/:user2", rt.unbanUser)
	rt.router.DELETE("/users/:username/photos/:photoID/likes/:liker", rt.unlikePhoto)
	rt.router.DELETE("/users/:username/photos/:photoID/comments/:commentID", rt.uncommentPhoto)
	rt.router.GET("/users/:username/photos/:photoID", rt.getPhoto)

	return rt.router
}
