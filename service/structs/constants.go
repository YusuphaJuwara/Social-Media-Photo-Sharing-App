package structs

/*
Here, I define some general constants needed throughout the project.
1. Status codes
2. Sql table creation statements' constants
3. Patterns, min and max length constants
*/

import (

	"errors"
)

var (
	UnAuthErr = errors.New("Unauthorized user")
	ForbiddenErr = errors.New("Forbidden")
	NotFoundErr = errors.New("Not found")
	BadReqErr 	= errors.New("Bad request")
)


// 1. Sql table creation statements' constants
const (

	UserS = `CREATE TABLE IF NOT EXISTS user (
		id TEXT NOT NULL PRIMARY KEY, 
		username TEXT NOT NULL, 
        private INT DEFAULT 0,
		profilename TEXT DEFAULT '', 
		profilemessage TEXT DEFAULT '',
		gender TEXT DEFAULT '', 
		birthdate text DEFAULT '',
		profilephotoid TEXT DEFAULT ''
		);`

	CommentS = `CREATE TABLE IF NOT EXISTS comment (
		id TEXT NOT NULL PRIMARY KEY, 
		postid TEXT NOT NULL REFERENCES post(id) On DELETE CASCADE,
		userid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE,
		message TEXT DEFAULT '',
		datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	// can store and retrieve TIMESTAMP as strings 2022-12-28T03:19:15Z

	PostS = `CREATE TABLE IF NOT EXISTS post (
		id TEXT NOT NULL PRIMARY KEY, 
		photoid TEXT NOT NULL,
		userid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE,
		caption TEXT DEFAULT '',
		datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	LikeS = `CREATE TABLE IF NOT EXISTS like (
		postid TEXT NOT NULL REFERENCES post(id) On DELETE CASCADE, 
		userid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE,
		PRIMARY KEY (postid, userid)
		);`

	FollowS = `CREATE TABLE IF NOT EXISTS follow (
		followerid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE, 
		followingid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE,
		PRIMARY KEY (followerid, followingid)
		);`
	
	BanS = `CREATE TABLE IF NOT EXISTS ban (
		bannerid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE, 
		bannedid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE,
		PRIMARY KEY (bannerid, bannedid)
		);`

	HashtagS = `CREATE TABLE IF NOT EXISTS hashtag (
		hashtag TEXT NOT NULL, 
		postid TEXT NOT NULL REFERENCES post(id) On DELETE CASCADE,
		PRIMARY KEY (hashtag, postid)
		);`
	
	SessionS = `CREATE TABLE IF NOT EXISTS session (
		id TEXT NOT NULL PRIMARY KEY, 
		userid TEXT NOT NULL REFERENCES user(id) On DELETE CASCADE
		);`

)

var SqlStmtList = []string{UserS, CommentS, PostS, LikeS, FollowS, BanS, HashtagS, SessionS}


// 3. Patterns, min and max length constants
const (
	// '^(?=.*?[a-zA-Z]).{8,20}$' # 8 to 20 chars of at least 1 alphabet. No new line char
	ProfileNamePattern			= 		"^.*[a-zA-Z].*$"
	ProfileNameMinLen 			= 		8
	ProfileNameMaxLen 			= 		20


	// '^(?=.*?[a-zA-Z])[a-zA-Z0-9]{8,20}$' # at least one alphabet and 7 or more other alphanumeric characters. 8-20
	UsernamePattern 			=		"^\"[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*\"$"
	UsernameMinLen 				= 		8
	UsernameMaxLen 				= 		20

	// string of length 1 to 20. ^(?=.*?[a-zA-Z])[a-zA-Z0-9]{1,20}$ -> at least 1 alphabet
	HashtagPattern 				=		"^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$"
	HashtagMinLen 				= 		1
	HashtagMaxLen 				= 		20

	// '(?=.*?[a-zA-Z])[^]'  # Any char including new lines of at least 2 alphabets (2-1000)
	MessagePattern 				=		"[a-zA-Z]"
	MessageMinLen 				= 		2
	MessageMaxLen 				= 		1000

	// valid date "2006-01-02"
	DatePattern 				=		"^[0-9]{4}-[0-9]{2}-[0-9]{2}$"

	// valid date and time "2006-01-02T15:04:05Z'
	DateTimePattern	 			=		"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$"

)