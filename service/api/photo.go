package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// Get the profile picture of the given user. The profile picture does not contain any like, comment, date-time posted, etc; no metadata attached.
func (rt *_router) getUserProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("getUserProfilePicture: Token Error: \n\ttoken: %s, \n\terr: %w", token, err)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Errorf("getUserProfilePicture: Server Error: \n\ttoken: %s, \n\terr: %w", token, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {
		ctx.Logger.WithError(err).Errorf("getUserProfilePicture: Bad Request Error for the user-id format: \n\tuserid: %s, \n\terr: %w", userID, err)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Errorf("getUserProfilePicture: Server Error: \n\tuserid: %s, \n\terr: %w", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	photoID, err := rt.db.GetUserProfilePicture(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.Printf("getUserProfilePicture: Not found: \n\tUserID: %s, \n\tPhotoID: %s, \n\terr: %w", userID, photoID, err)
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getUserProfilePicture: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	if photoID == "" {
		ctx.Logger.Printf("getUserProfilePicture: \n\tThe user does not yet have a profile photo -> PhotoID: %s, \n\terr: %w", photoID, err)
		w.WriteHeader(http.StatusNotFound) // --
		_, _ = w.Write([]byte("The user does not yet have a profile photo"))
		return

	}

	// Retrive the photo and send it.
	file := filepath.Join(structs.PicFolder, photoID+".png")

	img, err := os.Open(file)

	// if errors.Is(err, os.ErrNotExist) {
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfilePicture: Error reading the file")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	defer img.Close()

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	_, _ = io.Copy(w, img)

}

func (rt *_router) getSinglePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("getSinglePhoto: Token Error in getSinglePhoto: %s", token)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Errorf("getSinglePhoto: Server Error in getSinglePhoto: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	photoID := ps.ByName("photo-id")

	// Check the validity of the photo-id
	photoID, err = structs.UuidCheck(photoID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("getSinglePhoto: Bad Request Error for the photo-id format in getSinglePhoto: %s", photoID)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Errorf("getSinglePhoto: Server Error in getSinglePhoto: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.GetSinglePhoto(photoID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorizedin getSinglePhoto")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrNotFound) {

		ctx.Logger.WithError(err).Error("Not found in getSinglePhoto")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Error on our part in getSinglePhoto")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// Retrieve the photo and send it.
	file := filepath.Join(structs.PicFolder, photoID+".png")

	img, err := os.Open(file)

	// if errors.Is(err, os.ErrNotExist) {
	if err != nil {
		ctx.Logger.WithError(err).Error("Error reading the file in getSinglePhoto")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	defer img.Close()

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, img)
}

func (rt *_router) changeUserProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	photoID, valCreated, err := rt.db.ChangeUserProfilePicture(userID, token, r)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Not found")
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
	_ = json.NewEncoder(w).Encode(photoID)
}

func (rt *_router) deleteUserProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = rt.db.DeleteUserProfilePicture(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
