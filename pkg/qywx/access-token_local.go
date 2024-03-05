//go:build local
// +build local

package qywx

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func freshAccessToken(agent int, corpID, secret string) (error, *string) {
	log.Infoln("Begin to refresh access token via local mode directly")
	qywxAccessTokenAPI := fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		QYAPIEndpoint,
		corpID,
		secret)

	// 向qywxAccessTokenAPI发器get请求，然后读取json数据，并反序列化为AccessToken对象

	if resp, err := http.Get(qywxAccessTokenAPI); err == nil && resp.StatusCode == http.StatusOK {
		if bs, err := io.ReadAll(resp.Body); err == nil {
			accessToken := &AccessToken{}
			if err := json.Unmarshal(bs, accessToken); err == nil {
				if accessToken.Code == 0 {
					setCurrentAccessToken(accessToken.Token)
				} else {
					log.Warnf("Failed to get access token, %s", accessToken.Message)
				}
			} else {
				log.Warnf("Failed to unmarshal access token, %s", err)
			}
		}
	} else {
		log.Warnf("Failed to get access token, status:, error:%s", resp.Status, err)
	}

}
