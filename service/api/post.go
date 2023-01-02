package api

import (
	"encoding/json"
	"errors"
	_"io"
	"net/http"
	_"os"
	_"path/filepath"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// This sends the post with all its metadata attached
// It should be named GetPost, but due to the project requirements, it is called GetPhoto.
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	postID := ps.ByName("post-id")

	// Check the validity of the postID
	err = structs.UuidCheck(postID)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	post, err := rt.db.GetPhoto(postID, token)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}


// GetPhotos gets all posts of all users who did not set their profiles to private and did not ban the user.
func (rt *_router) getPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	posts, err := rt.db.GetPhotos( token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(posts)
}


// Get the list of posts posted by the given user's followings (including the user himself).
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the userID
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

	posts, err := rt.db.GetMyStream(userID, token)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Forbidden request")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}


// Get the list of posts posted by the given user.
func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Check the validity of the userID
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

	posts, err := rt.db.GetUserPhotos(userID, token)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}


// upload post
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	userID 		:= ps.ByName("user-id")

	// Check the validity of the userID
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

	// ParseMultipartForm parses a request body as multipart/form-data. 
	// The whole request body is parsed and up to a total of maxMemory bytes of its file parts are stored in memory,
	// with the remainder stored on disk in temporary files. ParseMultipartForm calls ParseForm if necessary. 
	// If ParseForm returns an error, ParseMultipartForm returns it but also continues parsing the request body. 
	// After one call to ParseMultipartForm, subsequent calls have no effect.
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
	// photo_file, _, err := r.FormFile("photo")
	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error getting file from form")
    //     w.WriteHeader(http.StatusBadRequest)
    //     return
	// }

	caption 	:= r.FormValue("caption")

	// r.Form returns values like map[string][]string.
	hashtags 	:= r.Form["hashtags"]

	if caption != "" {
		err = structs.PatternCheck(structs.MessagePattern, caption, structs.MessageMinLen, structs.MessageMaxLen)
		if errors.Is(err, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	if len(hashtags) > 0 {
		for _, hashtag := range hashtags {
			err = structs.PatternCheck(structs.HashtagPattern, hashtag, structs.HashtagMinLen, structs.HashtagMaxLen)
			if errors.Is(err, structs.BadReqErr) {

				ctx.Logger.WithError(err).Error("Bad Request Error format")
				w.WriteHeader(http.StatusBadRequest)
				return

			} else if err != nil {
				
				ctx.Logger.WithError(err).Error("Server Error")
				w.WriteHeader(http.StatusInternalServerError)
				return

			}
		}
	}

	postID, err := rt.db.UploadPhoto(userID, token, caption, hashtags, r)

	if errors.Is(err, structs.UnAuthErr ) {

		ctx.Logger.WithError(err).Error("User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ForbiddenErr) {

		ctx.Logger.WithError(err).Error("Forbidden Error")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// file := filepath.Join("./pictures", photoID + ".png")

	// Create creates or truncates the named file. If the file already exists, it is truncated. 
	// If the file does not exist, it is created with mode 0666 (before umask). 
	// If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. 
	// If there is an error, it will be of type *PathError.

	// img, err := os.Create(file)

	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error creating file")
    //     w.WriteHeader(http.StatusInternalServerError)
    //     return

	// }

	// defer img.Close()


	// Copy copies from src to dst until either EOF is reached on src or an error occurs. 
	// It returns the number of bytes copied and the first error encountered while copying, if any.

	// A successful Copy returns err == nil, not err == EOF. 
	// Because Copy is defined to read from src until EOF, it does not treat an EOF from Read as an error to be reported.

	// _, err = io.Copy(img, photo_file)

	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error copying file")
    //     w.WriteHeader(http.StatusInternalServerError)
    //     return

	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(postID)
}



func (rt *_router) modifyCaption(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
	postID := ps.ByName("post-id")

	caption := r.FormValue("message")

	// Check the validities
	for _, id := range [...]string{userID, postID} {
		err = structs.UuidCheck(id)
		if errors.Is(err, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = structs.PatternCheck(structs.MessagePattern, caption, structs.MessageMinLen, structs.MessageMaxLen)
	if errors.Is(err, structs.BadReqErr) {

		ctx.Logger.WithError(err).Error("Bad Request Error format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {
		
		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.ModifyCaption(userID, token, postID, caption)

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

	w.WriteHeader(http.StatusNoContent)
}


// Deletes a post with the given post ID together with the photo, caption, likes and comments, etc.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
	postID := ps.ByName("post-id")

	// Check the validities
	for _, id := range [...]string{userID, postID} {
		err = structs.UuidCheck(id)
		if errors.Is(err, structs.BadReqErr) {

			ctx.Logger.WithError(err).Error("Bad Request Error format")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {
			
			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.DeletePhoto(userID, token, postID)

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
	
	// file := filepath.Join("./pictures", photoID + ".png")

	// err = os.Remove(file)
    // if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error removing file")
	// 	w.WriteHeader(http.StatusInternalServerError)
    //     return

    // }

	w.WriteHeader(http.StatusNoContent)
}
