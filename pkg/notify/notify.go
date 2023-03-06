package notify

import (
	log "github.com/sirupsen/logrus"
	"jtdx-alarm/pkg/city"
)

type Notifier interface {
	Notify(de string, entry *city.DXCCEntry, msg string)
}

var NotifiersMap map[string]func() Notifier

type LogPrintNotifier struct {
}

func (n LogPrintNotifier) Notify(de string, entry *city.DXCCEntry, msg string) {
	log.Infoln(msg)
}

func init() {
	NotifiersMap = map[string]func() Notifier{
		"log": func() Notifier {
			return LogPrintNotifier{}
		},
		"wx": func() Notifier {
			return NewQYWXMessageNotifier(false)
		},
		"wx-debug": func() Notifier {
			return NewQYWXMessageNotifier(true)
		},
	}
}
