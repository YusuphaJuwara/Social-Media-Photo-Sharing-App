/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"net/http"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// For each user in our database, this gets it's information.
	// But the user trying to get users' info cannot get those users who banned him or those
	// whom he did not "follow" and set their profiles to private.
	GetAllUsers(token string) ([]structs.User, error)

	// Logs in the user and returns the user ID, and the newly created token. In addition:
	// 1. If the user exists before, the string in the return value will be "200".
	// 2. If this is the first time (sign up), signs up the user and returns "201"
	DoLogin(username string) (string, string, string, error)

	// This deletes the sign-in token
	LogOut(token string) error

	// See if the user profile is set to private or not: true or false.
	GetPrivate(userID, token string) (bool, error)

	SetPrivate(userID, token string) error

	SetPublic(userID, token string) error

	GetUserProfile(userID, token string) (*structs.User, error)

	// Successfully created some fields that did not exist: 201, else 204.
	UpdateUserProfile(userID, token string, pd *structs.ProfileDetail) (string, error)

	DeleteUser(userID, token string) error

	SetMyUserName(userID, token string, username string) error

	// The returned uuid is the photo ID to be retrieved from directory.
	GetUserProfilePicture(userID, token string) (string, error)

	// The returned uuid is the photo ID to be changed in directory if exists, else nil.
	ChangeUserProfilePicture(userID, token string, r *http.Request) (string, string, error)

	DeleteUserProfilePicture(userID, token string) error

	// First slice with user IDs and second slice with post IDs that corresponds to the search term.
	Search(token string, search string) ([]string, []string, error)

	// First slice with user followers' IDs and second slice with user followings' IDs
	GetUserFollows(userID, token string) ([]string, []string, error)

	// The user with userID follows the user with followID
	FollowUser(userID, followID, token string) error

	// The user with userID unfollows the user with followID
	UnfollowUser(userID, followID, token string) error

	GetBanUsers(userID, token string) ([]string, error)

	// The user with userID bans the user with banID
	BanUser(userID, banID, token string) error

	// The user with userID unbans the user with banID
	UnbanUser(userID, banID, token string) error

	// Check if the photoID exists and the user is allowed to access it -> nil, else error.
	GetSinglePhoto(photoID, token string) error

	// This sends the post with all its metadata attached
	// It should be named GetPost, but due to the project requirements, it is called GetPhoto.
	GetPhoto(postID, token string) (*structs.Post, error)

	// GetPhotos gets all posts of all users who did not set their profiles to private and did not ban the user.
	GetPhotos(token string) ([]structs.Post, error)

	// Get the list of posts posted by the given user's followings.
	GetMyStream(userID, token string) ([]structs.Post, error)

	// Get the list of posts posted by the given user.
	GetUserPhotos(userID, token string) ([]structs.Post, error)

	UploadPhoto(userID, token string, caption string, hashtags []string, r *http.Request) (string, error)

	ModifyCaption(userID, token, postID, caption string) error

	// Deletes the post with the given post ID together with the photo, caption, likes and comments, etc.
	DeletePhoto(userID, token, postID string) error

	// Return "204" if hashtag already exists, else "201".
	AddHashtag(userID, token, postID string, hashtag string) (string, error)

	DeleteHashtag(userID, token, postID string, hashtag string) error

	// Return "201" if at least one of the hashtags didn't already exist, else "204".
	// AddHashtags( userID, token, postID string, hashtags []string ) ( string, error )

	// DeleteHashtags( userID, token, postID string, hashtags []string ) error

	GetPostHashtags(token, postID string) ([]string, error)

	// Get the like-count and the user IDs who liked the post.
	GetLikes(token, postID string) (*structs.Like, error)

	LikePhoto(userID, token, postID string) error

	UnlikePhoto(userID, token, postID string) error

	// Get all the comments of a given post.
	GetPhotoComments(token, postID string) ([]structs.Comment, error)

	//  Places a new comment and returns the newly created comment ID.
	CommentPhoto(token, postID string, message string) (string, error)

	GetComment(token, commentID string) (*structs.Comment, error)

	//  Deletes a given comment by the user.
	// The user is either the one who placed the comment or on whose post the comment was placed.
	UncommentPhoto(token, commentID string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Create the tables that do not exist in the database
	for i, sqlStmt := range structs.SqlStmtList {
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w\n AT: %d", err, i)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
