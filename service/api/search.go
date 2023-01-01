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
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Token Error")
		w.WriteHeader(http.StatusBadRequest)
        return 

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
        return 

	}

	name_or_hashtag := r.URL.Query().Get("name-hashtag")

	// Check the validity thereof
	err1 := structs.PatternCheck(structs.ProfileNamePattern, name_or_hashtag, structs.ProfileNameMinLen, structs.ProfileNameMaxLen)
	err2 := structs.PatternCheck(structs.HashtagPattern, name_or_hashtag, structs.HashtagMinLen, structs.HashtagMaxLen)

	if err1 != nil && err2 != nil {
		if errors.Is(err1, structs.BadReqErr) || errors.Is(err2, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err1 != nil || err2 != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	userIDs, postIDs, err := rt.db.Search(token, name_or_hashtag)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")	// user banned
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
	json.NewEncoder(w).Encode(IDs)

}
