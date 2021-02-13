package main

import (
	"net/http"
	"net/url"
	"strings"
	"os"
)

func text() {
	accountSid := os.Getenv("ACCOUNT_SID")
	token := os.Getenv("TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", "9254467645")
	msgData.Set("From", os.Getenv("PHONE"))
	msgData.Set("Body", "HEHE UR CUTE!")
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(accountSid, token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
}

