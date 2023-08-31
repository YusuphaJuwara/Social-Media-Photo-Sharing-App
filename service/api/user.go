package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// For each user in our database, this gets it's information.
// But the user trying to get users' info cannot get those users who banned him or those whom he did not "follow" and set their profiles to private.
func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("Token Error: %v \n%s", err, token)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	users, err := rt.db.GetAllUsers(token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// Send the output to the user.
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

func (rt *_router) getPrivate(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("Bad Request Error for the user-id format: %s", userID)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	boolean, err := rt.db.GetPrivate(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authenticated")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("user-id not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(boolean)
}

func (rt *_router) setPrivate(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("Bad Request Error for the user-id format: %s", userID)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.SetPrivate(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authenticated")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) setPublic(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the user-id
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

	err = rt.db.SetPublic(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authenticated")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetUserProfile: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserProfile: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("GetUserProfile: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserProfile: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	user, err := rt.db.GetUserProfile(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("GetUserProfile: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("GetUserProfile: user-id not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("GetUserProfile: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) updateUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	profilename := r.FormValue("profile-name")
	profilemessage := r.FormValue("profile-message")
	gender := r.FormValue("gender")
	birthdate := r.FormValue("birth-date")

	ctx.Logger.Printf("UpdateUserProfile: \n\tprofilename: %s  \n\tprofilemessage: %s \n\tgender: %s \n\tbirthdate: %s", profilename, profilemessage, gender, birthdate)

	// If all four are empty, reject it
	if profilename == "" && profilemessage == "" && gender == "" && birthdate == "" {

		ctx.Logger.WithError(err).Errorf("UpdateUserProfile: At least one field must not be empty. \n\tprofilename: %s  \n\tprofilemessage: %s \n\tgender: %s \n\tbirthdate: %s", profilename, profilemessage, gender, birthdate)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// _ = json.NewDecoder(r.Body).Decode(userProf)
	userProf := structs.ProfileDetail{
		ProfileName:    profilename,
		ProfileMessage: profilemessage,
		Gender:         gender,
		BirthDate:      birthdate,
	}

	// Check the validities for the ones that are not empty
	if profilename != "" {
		err = structs.PatternCheck(structs.ProfileNamePattern, profilename, structs.ProfileNameMinLen, structs.ProfileNameMaxLen)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("UpdateUserProfile:: Bad Request Error for the user profile name format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	if profilemessage != "" {
		err = structs.PatternCheck(structs.MessagePattern, profilemessage, structs.MessageMinLen, structs.MessageMaxLen)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Bad Request Error for the user profile message name format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	if gender != "" {
		err = structs.GenderCheck(gender)
		if err != nil {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Bad Request Error for the user gender format")
			w.WriteHeader(http.StatusBadRequest)
			return

		}
	}

	if birthdate != "" && birthdate != "undefined" {
		err = structs.DateTimeCheck(structs.DatePattern, birthdate)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Bad Request Error for the user birthdate format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("UpdateUserProfile: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	valCreated, err := rt.db.UpdateUserProfile(userID, token, &userProf)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("UpdateUserProfile: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	if valCreated == "204" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Send this output to the user if any field was created.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(userProf)
}

func (rt *_router) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the user-id
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

	err = rt.db.DeleteUser(userID, token)

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

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("SetMyUserName: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("SetMyUserName: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("SetMyUserName: Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("SetMyUserName: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	username := r.FormValue("username")

	// Check the validity of the username
	err = structs.PatternCheck(structs.UsernamePattern, username, structs.UsernameMinLen, structs.UsernameMaxLen)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("SetMyUserName: Bad Request Error for the username format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("SetMyUserName: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.SetMyUserName(userID, token, username)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("SetMyUserName: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("SetMyUserName: Forbidden to modify another user's info")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("SetMyUserName: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
