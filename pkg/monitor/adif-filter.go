package monitor

import (
	"jtdx-alarm/pkg/adif"
	"jtdx-alarm/pkg/callsign"
)

type ADIFFilter struct {
}

func (f ADIFFilter) Filter(de *callsign.CallSign) bool {
	return adif.DoCheck(de.DXCC.DXCCName)
}

func NewADIFFilter() ADIFFilter {
	return ADIFFilter{}
}
