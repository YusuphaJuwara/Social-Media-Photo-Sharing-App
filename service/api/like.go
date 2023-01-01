package api

import (
	"encoding/json"
	"errors"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)


// Get the like-count and the user IDs who liked the post.
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
        return 

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
        return 

	}

	postID 		:= ps.ByName("post-id")

	// Check the validity of the postID
	err = structs.UuidCheck(postID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	like, err := rt.db.GetLikes( token, postID )

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(like)
}


func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
        return 

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
        return 

	}

	userID 		:= ps.ByName("user-id")
	postID 		:= ps.ByName("post-id")

	// Check the validity of the userID and postID
	for _, id := range [...]string{userID, postID} {

		err = structs.UuidCheck(id)
		if errors.Is(err, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.LikePhoto(userID, token, postID)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}


func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
        return 

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
        return 

	}

	userID 		:= ps.ByName("user-id")
	postID 		:= ps.ByName("post-id")

	// Check the validity of the userID and postID
	for _, id := range [...]string{userID, postID} {

		err = structs.UuidCheck(id)
		if errors.Is(err, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.UnlikePhoto(userID, token, postID)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
