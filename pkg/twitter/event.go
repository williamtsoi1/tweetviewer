package twitter

// SimpleTweet represents simple Twitter message
type SimpleTweet struct {
	CreatedAt string             `json:"created_at"`
	IDStr     string             `json:"id_str"`
	Text      string             `json:"text"`
	FullText  string             `json:"full_text"`
	User      *SimpleTwitterUser `json:"user"`
}

// SimpleTwitterUser represents author of simple Twitter message
type SimpleTwitterUser struct {
	DefaultProfileImage bool   `json:"default_profile_image"`
	Description         string `json:"description"`
	FollowersCount      int    `json:"followers_count"`
	IDStr               string `json:"id_str"`
	Name                string `json:"name"`
	ProfileImageURL     string `json:"profile_image_url"`
	ScreenName          string `json:"screen_name"`
	Verified            bool   `json:"verified"`
}
