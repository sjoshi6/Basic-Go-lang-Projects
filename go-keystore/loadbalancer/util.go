package loadbalancer

import (
	"encoding/json"
	"net"
	"net/http"
)

// BalancerCommunication : Reply sent by load balancer
type BalancerCommunication struct {
	IPAddr string `json:"ip_address"`
}

// BasicResponse : JSON reply for API Calls
type BasicResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// GetIPAddress : Used to get IP of the running machine
func GetIPAddress() string {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
			// process IP address
		}
	}

	return ""
}

// RespondOrThrowErr : Respond to general requests or exit with server err.
func RespondOrThrowErr(responseObj BasicResponse, w http.ResponseWriter) {

	responseJSON, err := json.Marshal(responseObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// ThrowInternalErrAndExit : Respond with internal server error
func ThrowInternalErrAndExit(w http.ResponseWriter) {

	responsecontent := BasicResponse{
		"Internal Server Error",
		500,
	}

	w.WriteHeader(http.StatusInternalServerError)
	RespondOrThrowErr(responsecontent, w)
}

// RespondSuccessAndExit : Repond with a success
func RespondSuccessAndExit(w http.ResponseWriter, msg string) {

	responsecontent := BasicResponse{
		msg,
		200,
	}
	w.WriteHeader(http.StatusOK)
	RespondOrThrowErr(responsecontent, w)

}
