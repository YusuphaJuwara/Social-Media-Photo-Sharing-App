package api

import (
	"encoding/json"
	"errors"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// First slice with user IDs and second slice with post IDs that correspond to the search term.
func (rt *_router) search(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	nameOrHashtag := r.URL.Query().Get("name-hashtag")

	// Check the validity thereof
	err1 := structs.PatternCheck(structs.ProfileNamePattern, nameOrHashtag, structs.ProfileNameMinLen, structs.ProfileNameMaxLen)
	err2 := structs.PatternCheck(structs.HashtagPattern, nameOrHashtag, structs.HashtagMinLen, structs.HashtagMaxLen)
	if err1 != nil && err2 != nil {
		if errors.Is(err1, structs.ErrBadReq) || errors.Is(err2, structs.ErrBadReq) {
			ctx.Logger.WithError(err).Error("Bad Request Error")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err1 != nil || err2 != nil {
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	userIDs, postIDs, err := rt.db.Search(token, nameOrHashtag)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrNotFound) {
		ctx.Logger.WithError(err).Error("Not found") // user banned
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	IDs := structs.IDs{
		UserIDs: userIDs,
		PostIDs: postIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(IDs)

}
