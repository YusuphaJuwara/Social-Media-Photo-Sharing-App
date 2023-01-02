package database

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"os"
	"time"
	"net/http"
	"io"

	"github.com/YusuphaJuwara/Social-Media-Photo-Sharing-App.git/service/structs"
	"github.com/gofrs/uuid"
)

// var (
// 	structs.UnAuthErr = structs.UnAuthError
// 	structs.ForbiddenErr = structs.ForbiddenError
// 	structs.NotFoundErr = structs.NotFoundError
// 	BadReqErr 	= structs.BadReqErr
// )

// If the server does not wish to make information available to the client (like accessing a user's details who banned the client), the status code 404 (Not Found) can be used instead.

// If user A sets his profile to private and user B does not follow him, then for user B, user A exists and can get only his profile info (including whether or not user A sets his profile to provate) but cannot get any other info apart from that.

func sessionCheck(token string, db *sql.DB) ( string, error ) {
	sqlSes := "SELECT userid FROM session WHERE id = ?"
	var userid string
	err := db.QueryRow(sqlSes, token).Scan(&userid)
	return userid, err
}

func banCheck(id1, id2 string, db *sql.DB) error {
	sqlBanStmt := "SELECT count(*) FROM ban WHERE bannerid = ? AND bannedid = ?"
	i := 0
	err := db.QueryRow(sqlBanStmt, id1, id2).Scan(&i)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
        	return err
		}
    }
	// The user is banned if i > 0, so he shouldn't be able to see any info of the banner.
	if i > 0 { 
		return structs.NotFoundErr
	}
	return nil
} 

// Users' profile information are available to all except those whom the users ban
func (db *appdbimpl) GetAllUsers(token string) ( []structs.User, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, userid)
	// if errors.Is(err, structs.NotFoundErr) {
	// 		return nil, structs.NotFoundErr
    // } else if err != nil {
    //     return nil, err 
	// } 

	// sqlGetAllUsers := `SELECT id, private, profilename, profilemessage, gender, birthdate, profilephotoid FROM user`

	sqlGetAllUsers := `
		SELECT user.id, profilename, profilemessage, gender, birthdate, profilephotoid
		FROM user WHERE NOT EXISTS ( SELECT * FROM ban WHERE bannerid = user.id AND bannedid = ?)`

	rows, err := db.c.Query(sqlGetAllUsers, userid)
	if err != nil {
        return nil, err
    }
	users := []structs.User{}
	defer rows.Close()

	for rows.Next() {
		user := structs.User{}
		err = rows.Scan(&user.ID, &user.ProfileName, &user.ProfileMessage, &user.Gender, &user.BirthDate, &user.ProfilePhotoID)
		if err != nil {
			return nil, err
		}
		sqlCount := "SELECT count(*) FROM %s WHERE %s = ?"
		sqlPost := fmt.Sprintf(sqlCount, "post", "userid")
		err = db.c.QueryRow(sqlPost, user.ID).Scan(&user.PostCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
           		return nil, err
			}
		}
		sqlfollowing := fmt.Sprintf(sqlCount, "follow", "followerid")
		err = db.c.QueryRow(sqlfollowing, user.ID).Scan(&user.FollowingCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
           		return nil, err
			}
		}
		sqlfollower := fmt.Sprintf(sqlCount, "follow", "followingid")
		err = db.c.QueryRow(sqlfollower, user.ID).Scan(&user.FollowerCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
           		return nil, err
			}
		}
		postids, err := getPostIds(user.ID, db.c)
		if err != nil {
			return nil, err
		}
		user.PostIDs = postids
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return users, nil
}

func getPostIds(user string, db *sql.DB) ( []string, error ) {
	sqlPostIds := `SELECT id FROM post WHERE userid = ?`
	rows, err := db.Query(sqlPostIds, user)
	if err != nil {
		return nil, err
	}
	var postids []string
	defer rows.Close()

	for rows.Next() {
		postid := ""
		err = rows.Scan(&postid)
		if err != nil {
			return nil, err
		}
		postids = append(postids, postid)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return postids, nil
}


// Logs in the user and returns the user ID, and the newly created token. In addition:
// 1. If the user exists before, the string in the third return value will be "200".
// 2. If this is the first time ( for sign up), signs up the user and returns "201" in the third return value.
func (db *appdbimpl) DoLogin(username string) ( string, string, string, error ) {
	nt, err := uuid.NewV4()
	if err != nil {
		return "", "", "", err
	}
	newtoken := nt.String()

	// atomic transaction
	tx, err := db.c.Begin()
	if err != nil {
		return "", "", "", err
	}
	
	defer tx.Rollback()
	
	var user1 string
	err = tx.QueryRow(`SELECT id FROM user WHERE username = ?`, username).Scan(&user1)

	// If the error is ErrNoRows, it means that the user does not exist in our database. So we need to sign him up.
	if errors.Is(err, sql.ErrNoRows) {
		nid, err := uuid.NewV4()
		if err != nil {
			return "", "", "", err
		}
		newid := nid.String()
		
		_, err = tx.Exec(`INSERT INTO user (id, username) VALUES (?, ?)`, newid, username)
		if err != nil {
			return "", "", "", err
		}
		_, err = tx.Exec(`INSERT INTO session (id, userid) VALUES (?, ?)`, newtoken, newid)
		if err != nil {
			return "", "", "", err
		}

		if err = tx.Commit(); err != nil {
			return "", "", "", err

		}
		return newid, newtoken, "201", nil

	} else if err != nil {
		return "", "", "", err
	}

	sqlLogin := `INSERT OR REPLACE INTO session (id, userid) VALUES (?, ?)
				WHERE NOT EXISTS (SELECT * FROM session WHERE userid = ?);`

    _, err = tx.Exec(sqlLogin, newtoken, user1, user1)
	if err != nil {
		return "", "", "", err

	}

	if err = tx.Commit(); err != nil {
		return "", "", "", err

	}

	return user1, newtoken, "200", nil
}

// This deletes the sign-in token
func (db *appdbimpl) LogOut(token string) error {

	_, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	row, err := db.c.Exec(`DELETE FROM session WHERE id = ?;`, token)

	if err != nil {
		return err

	}

	i, err := row.RowsAffected()
	if err != nil {
		return err

	}

	// if i = 0, it means that the user is not logged in
	if i == 0 {
		return structs.NotFoundErr

	}
	return nil
}

// See if the user profile is set to private or not: true or false.
func (db *appdbimpl) GetPrivate(userID, token string) ( bool, error ) {
	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return false, structs.UnAuthErr
	} else if err != nil {
		return false, err
	}

	// Even unfollowed users can get this info even if the user sets his profile to private. Only banned users cannot.
	var val int
	err = db.c.QueryRow(`SELECT private FROM user WHERE id = ? AND NOT EXISTS (
						SELECT * FROM ban WHERE bannerid = user.id AND bannedid = ?)`, userID, userid).Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return false, structs.NotFoundErr		// userID not found or userID banned userid
	} else if err != nil {	// Internal server error
		return false, err
	}
	if val == 1 {
		return true, nil
	}
	return false, nil
}

func (db *appdbimpl) SetPrivate(userID, token string) error {
	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	// If the owner is not the one requesting to modify the protected resource.
	if userid != userID {
		return structs.ForbiddenErr
	}
	_, err = db.c.Exec(`UPDATE user SET user.private = 1 WHERE id = ?`, userID)
	return err
}

func (db *appdbimpl) SetPublic(userID, token string) error {
	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	// If the owner is not the one requesting to modify the protected resource.
	if userid != userID {
		return structs.ForbiddenErr
	}
	_, err = db.c.Exec(`UPDATE user SET user.private = 0 WHERE id = ?`, userID)
	return err
}

func (db *appdbimpl) GetUserProfile(userID, token string) ( *structs.User, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( userID, userid, db.c)
	if errors.Is(err, structs.NotFoundErr) {
		return nil, structs.NotFoundErr
    } else if err != nil {
        return nil, err 
	} 

	user := structs.User{}

	sqlGetUser := `SELECT id, profilename, profilemessage, gender, birthdate, profilephotoid FROM user WHERE user.id = ?`
	
	err = db.c.QueryRow(sqlGetUser, userID).Scan(&user.ID, &user.ProfileName, &user.ProfileMessage, &user.Gender, &user.BirthDate, &user.ProfilePhotoID);

	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr
    } else if err != nil {
        return nil, err 	// Internal Server Error
	} 

	sqlCount := "SELECT count(*) FROM %s WHERE %s = ?"
	sqlPost := fmt.Sprintf(sqlCount, "post", "userid")
	err = db.c.QueryRow(sqlPost, user.ID).Scan(&user.PostCount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	sqlfollowing := fmt.Sprintf(sqlCount, "follow", "followerid")
	err = db.c.QueryRow(sqlfollowing, user.ID).Scan(&user.FollowingCount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	sqlfollower := fmt.Sprintf(sqlCount, "follow", "followingid")
	err = db.c.QueryRow(sqlfollower, user.ID).Scan(&user.FollowerCount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	postids, err := getPostIds(user.ID, db.c)
	if err != nil {
		return nil, err
	}
	user.PostIDs = postids
	return &user, nil
}

// Successfully created some fields that did not exist: 201, else 204.
func (db *appdbimpl) UpdateUserProfile(userID, token string, pd *structs.ProfileDetail) ( string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.UnAuthErr
	} else if err != nil {
		return "", err
	}
	if userID != userid {
		return "", structs.ForbiddenErr
	}

	// atomic transaction
	tx, err := db.c.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Update the propriate values.
	arr1 := [...]string{pd.ProfileName, pd.ProfileMessage, pd.Gender, pd.BirthDate}
	arr2 := [...]string{"profilename", "profilemessage", "gender", "birthdate"}
	valCreated := "204"

	for idx, elem := range arr1 {
		if elem != "" {

			// Check if the field was not empty
			str := ""
			sqlCheck := fmt.Sprintf("SELECT %s FROM user WHERE id = ?", arr2[idx])
			err = tx.QueryRow(sqlCheck, userID).Scan(&str)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return "", err
				}
				// valCreated = "201"
			}
			if str == "" {
				valCreated = "201"
			}

			// Update the field
			sqlUpdateUser := fmt.Sprintf("UPDATE user SET %s = ? WHERE id = ?", arr2[idx])
			_, err = tx.Exec(sqlUpdateUser, elem, userID)
			if err != nil {
				return "", err
			}
		}
	}

	// Commit the transaction.
    if err = tx.Commit(); err != nil {
		return "", err
	}
	return valCreated, nil
}

func (db *appdbimpl) DeleteUser(userID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr

	} else if err != nil {
		return err
	}

	// If the owner is not the one requesting to modify the protected resource.
	if userid != userID {
		return structs.ForbiddenErr
	}

	// atomic transaction
	tx, err := db.c.Begin()
	if err != nil {
		return err

	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT photoid FROM post WHERE userid = ?", userid)
	if err != nil {
        return err

    }
	defer rows.Close()

	// For each photo of the user to be deleted, delete it from disk
	for rows.Next() {

		var photoID string
        err = rows.Scan(&photoID)
        if err != nil {
            return err

        }

		file := filepath.Join("./pictures", photoID + ".png")

		err = os.Remove(file)
		if err != nil {
			return err

		}
	}

	if err = rows.Err(); err != nil {
        return err
	}

	// Cascade delete -> deletes any row in any table that has reference to the user in the database.
	_, err = tx.Exec(`DELETE FROM user WHERE id = ? ;`, userID)
	if err != nil {
        return err

	}

	err = tx.Commit()
    if err != nil {
        return err

	}
	
	return nil
}


func (db *appdbimpl) SetMyUserName(userID, token string, username string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	// If the owner is not the one requesting to modify the protected resource.
	if userid != userID {
		return structs.ForbiddenErr
	}
	_, err = db.c.Exec(`UPDATE user SET username = ? WHERE id = ? ;`, username, userID)
	return err
}

// The returned uuid is the photo ID to be retrieved from directory.
func (db *appdbimpl) GetUserProfilePicture(userID, token string) ( string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.UnAuthErr
	} else if err != nil {
		return "", err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( userID, userid, db.c)
	if errors.Is(err, structs.NotFoundErr) {
		return "", structs.NotFoundErr
    } else if err != nil {
        return "", err 
	} 

	// Everyone can get the user profile picture including those who did not follow him even if the user sets his profile to private except those banned
	sqlGetUser := `SELECT profilephotoid FROM user WHERE user.id = ?`

	var pid string
	err = db.c.QueryRow(sqlGetUser, userID).Scan(&pid);
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.NotFoundErr
    } else if err != nil {
        return "", err 
	} 

	return pid, nil
}

// The returned uuid is the photo ID to be sent. "204" if exists, else "201"
func (db *appdbimpl) ChangeUserProfilePicture(userID, token string, r *http.Request) ( string, string, error) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", structs.UnAuthErr
	} else if err != nil {
		return "", "", err
	}
	// If the owner is not the one requesting to modify the protected resource.
	if userid != userID {
		return "", "", structs.ForbiddenErr
	}

	tx, err := db.c.Begin()
	if err != nil {
        return "", "", err
    }
	defer tx.Rollback()


	pid := ""
	err = tx.QueryRow("SELECT profilephotoid FROM user WHERE id = ?", userid).Scan(&pid);
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", "", err
		}
	}

	if pid != "" {

		err = updateProfPic(pid, r)
		if err != nil {
            return "", "", err
        }

		return pid, "204", nil
	}

	uid, err := uuid.NewV4()
	if err != nil {
        return "", "", err
	}
	_, err = tx.Exec("UPDATE user SET profilephotoid = ? WHERE id = ?", uid.String(), userid)
	if err != nil {
        return "", "", err

	}

	err = updateProfPic(uid.String(), r)
	if err != nil {
		return "", "", err
	}

	err = tx.Commit()
    if err != nil {
		return "", "", err

	}

	return uid.String(), "201", nil
}

func updateProfPic(photoID string, r *http.Request) error {
	
	// Form contains the parsed form data, including both the URL field's query parameters and 
	// the PATCH, POST, or PUT form data. This field is only available after ParseForm is called.
	// But ParseMultipartForm automatically calls ParseForm

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
        return err

	}

	// (multipart.File, *multipart.FileHeader, error)
	// FormFile returns the first file for the provided form key. 
	// FormFile calls ParseMultipartForm and ParseForm if necessary.
	// _ for getting the filenames, extensions, etc.
	photo_file, _, err := r.FormFile("photo")
	if err != nil {
        return err

	}

	file := filepath.Join("./pictures", photoID + ".png")

	// Create creates or truncates the named file. If the file already exists, it is truncated. 
	// If the file does not exist, it is created with mode 0666 (before umask). 
	// If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. 
	// If there is an error, it will be of type *PathError.

	img, err := os.Create(file)

	if err != nil {
        return err

	}

	defer img.Close()


	// Copy copies from src to dst until either EOF is reached on src or an error occurs. 
	// It returns the number of bytes copied and the first error encountered while copying, if any.

	// A successful Copy returns err == nil, not err == EOF. 
	// Because Copy is defined to read from src until EOF, it does not treat an EOF from Read as an error to be reported.

	_, err = io.Copy(img, photo_file)

	if err != nil {
        return err

	}

	return nil
}


func (db *appdbimpl) DeleteUserProfilePicture(userID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr

	} else if err != nil {
		return err

	}
	
	if userID != userid {
		return structs.ForbiddenErr

	}

	tx, err := db.c.Begin()
	if err != nil {
        return err
    }

	defer tx.Rollback()

	var photoID string
	err = tx.QueryRow("SELECT profilephotoid FROM user WHERE id = ?", userid).Scan(&photoID)

	if err != nil {
		// if !errors.Is(err, sql.ErrNoRows) {
       	//  	return err
		// }

		return err
    }

	file := filepath.Join("./pictures", photoID + ".png")

	err = os.Remove(file)
    if err != nil {
        return err

    }


	_, err = tx.Exec("UPDATE user SET profilephotoid = '' WHERE id = ?", userid)
	if err != nil {
        return err

	}

	if err = tx.Commit(); err != nil {
		return err

	}

	return nil
}



// First slice with user IDs and second slice with post IDs that corresponds to the search term.
func (db *appdbimpl) Search(token string, search string) ( []string, []string , error) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, structs.UnAuthErr
	} else if err != nil {
		return nil, nil, err
	}

	var users []string

	// Notice that it is the profilename and not the username
	rows, err := db.c.Query("SELECT id FROM user WHERE profilename = ?", search)
	if err != nil {
        return nil, nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var id string
        err = rows.Scan(&id)
        if err != nil {
            return nil, nil, err
        }

		// If the client is banned, then for the client, the user does not exist. So return 404.
		err  = banCheck( id, userid, db.c)
		if err == nil {
			users = append(users, id)
		} else if errors.Is(err, structs.NotFoundErr) {
			continue
		} else if err != nil {
			return nil, nil, err 
		}
	}

	if err = rows.Err(); err != nil {
        return nil, nil, err
	}

	var posts []string
	rows, err = db.c.Query(`SELECT post.id FROM post INNER JOIN hashtag ON post.id = hashtag.postid WHERE hashtag = ?
							AND NOT EXISTS (SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?) `, search, userid)
	if err != nil {
        return nil, nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var pid string
        err = rows.Scan(&pid)
        if err != nil {
            return nil, nil, err
        }
		posts = append( posts, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, nil,  err
	}

	return users, posts, nil
}

// First slice with user followers' IDs and second slice with user followings' IDs
func (db *appdbimpl) GetUserFollows(userID, token string) ( []string, []string , error) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, structs.UnAuthErr
	} else if err != nil {
		return nil, nil, err
	}
	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( userID, userid, db.c)
	if errors.Is(err, structs.NotFoundErr) {
			return nil, nil, structs.NotFoundErr
    } else if err != nil {
        return nil, nil, err 
	} 

	// Anyone who does not follow a user and the user sets his profile to private, then cannot see any other details of the user except the user profile information including whether or not the user set it to private.
	sqlfollowCheck := `SELECT count(*) FROM follow as f1 INNER JOIN user ON f1.followingid = user.id WHERE 
							f1.followingid = ? AND user.private = 1 AND NOT EXISTS (
								SELECT * FROM follow as f2 WHERE f2.followingid = f1.followingid AND f2.followerid = ? )`

    i := 0
	err = db.c.QueryRow(sqlfollowCheck, userID, userid).Scan(&i) 
	if i > 0 { 		// if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, structs.NotFoundErr
	} else if err != nil {
		return nil, nil, err 
	} 

	var followerids []string
	rows, err := db.c.Query("SELECT f1.followerid FROM follow as f1 WHERE f1.followingid = ?", userID)
	if err != nil {
        return nil, nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var pid string
        err = rows.Scan(&pid)
        if err != nil {
            return nil, nil, err
        }
		followerids = append(followerids, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, nil, err
	}

	var followingids []string
	rows, err = db.c.Query("SELECT f1.followingid FROM follow as f1 WHERE f1.followerid = ?", userID)
	if err != nil {
        return nil, nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var pid string
        err = rows.Scan(&pid)
        if err != nil {
            return nil, nil, err
        }
		followingids = append(followingids, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, nil, err
	}

	return followerids, followingids, nil
}

// The user with userID follows the user with followID
func (db *appdbimpl) FollowUser(userID, followID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	if userID != userid {
		return structs.ForbiddenErr
	}
	if userID == followID {
		return structs.ForbiddenErr //errors.New("You cannot follow yourself")
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( followID, userID, db.c)
	if errors.Is(err, structs.NotFoundErr) {
		return structs.NotFoundErr
    } else if err != nil {
        return err 
	} 

	_, err = db.c.Exec(`INSERT OR IGNORE INTO follow (followerid, followingid) VALUES (?, ?)`, userID, followID,)

	return err
}

// The user with userID unfollows the user with followID
func (db *appdbimpl) UnfollowUser(userID, followID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return  structs.UnAuthErr
	} else if err != nil {
		return err
	}
	if userID != userid {
		return structs.ForbiddenErr
	}
	if userID == followID {
		return structs.ForbiddenErr //errors.New("You cannot unfollow yourself")
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( followID, userID, db.c)
	if errors.Is(err, structs.NotFoundErr) {
		return structs.NotFoundErr
    } else if err != nil {
        return err 
	} 

	_, err = db.c.Exec(`DELETE FROM follow WHERE followerid = ? AND followingid = ?`, userID, followID)

	return err		// Internal Server Error if error is not nil
}

// users who are banned by the userID
func (db *appdbimpl) GetBanUsers(userID, token string) ( []string, error) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}
	// Only the user himself can see his banned list
	if userID != userid {
		return nil, structs.ForbiddenErr
	}
	var bannedids []string
	rows, err := db.c.Query(`SELECT b1.bannedid FROM ban as b1 WHERE b1.bannerid = ?`, userID)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var pid string
        err = rows.Scan(&pid)
        if err != nil {
            return nil, err
        }
		bannedids = append(bannedids, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return bannedids, nil
}

// The user with userID bans the user with banID
func (db *appdbimpl) BanUser(userID, banID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	if userID != userid {
		return structs.ForbiddenErr
	}

	if userID == banID {
		return structs.ForbiddenErr //errors.New("You cannot ban yourself")
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, banID, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 	return structs.NotFoundErr
    // } else if err != nil {
    //     return err 
	// } 

	sqlBan := `INSERT OR IGNORE INTO ban (bannerid, bannedid) VALUES (?, ?)`
	_, err = db.c.Exec(sqlBan, userID, banID)

	return err
}

// The user with userID unbans the user with banID
func (db *appdbimpl) UnbanUser(userID, banID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	if userID != userid {
		return structs.ForbiddenErr
	}
	if userID == banID {
		return structs.ForbiddenErr //errors.New("You cannot ban/unban yourself")
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, userID, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 	return structs.NotFoundErr
    // } else if err != nil {
    //     return err 
	// } 

	sqlUnban := `DELETE FROM ban WHERE bannerid = ? AND bannedid = ?`
	_, err = db.c.Exec(sqlUnban, userID, banID)

	return err
}

// Check if the photoID exists and the user is allowed to access it -> nil, else error.
func (db *appdbimpl) GetSinglePhoto(photoID, token string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	// err  = banCheck( userID, userID, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 	return structs.NotFoundErr
    // } else if err != nil {
    //     return err 
	// } 

	sqlGetSinglephoto := `
	SELECT post.photoid FROM post INNER JOIN user ON post.userid = user.id WHERE post.photoid = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`
	var pid string
	err = db.c.QueryRow(sqlGetSinglephoto, photoID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
        return err
    }
	if pid == photoID { // --
        return nil
    }
    return nil
}

// This sends the post with all its metadata attached
// It should be named GetPost, but due to the project requirements, it is called GetPhoto.
func (db *appdbimpl) GetPhoto(postID, token string) ( *structs.Post, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, userID, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 	return structs.NotFoundErr
    // } else if err != nil {
    //     return err 
	// } 

	sqlGetPhoto := `
	SELECT * FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`
	var post = structs.Post{}
	err = db.c.QueryRow(sqlGetPhoto, postID, userid, userid, userid).Scan(&post.ID, &post.PhotoID, &post.UserID, &post.Caption, &post.DateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr	// photoID not found or user banned or (private = true and unfollow)
	} else if err != nil {
		return nil, err
	}
	var sqlCount = "SELECT count(*) FROM %s WHERE postid = ?"
	sqlLikeCount := fmt.Sprintf(sqlCount, "like")
	err = db.c.QueryRow(sqlLikeCount, postID).Scan(&post.LikeCount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	sqlCommentCount := fmt.Sprintf(sqlCount, "comment")
	err = db.c.QueryRow(sqlCommentCount, postID).Scan(&post.CommentCount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	hashtags, err := getHashtags(postID, db.c)

	if err != nil {
		return nil, err

	}
	
	post.Hashtags = hashtags

	return &post, nil
}

func getHashtags ( postID string, db *sql.DB) ( []string, error ) {
	var hashtags []string
	rows, err := db.Query("SELECT hashtag FROM hashtag WHERE postid = ?", postID)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var pid string
        err = rows.Scan(&pid)
        if err != nil {
            return nil, err
        }
		hashtags = append(hashtags, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return hashtags, nil
}

// GetPhotos gets all posts of all users who did not set their profiles to private and did not ban the user.
func (db *appdbimpl) GetPhotos( token string ) ( []structs.Post, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, userID, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 	return structs.NotFoundErr
    // } else if err != nil {
    //     return err 
	// } 

	sqlGetPhotos := `
	SELECT * FROM post INNER JOIN user ON post.userid = user.id WHERE post.userid = ? OR 
					(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
						AND (user.private = 0 OR (user.private = 1 AND 
							EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
													)
							)
					)`
	rows, err := db.c.Query(sqlGetPhotos, userid, userid, userid)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	var posts = []structs.Post{}

	for rows.Next() {
		var post = structs.Post{}
        err = rows.Scan(&post.ID, &post.PhotoID, &post.UserID, &post.Caption, &post.DateTime)
        if err != nil {
            return nil, err
        }
		var sqlCount = "SELECT count(*) FROM %s WHERE postid = ?"
		sqlLikeCount := fmt.Sprintf(sqlCount, "like")
		err = db.c.QueryRow(sqlLikeCount, post.ID).Scan(&post.LikeCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		sqlCommentCount := fmt.Sprintf(sqlCount, "comment")
		err = db.c.QueryRow(sqlCommentCount, post.ID).Scan(&post.CommentCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}

		hashtags, err := getHashtags(post.ID, db.c)
		if err != nil {
			return nil, err
	
		}	
		
		post.Hashtags = hashtags
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}
	
	return posts, nil
}

// Get the list of posts posted by the given user's followings (including the user himself).
func (db *appdbimpl) GetMyStream(userID, token string ) ( []structs.Post, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}
	// only the user himself can get the posts of his followings.
	if userid != userID {
		return nil, structs.ForbiddenErr
	}
	sqlGetStream := `
	SELECT * FROM post WHERE post.userid = ? OR 
						(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
						AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)`

	rows, err := db.c.Query(sqlGetStream, userid, userid, userid)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	var posts = []structs.Post{}

	for rows.Next() {
		var post = structs.Post{}
        err = rows.Scan(&post.ID, &post.PhotoID, &post.UserID, &post.Caption, &post.DateTime)
        if err != nil {
            return nil, err
        }
		var sqlCount = "SELECT count(*) FROM %s WHERE postid = ?"
		sqlLikeCount := fmt.Sprintf(sqlCount, "like")
		err = db.c.QueryRow(sqlLikeCount, post.ID).Scan(&post.LikeCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		sqlCommentCount := fmt.Sprintf(sqlCount, "comment")
		err = db.c.QueryRow(sqlCommentCount, post.ID).Scan(&post.CommentCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}

		hashtags, err := getHashtags(post.ID, db.c)
		if err != nil {
			return nil, err
	
		}
		
		
		post.Hashtags = hashtags
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}
	
	return posts, nil
}

// Get the list of posts posted by the given user.
func (db *appdbimpl) GetUserPhotos(userID, token string ) ( []structs.Post, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	err  = banCheck( userID, userid, db.c)
	if errors.Is(err, structs.NotFoundErr) {
		return nil, structs.NotFoundErr
    } else if err != nil {
        return nil, err 
	} 

	sqlGetUserPhotos := `
	SELECT * FROM post INNER JOIN user ON post.userid = user.id WHERE post.userid = ? AND 
	(post.userid = ? OR user.private = 0 OR (user.private = 1 AND 
								EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
											)
	)`

	rows, err := db.c.Query(sqlGetUserPhotos, userID, userid, userid)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	var posts = []structs.Post{}

	for rows.Next() {
		var post = structs.Post{}
        err = rows.Scan(&post.ID, &post.PhotoID, &post.UserID, &post.Caption, &post.DateTime)
        if err != nil {
            return nil, err
        }
		var sqlCount = "SELECT count(*) FROM %s WHERE postid = ?"
		sqlLikeCount := fmt.Sprintf(sqlCount, "like")
		err = db.c.QueryRow(sqlLikeCount, post.ID).Scan(&post.LikeCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		sqlCommentCount := fmt.Sprintf(sqlCount, "comment")
		err = db.c.QueryRow(sqlCommentCount, post.ID).Scan(&post.CommentCount)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}

		hashtags, err := getHashtags(post.ID, db.c)
		if err != nil {
			return nil, err
	
		}
		
		
		post.Hashtags = hashtags
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}
	
	return posts, nil
}

func (db *appdbimpl) UploadPhoto( userID, token string, caption string, hashtags []string, r *http.Request ) ( string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.UnAuthErr
	} else if err != nil {
		return "", err
	}
	
	if userid != userID {
		return "", structs.ForbiddenErr
	}

	uid, err := uuid.NewV4()
	if err != nil {
        return "", err
    }
	photoID := uid.String()

	uid, err = uuid.NewV4()
	if err != nil {
        return "", err
    }
	postID := uid.String()

	datetime := time.Now().Format("2006-01-02T15:04:05Z")

    sqlUploadPhoto := `
		INSERT INTO post (id, photoid, userid, caption, datetime)
		VALUES (?, ?, ?, ?, ?);
	`

	// atomic transaction
	tx, err := db.c.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(sqlUploadPhoto, postID, photoID, userID, caption, datetime)
	if err != nil {
        return "", err
    }

	for _, hashtag := range hashtags {
		_, err = tx.Exec("INSERT INTO hashtag (hashtag, postid) VALUES (?,?);", postID, hashtag)
        if err != nil {
            return "", err
        }
	}

	err = updateProfPic(postID, r)
	if err != nil {
        return "", err

    }

	if err = tx.Commit(); err != nil {
		return "", err
	}
	
	return postID, nil
}


func (db *appdbimpl) ModifyCaption( userID, token, postID, caption string ) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	if userid != userID {
		return structs.ForbiddenErr
	}
	_, err = db.c.Exec("UPDATE post SET caption = ? WHERE id = ?", caption, postID)
	return err
}

// Deletes a post with the given post ID together with the photo, caption, likes and comments, etc.
func (db *appdbimpl) DeletePhoto( userID, token, postID string) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	if userid != userID {
		return structs.ForbiddenErr
	}

	// atomic transaction
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	
	defer tx.Rollback()

	// The photo id to be deleted from disk
	var photoID string
	err = tx.QueryRow("SELECT photoid FROM post WHERE post.id = ?", postID).Scan(&photoID)

	if err != nil {
		// if !errors.Is(err, sql.ErrNoRows) {
       	//  	return err
		// }

		return err
    }


	// sqlDeleteP := `
	// 				DELETE FROM post WHERE post.id = ?;
	// 				DELETE FROM comment WHERE postid =?;
	// 				DELETE FROM like WHERE postid = ?;
	// 				DELETE FROM hashtag WHERE postid = ?;
	// `
	// _, err = tx.Exec(sqlDeleteP, postID, postID, postID, postID)

	// Cascade delete
	sqlDeleteP := `DELETE FROM post WHERE post.id = ?;`
	_, err = tx.Exec(sqlDeleteP, postID)
	if err != nil {
		return err
	}

	file := filepath.Join("./pictures", photoID + ".png")

	err = os.Remove(file)
    if err != nil {
        return err

    }

	// Commit the transaction.
    if err = tx.Commit(); err != nil {
        return err
    }

	return nil
}

// Return "204" if hashtag already exists, else "201".
func (db *appdbimpl) AddHashtag( userID, token, postID string, hashtag string ) ( string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.UnAuthErr
	} else if err != nil {
		return "", err
	}
	if userid != userID {
		return "", structs.ForbiddenErr
	}

	// Check if the poster is the one trying to add the hashtag.
	var i int64 = 0
	err = db.c.QueryRow("SELECT count(*) FROM post WHERE id = ? AND userid = ?", postID, userID).Scan(&i)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
        	return "", err
		}
    }
	if i <= 0 { 
		return "", structs.ForbiddenErr
	}

	row, err := db.c.Exec("INSERT OR IGNORE INTO hashtag (hashtag, postid) VALUES (?, ?)", hashtag, postID)
	if err != nil {
        return "", err
    }
	i, err = row.RowsAffected()
	if err != nil {
		return "", err
	}
	// if i = 0, it means that the value already exists and so the db ignores it, so return 204
	if i == 0 {
		return "204", nil
	}
	return "201", nil
}

func (db *appdbimpl) DeleteHashtag( userID, token, postID string, hashtag string ) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}
	if userid != userID {
		return structs.ForbiddenErr
	}
	// Check if the poster is the one trying to delete the hashtag.
	i := 0
	err = db.c.QueryRow("SELECT count(*) FROM post WHERE id = ? AND userid = ?", postID, userID).Scan(&i)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
        	return err
		}
    }
	if i <= 0 { 
		return structs.ForbiddenErr
	}

	_, err = db.c.Exec("DELETE FROM hashtag WHERE hashtag = ? AND postID = ? ", hashtag, postID)

	return err
}


func (db *appdbimpl) GetPostHashtags( token, postID string ) ( []string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// If the client is banned, then for the client, the user does not exist. So return 404.
	// err  = banCheck( userID, userid, db.c)
	// if errors.Is(err, structs.NotFoundErr) {
	// 		return nil, structs.NotFoundErr
    // } else if err != nil {
    //     return nil, err 
	// } 
	

	// Check everything first.
	sqlGetPostH := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`
	var pid string
	err = db.c.QueryRow(sqlGetPostH, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return nil, err
	}

	// Everything is ok. Now get the hashtags
	var hashtags []string
	rows, err := db.c.Query("SELECT hashtag FROM hashtag WHERE postid = ? ", postID)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		pid = ""
        err = rows.Scan(&pid)
        if err != nil {
            return nil, err
        }
		hashtags = append(hashtags, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return hashtags, nil
}

// Get the like-count and the user IDs who liked the post.
func (db *appdbimpl) GetLikes( token, postID string ) ( *structs.Like, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// Check everything first.
	sqlGetLikes := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlGetLikes, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return nil, err
	}

	// Everything is ok. Now get the requested data
	var userids []string
	rows, err := db.c.Query("SELECT userid FROM like WHERE postid = ? ", postID)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		pid = ""
        err = rows.Scan(&pid)
        if err != nil {
            return nil, err
        }
		userids = append(userids, pid)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	likes := structs.Like{}
	likes.LikeCount = len(userids)
	likes.UserIDs = userids

	return &likes, nil
}

func (db *appdbimpl) LikePhoto( userID, token, postID string ) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	if userid != userID {
		return structs.ForbiddenErr
	}

	// Check everything first.
	sqlGetLikes := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlGetLikes, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnlikePhoto( userID, token, postID string ) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	if userid != userID {
		return structs.ForbiddenErr
	}

	// Check everything first.
	sqlUnLikes := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlUnLikes, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return err
	}

	return nil
}

// Get all the comments of a given post.
func (db *appdbimpl) GetPhotoComments( token, postID string ) ( []structs.Comment, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	// Check everything first.
	sqlCommentsCheck := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlCommentsCheck, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return nil, err
	}

	// Everything is ok. Now get the requested data
	var comments = []structs.Comment{}
	rows, err := db.c.Query("SELECT * FROM comment WHERE postid = ? ", postID)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var comment = structs.Comment{}
        err = rows.Scan(&comment.ID , &comment.PostID , &comment.UserID , &comment.Message , &comment.DateTime)
        if err != nil {
            return nil, err
        }
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}

	return comments, nil
}

//  Places a new comment and returns the newly created comment ID.
func (db *appdbimpl) CommentPhoto( token, postID string, message string ) ( string, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.UnAuthErr
	} else if err != nil {
		return "", err
	}

	// Check everything first.
	sqlComment := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlComment, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return "", structs.NotFoundErr	// not found or banned or (private = true and unfollow)
	} else if err != nil {
		return "", err
	}

	uid, err := uuid.NewV4()
	if err != nil {
        return "", err
    }
	commentID := uid.String()
	datetime := time.Now().Format("2006-01-02T15:04:05Z")

    sqlcomment := `
		INSERT INTO comment (id, postid, userid, message, datetime)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err = db.c.Exec(sqlcomment, commentID, postID, userid, message, datetime)
	if err != nil {
        return "", err
    }

	return commentID, nil
}

func (db *appdbimpl) GetComment( token, commentID string ) ( *structs.Comment, error ) {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.UnAuthErr
	} else if err != nil {
		return nil, err
	}

	var comment = structs.Comment{}
	err = db.c.QueryRow("SELECT * FROM comment WHERE comment.id = ? ", commentID).Scan(&comment.ID , &comment.PostID , &comment.UserID , &comment.Message , &comment.DateTime)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr
	} else if err != nil {
		return nil, err
	}

	// Even the person who placed the comment cannot get the comment if the owner of the post bans him... So late test
	sqlCommentsCheck := `
	SELECT post.id FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlCommentsCheck, comment.PostID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, structs.NotFoundErr	// banned or (private = true and unfollow)
	} else if err != nil {
		return nil, err
	}

	return &comment, nil
}

//  Deletes a given comment by the user. 
// The user is either the one who placed the comment or on whose post the comment was placed.
func (db *appdbimpl) UncommentPhoto( token, commentID string ) error {

	userid, err := sessionCheck(token, db.c)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.UnAuthErr
	} else if err != nil {
		return err
	}

	var postID, userID string
	err = db.c.QueryRow("SELECT postid, userid FROM comment WHERE comment.id = ? ", commentID).Scan(&postID, &userID)

	if errors.Is(err, sql.ErrNoRows) {
		return structs.NotFoundErr
	} else if err != nil {
		return err
	}

	// Even the person who placed the comment cannot delete/get/see the comment if the owner of the post bans him...
	sqlCommentsCheck := `
	SELECT post.userid FROM post INNER JOIN user ON post.userid = user.id WHERE post.id = ? AND 
		( post.userid = ? OR 
				(NOT EXISTS ( SELECT * FROM ban WHERE bannerid = post.userid AND bannedid = ?)
				AND (user.private = 0 OR 
						(user.private = 1 AND EXISTS (SELECT * FROM follow WHERE followerid = ? AND followingid = post.userid)
						)
					)
				)
		)`

	var pid string
	err = db.c.QueryRow(sqlCommentsCheck, postID, userid, userid, userid).Scan(&pid)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.NotFoundErr	// banned or (private = true and unfollow)
	} else if err != nil {
		return err
	}

	// Check that the user trying to delete the post is the one who placed it or on whose post it was placed
	if userID != userid || pid != userid {
		return structs.ForbiddenErr

	}

	_, err = db.c.Exec("DELETE FROM comment WHERE comment.id = ? ", commentID)

	return err

}