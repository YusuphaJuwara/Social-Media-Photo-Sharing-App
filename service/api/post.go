package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/api/reqcontext"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/julienschmidt/httprouter"
)

// This sends the post with all its metadata attached
// It should be named GetPost, but due to the project requirements, it is called GetPhoto.
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	postID := ps.ByName("post-id")

	// Check the validity of the postID
	postID, err = structs.UuidCheck(postID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Bad Request Error for the user-id format")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	post, err := rt.db.GetPhoto(postID, token)

	if errors.Is(err, structs.ErrUnAuth) {

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
	_ = json.NewEncoder(w).Encode(post)
}

// GetPhotos gets all posts of all users who did not set their profiles to private and did not ban the user.
func (rt *_router) getPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("getPhotos: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getPhotos: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	posts, err := rt.db.GetPhotos(token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("getPhotos: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getPhotos: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

// Get the list of posts posted by the given user's followings (including the user himself).
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("getMyStream: Token Error: %s", r.Header.Get("authorization"))
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getMyStream: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")

	// Check the validity of the userID
	userID, err = structs.UuidCheck(userID)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("getMyStream: Bad Request Error for the user-id format: %s", userID)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getMyStream: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	posts, err := rt.db.GetMyStream(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("getMyStream: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("getMyStream: Forbidden request")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("getMyStream: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

// Get the list of posts posted by the given user.
func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("Token Error: " + r.Header.Get("authorization"))
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

	posts, err := rt.db.GetUserPhotos(userID, token)

	if errors.Is(err, structs.ErrUnAuth) {

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
	_ = json.NewEncoder(w).Encode(posts)
}

// upload post
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	caption := r.FormValue("caption")

	if caption != "" {
		err = structs.PatternCheck(structs.MessagePattern, caption, structs.MessageMinLen, structs.MessageMaxLen)
		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Error("Bad Request Error format for the caption")
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	// r.Form returns values like map[string][]string.
	hashtags := r.Form["hashtags"]

	//if len(hashtags) > 0 {
	//ctx.Logger.Printf("All Hashtag: %s, Type: %t", hashtags, hashtags)
	for _, hashtag := range hashtags {
		if hashtag != "" {
			//ctx.Logger.Printf("one Hashtag: %s, Type: %t", hashtag, hashtag)
			err = structs.PatternCheck(structs.HashtagPattern, hashtag, structs.HashtagMinLen, structs.HashtagMaxLen)
			if errors.Is(err, structs.ErrBadReq) {

				ctx.Logger.WithError(err).Errorf("Bad Request Error format for the hashtag: %s", hashtag)
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

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Errorf("User Not Authorized\n postID: %s \nerr %v", postID, err)

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Errorf("Forbidden Error\n postID: %s \nerr %v", postID, err)
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Errorf("Error on our part\n postID: %s \nerr %v", postID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(postID)
}

func (rt *_router) modifyCaption(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("modifyCaption: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("modifyCaption: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	postID := ps.ByName("post-id")

	caption := r.FormValue("message")

	// Check the validities
	for idx, id := range [...]string{userID, postID} {
		uid, err := structs.UuidCheck(id)
		if idx == 0 {
			userID = uid
		} else {
			postID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Errorf("modifyCaption: Bad Request Error for UuidCheck \n\tpostID or userID: %s \n\t Error: %w", id, err)
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("modifyCaption: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = structs.PatternCheck(structs.MessagePattern, caption, structs.MessageMinLen, structs.MessageMaxLen)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Errorf("modifyCaption: Bad Request Error format for caption PatternCheck: \n\tCaption: %s \n\tError: $w", caption, err)
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("modifyCaption: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = rt.db.ModifyCaption(userID, token, postID, caption)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("modifyCaption: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("modifyCaption: forbidden error")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("modifyCaption: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Deletes a post with the given post ID together with the photo, caption, likes and comments, etc.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := structs.TokenCheck(r)
	if errors.Is(err, structs.ErrBadReq) {

		ctx.Logger.WithError(err).Error("deletePhoto: Token Error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else if err != nil {

		ctx.Logger.WithError(err).Error("deletePhoto: Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	userID := ps.ByName("user-id")
	postID := ps.ByName("post-id")

	// Check the validities
	for idx, id := range [...]string{userID, postID} {
		uid, err := structs.UuidCheck(id)
		if idx == 0 {
			userID = uid
		} else {
			postID = uid
		}

		if errors.Is(err, structs.ErrBadReq) {

			ctx.Logger.WithError(err).Errorf("deletePhoto: Bad Request Error for UuidCheck \n\tpostID or userID: %s \n\t Error: %w", id, err)
			w.WriteHeader(http.StatusBadRequest)
			return

		} else if err != nil {

			ctx.Logger.WithError(err).Error("deletePhoto: Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}

	err = rt.db.DeletePhoto(userID, token, postID)

	if errors.Is(err, structs.ErrUnAuth) {

		ctx.Logger.WithError(err).Error("deletePhoto: User Not Authorized")

		w.Header().Set("WWW-Authenticate", "Bearer ")
		// w.Header().Add("www-authenticate", "Bearer ")

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Must be authorized to access this website"))
		return

		// } else if errors.Is(err, structs.ErrNotFound) {

		// 	// ctx.Logger.WithError(err).Error("The post is not found. But return No Content")
		// 	w.WriteHeader(http.StatusNoContent)
		// 	return

	} else if errors.Is(err, structs.ErrForbidden) {

		ctx.Logger.WithError(err).Error("deletePhoto: Forbidden error")
		w.WriteHeader(http.StatusForbidden)
		return

	} else if err != nil {
		// if !errors.Is(err, structs.ErrNotFound)  {
		ctx.Logger.WithError(err).Error("deletePhoto: Error on our part")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if photoID == "" {
	// 	w.WriteHeader(http.StatusNotFound)		// --
	// 	_, _ = w.Write([]byte("No such photo"))
	//     return

	// }

	// file := filepath.Join(structs.PicFolder, photoID + ".png")

	// err = os.Remove(file)
	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error removing file")
	// 	w.WriteHeader(http.StatusInternalServerError)
	//     return

	// }

	w.WriteHeader(http.StatusNoContent)
}
