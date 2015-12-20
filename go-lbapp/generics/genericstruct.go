package generics

// SyncMsg : Used to send sync msgs between functions and main driver
type SyncMsg struct{}

// SignUpData : Generic struct to use for signup requests
type SignUpData struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginData : JSON format data for incoming login request
type LoginData struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}
