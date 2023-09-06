package structs

type ProfileDetail struct {
	ProfileName    string `json:"profile-name"`
	ProfileMessage string `json:"profile-message"`

	// comes as a string, then checked to see if it's a valid enum in {Male, Female}
	Gender string `json:"gender"`

	// comes as a string, then parsed to see if it's a valid date. "2006-01-02"
	BirthDate string `json:"birth-date"`
}

type User struct {
	ID             string   `json:"user-id"`
	ProfileName    string   `json:"profile-name"`
	Gender         string   `json:"gender"`
	BirthDate      string   `json:"birth_date"`
	ProfilePhotoID string   `json:"profile-photo-id"`
	PostIDs        []string `json:"user-post-ids"`
	PostCount      int      `json:"post-count"`
	FollowerCount  int      `json:"follower-count"`
	FollowingCount int      `json:"following-count"`
	ProfileMessage string   `json:"profile-message"`
}

type Comment struct {
	ID       string `json:"comment-id"`
	PostID   string `json:"post-id"`
	UserID   string `json:"user-id"`
	Message  string `json:"message"`
	DateTime string `json:"date-time"` // Valid date and time as "2006-01-02T15:04:05Z
}

type Post struct {
	ID           string `json:"post-id"`
	PhotoID      string `json:"photo-id"`
	UserID       string `json:"user-id"`
	Caption      string `json:"caption"`   // Just a message
	DateTime     string `json:"date-time"` // Valid date and time as "2006-01-02T15:04:05Z
	LikeCount    int    `json:"like-count"`
	CommentCount int    `json:"comment-count"`

	// Hashtag contains at least one alphabet and zero or more other alphanumeric characters in [a-zA-Z0-9].
	// pattern: /^(?=.*?[a-zA-Z])[a-zA-Z0-9]{1,20}$/
	Hashtags []string `json:"hashtags"`
}

// type ModHashtags struct {
// 	Add bool	`json:"add"`
// 	Hashtags []string	`json:"hashtags"`
// }

type Like struct {
	LikeCount int      `json:"like-count"`
	UserIDs   []string `json:"user-ids"`
}

type IDs struct {
	UserIDs []string `json:"userIDs"`
	PostIDs []string `json:"postIDs"`
}

type Follow struct {
	FollowerIDs  []string `json:"followers-array"`
	FollowingIDs []string `json:"followings-array"`
}

type Ban struct {
	BannerIDs []string `json:"banners"`
	BannedIDs []string `json:"banneds"`
}
