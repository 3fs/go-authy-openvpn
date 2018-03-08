package main

import (
	"flag"
	"log"
	"os"

	"github.com/dcu/go-authy"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[Authy] ")

	apiKey := flag.String("a", "", "Authy API key")
	config := flag.String("c", "/etc/openvpn/authy/authy-vpn.conf", "Authy config file")
	flag.Parse()

	data := authyVPNData{
		username:   os.Getenv("username"),
		password:   os.Getenv("password"),
		commonName: os.Getenv("common_name"),
		config:     *config,
	}
	controlFile := os.Getenv("auth_control_file")

	if *apiKey == "" {
		log.Println("Error: API key is required")
		writeStatus(false, data.username, controlFile)
		return
	}

	data.authyAPI = authy.NewAuthyAPI(*apiKey)
	success := data.authenticate()

	writeStatus(success, data.username, controlFile)
}

func writeStatus(success bool, username, controlFile string) {
	file, err := os.OpenFile(controlFile, os.O_RDWR, 0755)
	if err != nil {
		log.Printf("Error opening control file: %s\n", err.Error())
		return
	}
	defer file.Close()

	if success {
		log.Printf("Authorization was successful for user %s\n", username)
		file.WriteString("1")
	} else {
		log.Printf("Authorization WAS NOT successful for user %s\n", username)
		file.WriteString("0")
	}
}
