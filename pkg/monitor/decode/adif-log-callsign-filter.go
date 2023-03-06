package decode

import (
	log "github.com/sirupsen/logrus"
	"jtdx-alarm/pkg/adif"
	"jtdx-alarm/pkg/callsign"
)

type ADIFLogCallSignFilter struct {
}

func (f ADIFLogCallSignFilter) Filter(de *callsign.CallSign) bool {

	found := adif.DoCheck(de.DXCC.DXCCName)
	if found {
		log.Infof("ADIFFilter: %s(%s) found in cache, should be filtered", de.DXCC.DXCCName, de.Number)
	}

	return found
}

func NewADIFFilter() ADIFLogCallSignFilter {
	return ADIFLogCallSignFilter{}
}
