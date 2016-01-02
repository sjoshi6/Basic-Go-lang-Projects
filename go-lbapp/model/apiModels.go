package models

// User : This is used only for the API response as a model
type User struct {
	UserID      string `json:"userid"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Gender      string `json:"gender"`
	Age         string `json:"age"`
	PhoneNumber string `json:"phonenumber"`
}

// Event : to wrap data for returning
type Event struct {
	ID           string `json:"id"`
	EventName    string `json:"eventname"`
	Lat          string `json:"latitude"`
	Long         string `json:"longitude"`
	Creationtime string `json:"creationtime"`
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
