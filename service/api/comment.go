package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postID := ps.ByName("post-id")

	// Check the validity of the post-id
	postID, err = structs.UuidCheck(postID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comments, err := rt.db.GetPhotoComments(token, postID)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("GetPhotoComments: User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrNotFound) {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("GetPhotoComments: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

// Places a new comment and returns the newly created comment ID.
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Token Error \n\ttoken: %s \n\terr %v", token, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Server Error \n\ttoken: %s \n\terr %v", token, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postID := ps.ByName("post-id")
	message := r.FormValue("message")

	// Check the validities
	postID, err = structs.UuidCheck(postID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Bad Request Error for the post-id format: \n\tpostID: %s \n\terr %v", postID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Server Error: \n\tpostID: %s \n\terr %v", postID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = structs.PatternCheck(structs.MessagePattern, message, structs.MessageMinLen, structs.MessageMaxLen)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Bad Request Message format: \n\tmsgsg: %s \n\terr %v", message, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Server Error: \n\tmsg: %s \n\terr %v", message, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	commentID, err := rt.db.CommentPhoto(token, postID, message)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Errorf("commentPhoto: User Not Authorized: \n\ttoken: %s \n\terr %v", token, err)
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrNotFound) {
		ctx.Logger.WithError(err).Error("commentPhoto: Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Errorf("commentPhoto: Error on our part: \n\ttoken: %s \n\terr %v", token, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(commentID)

}

func (rt *_router) getComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	commentID := ps.ByName("comment-id")

	// Check the validity of the comment-id
	commentID, err = structs.UuidCheck(commentID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment, err := rt.db.GetComment(token, commentID)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
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
	_ = json.NewEncoder(w).Encode(comment)

}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	commentID := ps.ByName("comment-id")

	// Check the validity of the comment-id
	commentID, err = structs.UuidCheck(commentID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.UncommentPhoto(token, commentID)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("uncommentPhoto: User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrForbidden) {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Forbidden Error")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if errors.Is(err, structs.ErrNotFound) {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
