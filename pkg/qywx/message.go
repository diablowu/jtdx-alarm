package qywx

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Service struct {
	call string
}

var targetCallSign string
var agentID int

func Setup(agent int, call string) {
	targetCallSign = call
	agentID = agent
	FreshTokenTask(time.Minute * 15)
}

func SendAgentMessage(message string) {

	url := QYAPIEndpoint + "/cgi-bin/message/send?access_token=" + GetAccessToken()

	tm := TextMessage{
		To:      targetCallSign,
		Type:    "text",
		AgentID: agentID,
		Text: struct {
			Content string `json:"content"`
		}{
			Content: message,
		},
	}

	if bs, err := json.Marshal(tm); err == nil {
		if resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs)); err == nil {
			bs, _ := ioutil.ReadAll(resp.Body)
			log.Println(string(bs))
		} else {
			log.Warnf("Failed to send agent message, %s", err)
		}

	}

}

type TextMessage struct {
	To      string `json:"touser"`
	Type    string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}
