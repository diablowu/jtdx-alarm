package qywx

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AccessToken struct {
	Code      int    `json:"errcode"`
	Message   string `json:"errmsg"`
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

const (
	QYAPIEndpoint = "https://qyapi.weixin.qq.com"
	TokenFileURL  = "http://js.tkpdevops.com/k.txt"
)

var accessToken *string

func GetAccessToken() string {
	return *accessToken
}

func freshAccessToken() {
	resp, err := http.Get(TokenFileURL)
	if err == nil {
		if bs, err := ioutil.ReadAll(resp.Body); err == nil {
			tokenString := strings.TrimSpace(string(bs))
			accessToken = &tokenString
		}
	}
}

func FreshTokenTask(interval time.Duration) {

	ticker := time.NewTicker(interval)
	freshAccessToken()
	log.Println("Success to fresh access token")
	go func() {
		for {
			<-ticker.C
			log.Println("Begin to refresh access token")
			freshAccessToken()
			log.Printf("New access token is %s", *accessToken)
		}
	}()
}
