package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// Get the like-count and the user IDs who liked the post.
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetLikes: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetLikes: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	postID := ps.ByName("post-id")

	// Check the validity of the postID
	postID, err = structs.UuidCheck(postID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetLikes: Bad Request Error for the post-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetLikes: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	like, err := rt.db.GetLikes(token, postID)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("GetLikes: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("GetLikes: Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetLikes: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(like)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("likePhoto: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("likePhoto: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	postID := ps.ByName("post-id")

	// Check the validity of the userID and postID
	for idx, id := range [...]string{userID, postID} {

		uid, err := structs.UuidCheck(id)
		if idx == 0 {
			userID = uid
		} else {
			postID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Errorf("likePhoto: Bad Request Error for the %s format", id)
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("likePhoto: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.LikePhoto(userID, token, postID)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("likePhoto: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("likePhoto: Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("likePhoto: Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("likePhoto: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("unlikePhoto: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("unlikePhoto: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	postID := ps.ByName("post-id")

	// Check the validity of the userID and postID
	for idx, id := range [...]string{userID, postID} {

		uid, err := structs.UuidCheck(id)
		if idx == 0 {
			userID = uid
		} else {
			postID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("unlikePhoto: Bad Request Error for the %s format", id)
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("unlikePhoto: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.UnlikePhoto(userID, token, postID)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("unlikePhoto: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("unlikePhoto: Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("unlikePhoto: Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("unlikePhoto: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
