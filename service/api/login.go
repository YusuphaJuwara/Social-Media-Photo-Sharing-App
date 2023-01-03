package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)



func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := r.FormValue("username")

	// Check the validity of the username
	err := structs.PatternCheck(structs.UsernamePattern, username, structs.UsernameMinLen, structs.UsernameMaxLen)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the username format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	
	ID, newtoken, valCreated, err := rt.db.DoLogin(username)

	if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Authorization", "Bearer " + newtoken)
	
	if valCreated == "201" {
		w.WriteHeader(http.StatusCreated)

	} else if valCreated == "200" {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(ID)
}

func (rt *_router) logOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = rt.db.LogOut(token)

	// if errors.Is(err, structs.ErrNotFound) {
	// 	ctx.Logger.WithError(err).Error("Token not found")
	// 	w.WriteHeader(http.StatusNotFound)
    //     return 
	// } else 
	if err != nil {
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
        return 
	}
	w.WriteHeader(http.StatusNoContent)
}

