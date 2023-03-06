package notify

import (
	log "github.com/sirupsen/logrus"
	"jtdx-alarm/pkg/city"
	"jtdx-alarm/pkg/qywx"
	"strings"
	"time"
)

var histories = make(map[string]int64)

const sendIntervalInSecond = 60

type QYWXMessageNotifier struct {
	AgentID      int
	messageQueue chan messageEntry
	dryRun       bool
}

type messageEntry struct {
	de      string
	message string
}

func NewQYWXMessageNotifier(dryRun bool) Notifier {

	n := QYWXMessageNotifier{
		messageQueue: make(chan messageEntry, 500),
		dryRun:       dryRun,
	}

	go func() {
		n.startSendTask()
	}()
	return &n
}

func (n *QYWXMessageNotifier) startSendTask() {

	ticker := time.NewTicker(time.Second * 30)
	for {

		<-ticker.C
		log.Infof("Message queue: len(%d), cap(%d)", len(n.messageQueue), cap(n.messageQueue))
		var contents []string

		for i := 0; i < 25; i++ {
			select {
			case <-time.After(time.Second * 5):
				{
					break
				}
			case m := <-n.messageQueue:
				{
					contents = append(contents, m.message)
				}
			}
		}
		if !n.dryRun {
			if len(contents) > 0 {
				qywx.SendAgentMessage("Found DXCC\n\n" + strings.Join(contents, "\n"))
			}
		} else {
			log.Infof("Following message was sent via qywx: %v", contents)
		}

	}
}

func (n *QYWXMessageNotifier) Notify(de string, entry *city.DXCCEntry, msg string) {

	now := time.Now().Unix()
	canBeEnQueue := false
	if lastTS, found := histories[de]; found {
		interval := now - lastTS
		if interval >= sendIntervalInSecond {
			canBeEnQueue = true
			log.Infof("Call[%s] was be enqueue , interval :%d", de, interval)
		} else {
			log.Infof("Call[%s] was be blocked , interval :%d", de, interval)
		}
	} else {
		canBeEnQueue = true
		log.Infof("Call[%s] was be enqueue , not found in cache", de)
	}

	if canBeEnQueue {
		histories[de] = now
		n.messageQueue <- messageEntry{
			de:      de,
			message: msg,
		}
	}
}
