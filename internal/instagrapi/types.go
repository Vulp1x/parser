package instagrapi

type UserShort struct {
	Pk       string `json:"pk"`
	Username string `json:"username"`
}

type UserTag struct {
	User UserShort `json:"user"`
	X    float64   `json:"x"`
	Y    float64   `json:"y"`
}
