package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getBanUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Check the validity of the userID
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bannedUserIDs, err := rt.db.GetBanUsers(userID, token)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrForbidden) {
		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bannedUserIDs)
}

// The user with userID bans the user with banID
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	banID := ps.ByName("ban-user")

	// Check the validity of the userID and banID
	for idx, userid := range [...]string{userID, banID} {
		uid, err := structs.UuidCheck(userid)
		if idx == 0 {
			userID = uid
		} else {
			banID = uid
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

	err = rt.db.BanUser(userID, banID, token)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrForbidden) {
		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// The user with userID bans the user with banID
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	banID := ps.ByName("ban-user")

	// Check the validity of the userID and banID
	for idx, userid := range [...]string{userID, banID} {
		uid, err := structs.UuidCheck(userid)
		if idx == 0 {
			userID = uid
		} else {
			banID = uid
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

	err = rt.db.UnbanUser(userID, banID, token)
	if errors.Is(err, structs.ErrUnAuth) {
		ctx.Logger.WithError(err).Error("User Not Authorized")
		w.Header().Set("WWW-Authenticate", "Bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return
	} else if errors.Is(err, structs.ErrForbidden) {
		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
