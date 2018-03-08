package main

import (
	"log"
	"net/url"
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
	config     string
	username   string
	password   string
	commonName string
	location   string
	authyAPI   authyAPI
}

func (d *authyVPNData) authenticate() bool {
	id, cn, err := getAuthyID(d.config, d.username)
	if err != nil {
		log.Printf("Error getting authy ID from config: %s", err.Error())
		return false
	}

	if cn != "" && d.commonName != cn {
		log.Printf("Error: client common name %s does not match the configuration file common name %s", d.commonName, cn)
		return false
	}

	authyID := strconv.Itoa(id)

	switch d.password {
	case "onetouch":
		log.Printf("Starting OneTouch authorization for user %s with Authy ID %s", d.username, authyID)

		details := authy.Details{
			"User":     d.username,
			"Location": d.location,
		}

		approvalRequest, err := d.authyAPI.SendApprovalRequest(authyID, "OpenVPN login", details, url.Values{"seconds_to_expire": {"60"}})
		if err != nil {
			log.Printf("Error in authy.SendApprovalRequest: %s", err.Error())
			return false
		}

		status, err := d.authyAPI.WaitForApprovalRequest(approvalRequest.UUID, 60*time.Second, url.Values{})
		if err != nil {
			log.Printf("Error in authy.WaitForApprovalRequest: %s", err.Error())
			return false
		}

		if status == authy.OneTouchStatusApproved {
			return true
		}

	case "sms":
		log.Printf("Sending SMS for user %s with Authy ID %s", d.username, authyID)

		sms, err := d.authyAPI.RequestSMS(authyID, url.Values{})
		if err != nil {
			log.Printf("Error in authy.RequestSMS: %s", err.Error())
			return false
		}

		if !sms.Valid() {
			log.Printf("Error: request for SMS was invalid")
			return false
		}

	case "call":
		log.Printf("Calling user %s with Authy ID %s", d.username, authyID)

		call, err := d.authyAPI.RequestPhoneCall(authyID, url.Values{})
		if err != nil {
			log.Printf("Error in authy.RequestPhoneCall: %s", err.Error())
			return false
		}

		if !call.Valid() {
			log.Printf("Error: request for call was invalid")
			return false
		}

	default:
		log.Printf("Verifying token for user %s with Authy ID %s", d.username, authyID)

		verification, err := d.authyAPI.VerifyToken(authyID, d.password, url.Values{})
		if err != nil {
			log.Printf("Error in authy.VerifyToken: %s", err.Error())
			return false
		}
		return verification.Valid()
	}

	return false
}
