package api

import (
	"encoding/json"
	"errors"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)



// return "204" if hashtag already exists, else "201".
func (rt *_router) addHashtag(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	userID 		:= ps.ByName("user-id")
	postID 		:= ps.ByName("post-id")
	hashtag 	:= ps.ByName("hashtag")

	// Check the validity of the userID and postID
	for _, id := range [...]string{userID, postID} {
		err = structs.UuidCheck(id)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("Bad Request format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	// Check the validity of the hashtag
	err = structs.PatternCheck(structs.HashtagPattern, hashtag, structs.HashtagMinLen, structs.HashtagMaxLen)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Bad Request Error format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	valCreated, err := rt.db.AddHashtag(userID, token, postID, hashtag)

	if errors.Is(err, structs.ErrUnAuth ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden Error")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	if valCreated == "204" {
		w.WriteHeader(http.StatusNoContent)
		return 

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(hashtag)
}

func (rt *_router) deleteHashtag(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	userID 		:= ps.ByName("user-id")
	postID 		:= ps.ByName("post-id")
	hashtag 	:= ps.ByName("hashtag")

	// Check the validity of the userID and postID
	for _, id := range [...]string{userID, postID} {
		err = structs.UuidCheck(id)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("Bad Request format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	// Check the validity of the hashtag
	err = structs.PatternCheck(structs.HashtagPattern, hashtag, structs.HashtagMinLen, structs.HashtagMaxLen)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Bad Request Error format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.DeleteHashtag(userID, token, postID, hashtag)

	if errors.Is(err, structs.ErrUnAuth ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden Error")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}


func (rt *_router) getPostHashtags(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	postID 	:= ps.ByName("post-id")

	// Check the validity of the postID
	err = structs.UuidCheck(postID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Bad Request format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	hashtags, err := rt.db.GetPostHashtags(token, postID)

	if errors.Is(err, structs.ErrUnAuth ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hashtags)
}