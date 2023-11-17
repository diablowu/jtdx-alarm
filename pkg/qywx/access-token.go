package qywx

import (
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	QYAPIEndpoint = "https://qyapi.weixin.qq.com"
)

type AccessToken struct {
	Code      int    `json:"errcode"`
	Message   string `json:"errmsg"`
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

var accessToken *string

func GetAccessToken() string {
	return *accessToken
}

func setCurrentAccessToken(token string) {
	accessToken = &token
}

func FreshTokenTask(agent int, corpID, secret string, interval time.Duration) {
	log.Infoln("Begin to fresh access token task")
	ticker := time.NewTicker(interval)
	freshAccessToken(agent, corpID, secret)
	log.Infoln("New access token is %s", *accessToken)
	go func() {
		for {
			<-ticker.C
			log.Infoln("Begin to refresh access token")
			freshAccessToken(agent, corpID, secret)
			log.Infof("New access token is %s", *accessToken)
		}
	}()
}
