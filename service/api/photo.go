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
	if errors.Is(err, structs.BadReqErr) {

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
	err = structs.UuidCheck(userID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	photoID, err := rt.db.GetUserProfilePicture(userID, token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// Retrive the photo and send it.
	file := filepath.Join("./pictures", photoID + ".png")

	img, err := os.Open(file)

	// if errors.Is(err, os.ErrNotExist) {
    if err != nil {
        ctx.Logger.WithError(err).Error("Error reading the file")
		w.WriteHeader(http.StatusInternalServerError)
		return

    }

    defer img.Close()

    w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

    io.Copy(w, img)
	
}



func (rt *_router) getSinglePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	photoID := ps.ByName("photo-id")

	// Check the validity of the photo-id
	err = structs.UuidCheck(photoID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.GetSinglePhoto(photoID, token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.NotFoundErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// Retrieve the photo and send it.
	file := filepath.Join("./pictures", photoID + ".png")

	img, err := os.Open(file)

	// if errors.Is(err, os.ErrNotExist) {
    if err != nil {
        ctx.Logger.WithError(err).Error("Error reading the file")
		w.WriteHeader(http.StatusInternalServerError)
		return

    }

    defer img.Close()

    w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
    io.Copy(w, img)
}



func (rt *_router) changeUserProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	err = structs.UuidCheck(userID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}


	// Form contains the parsed form data, including both the URL field's query parameters and 
	// the PATCH, POST, or PUT form data. This field is only available after ParseForm is called.
	// But ParseMultipartForm automatically calls ParseForm

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error parsing multipart form")
        w.WriteHeader(http.StatusBadRequest)
        return

	}

	// (multipart.File, *multipart.FileHeader, error)
	// FormFile returns the first file for the provided form key. 
	// FormFile calls ParseMultipartForm and ParseForm if necessary.
	// _ for getting the filenames, extensions, etc.
	photo_file, _, err := r.FormFile("photo")
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting file from form")
        w.WriteHeader(http.StatusBadRequest)
        return

	}

	photoID, valCreated, err := rt.db.ChangeUserProfilePicture(userID, token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	file := filepath.Join("./pictures", photoID + ".png")

	// Create creates or truncates the named file. If the file already exists, it is truncated. 
	// If the file does not exist, it is created with mode 0666 (before umask). 
	// If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. 
	// If there is an error, it will be of type *PathError.

	img, err := os.Create(file)

	if err != nil {
		ctx.Logger.WithError(err).Error("Error creating file")
        w.WriteHeader(http.StatusInternalServerError)
        return

	}

	defer img.Close()


	// Copy copies from src to dst until either EOF is reached on src or an error occurs. 
	// It returns the number of bytes copied and the first error encountered while copying, if any.

	// A successful Copy returns err == nil, not err == EOF. 
	// Because Copy is defined to read from src until EOF, it does not treat an EOF from Read as an error to be reported.

	_, err = io.Copy(img, photo_file)

	if err != nil {
		ctx.Logger.WithError(err).Error("Error copying file into img")
        w.WriteHeader(http.StatusInternalServerError)
        return

	}

	if valCreated == "204" {
		w.WriteHeader(http.StatusNoContent)
        return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photoID)
}



func (rt *_router) deleteUserProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	userID := ps.ByName("user-id")

	// Check the validity of the user-id
	err = structs.UuidCheck(userID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}


	photoID, err := rt.db.DeleteUserProfilePicture(userID, token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Not found")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	file := filepath.Join("./pictures", photoID + ".png")

	err = os.Remove(file)
    if err != nil {
		ctx.Logger.WithError(err).Error("Error removing file")
		w.WriteHeader(http.StatusInternalServerError)
        return

    }

	w.WriteHeader(http.StatusNoContent)
}


