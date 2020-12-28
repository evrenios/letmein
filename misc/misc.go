package misc

import "os"

var (
	// Secret is the common secret that server-client uses
	Secret = getSecret()
)

// AuthReq holds the properties for a common new authorization request
type AuthReq struct {
	IP        string
	Secret    string
	Hour      int
	Name      string
	Timestamp int64
}

func getSecret() string {
	secret, ok := os.LookupEnv("LETMEIN_SECRET")
	if ok {
		return secret
	}
	return "superdupersecret"
}
