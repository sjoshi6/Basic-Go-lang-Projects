package settings

/*
   Contains all common configuration parameters for the project
*/
const (

	// ServerPort : common setting for all storage nodes to accept incoming rpc calls
	ServerPort       string = "9376"
	DBName           string = "storagenode"
	DBHostName       string = "localhost"
	DBPort           string = "3306"
	RedisPORT        string = "6379"
	LoadBalancerPort string = "5000"
)
