package notify

import (
	"jtdx-alarm/pkg/city"
	"log"
)

type Notifier interface {
	Notify(de string, entry *city.DXCCEntry, msg string)
}

var NotifiersMap map[string]func() Notifier

type LogPrintNotifier struct {
}

func (n LogPrintNotifier) Notify(de string, entry *city.DXCCEntry, msg string) {
	log.Printf(msg)
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
