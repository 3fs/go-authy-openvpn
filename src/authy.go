package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dcu/go-authy"
)

type authyAPI interface {
	SendApprovalRequest(userID string, message string, details authy.Details, params url.Values) (*authy.ApprovalRequest, error)
	WaitForApprovalRequest(uuid string, maxDuration time.Duration, params url.Values) (authy.OneTouchStatus, error)
	VerifyToken(userID string, token string, params url.Values) (*authy.TokenVerification, error)
	RequestSMS(userID string, params url.Values) (*authy.SMSRequest, error)
	RequestPhoneCall(userID string, params url.Values) (*authy.PhoneCallRequest, error)
}

type authyVPNData struct {
	config      string
	username    string
	password    string
	controlFile string
	authyAPI    authyAPI
}

func (d *authyVPNData) writeStatus(success *bool) {
	file, err := os.OpenFile(d.controlFile, os.O_RDWR, 0755)
	if err != nil {
		logError(err)
	}
	defer file.Close()

	if *success {
		log.Printf("Authorization was successful for user %s\n", d.username)
		file.WriteString("1")
	} else {
		log.Printf("Authorization wasn't successful for user %s\n", d.username)
		file.WriteString("0")
	}
}

func (d *authyVPNData) authenticate() bool {
	id, err := getAuthyID(d.config, d.username)
	if err != nil {
		return logError(err)
	}

	authyID := strconv.Itoa(id)

	switch d.password {
	case "onetouch":
		log.Printf("Starting OneTouch authorization for user %s with Authy ID %s", d.username, authyID)

		details := authy.Details{
			"User":       d.username,
			"IP Address": os.Getenv("untrusted_ip"),
		}
		approvalRequest, err := d.authyAPI.SendApprovalRequest(authyID, "OpenVPN login", details, url.Values{"seconds_to_expire": {"60"}})
		if err != nil {
			return logError(err)
		}

		status, err := d.authyAPI.WaitForApprovalRequest(approvalRequest.UUID, 60*time.Second, url.Values{})
		if err != nil {
			return logError(err)
		}

		if status == authy.OneTouchStatusApproved {
			return true
		}

	case "sms":
		log.Printf("Sending SMS for user %s with Authy ID %s", d.username, authyID)

		sms, err := d.authyAPI.RequestSMS(authyID, url.Values{})

		if err != nil {
			return logError(err)
		}
		if !sms.Valid() {
			return logError(errors.New("Request for SMS failed"))
		}

	case "call":
		log.Printf("Calling user %s with Authy ID %s", d.username, authyID)

		call, err := d.authyAPI.RequestPhoneCall(authyID, url.Values{})

		if err != nil {
			return logError(err)
		}
		if !call.Valid() {
			return logError(errors.New("Request for call failed"))
		}

	default:
		log.Printf("Verifying token for user %s with Authy ID %s", d.username, authyID)

		verification, err := d.authyAPI.VerifyToken(authyID, d.password, url.Values{})
		if err != nil {
			return logError(err)
		}
		return verification.Valid()
	}

	return false
}
