package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// First slice with user followers' IDs and second slice with user followings' IDs
func (rt *_router) getUserFollows(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetUserFollows: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserFollows: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetUserFollows: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserFollows: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	FollowerIDs, FollowingIDs, err := rt.db.GetUserFollows(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("GetUserFollows: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("GetUserFollows: Not Found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserFollows: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	follows := structs.Follow{
		FollowerIDs:  FollowerIDs,
		FollowingIDs: FollowingIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(follows)

}

// The user with userID follows the user with followID
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	followID := ps.ByName("follow-id")

	// Check the validity of the userID and followID
	for idx, userid := range [...]string{userID, followID} {

		uid, err := structs.UuidCheck(userid)
		if idx == 0 {
			userID = uid
		} else {
			followID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.FollowUser(userID, followID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.ErrNotFound) {

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

// The user with userID unfollows the user with followID
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	followID := ps.ByName("follow-id")

	// Check the validity of the userID and followID
	for idx, userid := range [...]string{userID, followID} {

		uid, err := structs.UuidCheck(userid)
		if idx == 0 {
			userID = uid
		} else {
			followID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.UnfollowUser(userID, followID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if errors.Is(err, structs.ErrNotFound) {

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
