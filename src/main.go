package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dcu/go-authy"
)

type vpnData struct {
	username    string
	password    string
	controlFile string
}

func (d *vpnData) writeStatus(success bool) {
	file, err := os.OpenFile(d.controlFile, os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if success {
		fmt.Printf("[Authy] Authorization was successful for user %s\n", d.username)
		file.WriteString("1")
	} else {
		fmt.Printf("[Authy] Authorization wasn't successful for user %s\n", d.username)
		file.WriteString("0")
	}
}

func main() {
	fmt.Printf("%v\n", os.Args)
	apiKey := os.Args[1]
	config := os.Args[2]

	data := vpnData{
		username:    os.Getenv("username"),
		password:    os.Getenv("password"),
		controlFile: os.Getenv("auth_control_file"),
	}

	fmt.Printf("[Authy] Starting with API key %s and config %s for user %s\n", apiKey, config, data.username)

	id, err := getAuthyID(config, data.username)
	if err != nil {
		data.writeStatus(false)
		log.Fatal(err)
	}

	authyID := strconv.Itoa(id)
	authyAPI := authy.NewAuthyAPI(apiKey)

	switch data.password {
	case "onetouch":
		details := authy.Details{
			"User":       data.username,
			"IP Address": os.Getenv("untrusted_ip"),
		}
		approvalRequest, err := authyAPI.SendApprovalRequest(authyID, "OpenVPN login", details, url.Values{"seconds_to_expire": {"60"}})
		if err != nil {
			data.writeStatus(false)
			log.Fatal(err)
		}

		status, err := authyAPI.WaitForApprovalRequest(approvalRequest.UUID, 60*time.Second, url.Values{})
		if err != nil {
			data.writeStatus(false)
			log.Fatal(err)
		}

		if status == authy.OneTouchStatusApproved {
			data.writeStatus(true)
		} else {
			data.writeStatus(false)
		}

	case "sms":
		sms, err := authyAPI.RequestSMS(authyID, url.Values{})
		data.writeStatus(false) // making SMS request will always fail authorizatoin

		if err != nil {
			log.Fatal(err)
		}
		if !sms.Valid() {
			log.Fatal("[Authy] Request for SMS failed")
		}

	case "call":
		call, err := authyAPI.RequestPhoneCall(authyID, url.Values{})
		data.writeStatus(false) // making call request will always fail authorizatoin

		if err != nil {
			log.Fatal(err)
		}
		if !call.Valid() {
			log.Fatal("[Authy] Request for call failed")
		}

	default:
		verification, err := authyAPI.VerifyToken(authyID, data.password, url.Values{})
		if err != nil {
			data.writeStatus(false)
			log.Fatal(err)
		}
		data.writeStatus(verification.Valid())
	}
}
