package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/users/", rt.wrap(rt.getAllUsers))
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.DELETE("/session", rt.wrap(rt.logOut))
	rt.router.GET("/users/:user-id/private-profile", rt.wrap(rt.getPrivate))
	rt.router.POST("/users/:user-id/private-profile", rt.wrap(rt.setPrivate))
	rt.router.DELETE("/users/:user-id/private-profile", rt.wrap(rt.setPublic))
	rt.router.GET("/users/:user-id", rt.wrap(rt.getUserProfile))
	rt.router.PATCH("/users/:user-id", rt.wrap(rt.updateUserProfile))
	rt.router.DELETE("/users/:user-id", rt.wrap(rt.deleteUser))
	rt.router.PUT("/users/:user-id/username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:user-id/profile-picture", rt.wrap(rt.getUserProfilePicture))
	rt.router.PUT("/users/:user-id/profile-picture", rt.wrap(rt.changeUserProfilePicture))
	rt.router.DELETE("/users/:user-id/profile-picture", rt.wrap(rt.deleteUserProfilePicture))
	rt.router.GET("/search-name-or-hashtag", rt.wrap(rt.search))
	rt.router.GET("/users/:user-id/follow", rt.wrap(rt.getUserFollows))
	rt.router.PUT("/users/:user-id/follow/:follow-id", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:user-id/follow/:follow-id", rt.wrap(rt.unfollowUser))
	rt.router.GET("/users/:user-id/ban", rt.wrap(rt.getBanUsers))
	rt.router.PUT("/users/:user-id/ban/:ban-user", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:user-id/ban", rt.wrap(rt.unbanUser))
	rt.router.GET("/photos/:photo-id", rt.wrap(rt.getSinglePhoto))
	rt.router.GET("/posts/:post-id", rt.wrap(rt.getPhoto))
	rt.router.GET("/posts/", rt.wrap(rt.getPhotos))
	rt.router.GET("/users/:user-id/posts/stream/", rt.wrap(rt.getMyStream))
	rt.router.GET("/users/:user-id/posts/", rt.wrap(rt.getUserPhotos))
	rt.router.POST("/users/:user-id/posts/", rt.wrap(rt.uploadPhoto))
	rt.router.PUT("/users/:user-id/posts/:post-id", rt.wrap(rt.modifyCaption))
	rt.router.DELETE("/users/:user-id/posts/:post-id", rt.wrap(rt.deletePhoto))
	rt.router.PUT("/users/:user-id/posts/:post-id/hashtags/:hashtag", rt.wrap(rt.addHashtag))
	rt.router.DELETE("/users/:user-id/posts/:post-id/hashtags/:hashtag", rt.wrap(rt.deleteHashtag))
	rt.router.GET("/posts/:post-id/hashtags/", rt.wrap(rt.getPostHashtags))
	rt.router.GET("/posts/:post-id/likes/", rt.wrap(rt.getLikes))
	rt.router.PUT("/posts/:post-id/likes/:user-id", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/posts/:post-id/likes/:user-id", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/posts/:post-id/comments/", rt.wrap(rt.getPhotoComments))
	rt.router.PUT("/posts/:post-id/comments/", rt.wrap(rt.commentPhoto))
	rt.router.GET("/comments/:comment-id", rt.wrap(rt.getComment))
	rt.router.DELETE("/comments/:comment-id", rt.wrap(rt.uncommentPhoto))

	
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}