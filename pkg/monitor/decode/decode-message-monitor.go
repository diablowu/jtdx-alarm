package decode

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"jtdx-alarm/pkg/callsign"
	"jtdx-alarm/pkg/city"
	"jtdx-alarm/pkg/notify"
	"jtdx-alarm/pkg/wsjtx"
	"strings"
)

type DecodeMessageMonitors []DecodeMessageMonitor

func (dms DecodeMessageMonitors) Do(msg wsjtx.DecodeMessage) {
	for _, m := range []DecodeMessageMonitor(dms) {
		m.Monit(msg)
	}
}

func CreateDecodeMessageMonitors(monitors ...DecodeMessageMonitor) *DecodeMessageMonitors {
	var ma []DecodeMessageMonitor
	for _, mo := range monitors {
		ma = append(ma, mo)
	}
	mds := DecodeMessageMonitors(ma)
	return &mds
}

type DecodeMessageMonitor interface {
	Monit(msg wsjtx.DecodeMessage)
}

type DeCallSignMonitor struct {
	filter    CallSignFilter
	notifiers []notify.Notifier
}

func NewDefaultMonitor(filter CallSignFilter, notifiers []string) DecodeMessageMonitor {

	log.Infof("Regiestered notifiers: %v", notify.NotifiersMap)
	log.Infof("Requested notifiers: %v", notifiers)
	var ns []notify.Notifier
	for _, name := range notifiers {
		if n, found := notify.NotifiersMap[strings.TrimSpace(name)]; found {
			ns = append(ns, n())
		}
	}

	log.Infof("Notifiers: %v was be enabled", ns)

	return DeCallSignMonitor{
		filter:    filter,
		notifiers: ns,
	}
}

func (monitor DeCallSignMonitor) Monit(msg wsjtx.DecodeMessage) {
	if de, _, err := callsign.ExtractCallSignFromMessage(msg.Message, true); err == nil {
		deDXCC := city.FindDXCC(de.Number)
		if !monitor.filter.Filter(de) {
			for _, n := range monitor.notifiers {
				n.Notify(de.Number, deDXCC, fmt.Sprintf("%s %s(%s)", de.Number, deDXCC.City, deDXCC.DXCCName))
			}
		}

	} else {
		log.Warnf("Failed to ExtractCallSignFromMessage: %s", err)
	}
}
