package generics

import "go-lbapp/model"

// SyncMsg : Used to send sync msgs between functions and main driver
type SyncMsg struct{}

// SignUpData : Generic struct to use for signup requests
type SignUpData struct {
	UserID      string `json:"userid"`
	Password    string `json:"password"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Gender      string `json:"gender"`
	Age         string `json:"age"`
	PhoneNumber string `json:"phonenumber"`
}

// LoginData : JSON format data for incoming login request
type LoginData struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

// EventCreationData : JSON format data for event creation request
type EventCreationData struct {
	EventName    string `json:"eventname"`
	Lat          string `json:"latitude"`
	Long         string `json:"longitude"`
	Creationtime string `json:"creation_time"`
	Creatorid    string `json:"creatorid"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	MaxMem       string `json:"max_mem"`
	MinMem       string `json:"min_mem"`
	FriendOnly   string `json:"friend_only"`
	Gender       string `json:"gender"`
	MinAge       string `json:"min_age"`
	MaxAge       string `json:"max_age"`
}

// Events : slice of EventFmt
type Events []models.Event

// SearchEventsData : JSON format data for finding events
type SearchEventsData struct {
	Lat    string `json:"latitude"`
	Long   string `json:"longitude"`
	Radius string `json:"radius"`
}

// SearchResults : Fmt of JSON for returning results
type SearchResults struct {
	Events Events `json:"events"`
}
