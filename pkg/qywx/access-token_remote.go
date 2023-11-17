//go:build !local

package qywx

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

const (
	TokenFileURL = "http://js.tkpdevops.com/k.txt"
)

func freshAccessToken(agent int, corpID, secret string) {
	log.Infoln("Begin to refresh access token via remote mode")
	resp, err := http.Get(TokenFileURL)
	if err == nil {
		if bs, err := io.ReadAll(resp.Body); err == nil {
			tokenString := strings.TrimSpace(string(bs))
			accessToken = &tokenString
		}
	}
}
