package api

// BasicResponse : JSON reply for API Calls
type BasicResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
