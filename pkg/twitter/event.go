package twitter

// SimpleTweet represents simple Twitter message
type SimpleTweet struct {
	CreatedAt string             `json:"created_at"`
	IDStr     string             `json:"id_str"`
	Text      string             `json:"text"`
	User      *SimpleTwitterUser `json:"user"`
}

// SimpleTwitterUser represents author of simple Twitter message
type SimpleTwitterUser struct {
	ProfileImageURL string `json:"profile_image_url"`
	ScreenName      string `json:"screen_name"`
}
