package main

import (
	"errors"
	"log"
	"os"

	"github.com/dcu/go-authy"
)

func logError(err error) bool {
	log.Println("Error:", err)
	return false
}

func main() {
	data := authyVPNData{
		username:    os.Getenv("username"),
		password:    os.Getenv("password"),
		controlFile: os.Getenv("auth_control_file"),
	}

	log.SetFlags(0)
	log.SetPrefix("[Authy] ")

	if len(os.Args) == 1 {
		status := logError(errors.New("First argument (API key) is required"))
		data.writeStatus(&status)
		return
	}
	data.authyAPI = authy.NewAuthyAPI(os.Args[1])

	if len(os.Args) == 2 {
		data.config = "/etc/openvpn/authy/authy-vpn.conf"
		log.Printf("Using default config (%s)", data.config)
	} else {
		data.config = os.Args[2]
	}

	status := false
	defer data.writeStatus(&status)

	status = data.authenticate()
}
