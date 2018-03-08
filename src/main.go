package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dcu/go-authy"
	geoip2 "github.com/oschwald/geoip2-golang"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[Authy] ")

	apiKey := flag.String("a", "", "Authy API key")
	config := flag.String("c", "/etc/openvpn/authy/authy-vpn.conf", "Authy config file")
	geoipDBpath := flag.String("g", "", "MaxMind GeoLite2 DB path")
	flag.Parse()

	data := authyVPNData{
		username:   os.Getenv("username"),
		password:   os.Getenv("password"),
		commonName: os.Getenv("common_name"),
		location:   os.Getenv("untrusted_ip"),
		config:     *config,
	}
	controlFile := os.Getenv("auth_control_file")

	if *apiKey == "" {
		log.Println("Error: API key is required")
		writeStatus(false, data.username, controlFile)
		return
	}

	if *geoipDBpath != "" {
		location, _ := getLocation(*geoipDBpath, data.location)
		if location != "" {
			data.location = fmt.Sprintf("%s (%s)", location, data.location)
		}
	}

	data.authyAPI = authy.NewAuthyAPI(*apiKey)
	success := data.authenticate()

	writeStatus(success, data.username, controlFile)
}

func getLocation(geoDB, ip string) (string, error) {
	db, err := geoip2.Open(geoDB)
	if err != nil {
		log.Printf("Error opening GeoIP database: %s\n", err.Error())
		return "", err
	}
	defer db.Close()

	record, err := db.City(net.ParseIP(ip))
	if err != nil {
		log.Printf("Error geting city from IP: %s\n", err.Error())
		return "", nil
	}

	if record.City.Names["en"] == "" {
		return record.Country.Names["en"], nil
	}

	return fmt.Sprintf("%s, %s", record.City.Names["en"], record.Country.Names["en"]), nil
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
