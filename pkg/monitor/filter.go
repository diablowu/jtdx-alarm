package monitor

import (
	"jtdx-alarm/pkg/callsign"
	"strings"
)

type MessageFilter interface {
	Filter(de *callsign.CallSign) bool
}

type CompositeFilter struct {
	filters []MessageFilter
}

func (f CompositeFilter) Filter(de *callsign.CallSign) bool {
	ret := false

	for _, f := range f.filters {
		ret = ret && f.Filter(de)
	}

	return ret
}

func NewCompositeFilter(filter ...MessageFilter) CompositeFilter {
	return CompositeFilter{filters: filter}
}

type DXCCFilter struct {
	dxcc map[string]int
}

func NewDXCCFilter(dxcc []string) *DXCCFilter {
	dm := make(map[string]int, len(dxcc))
	for _, dxccName := range dxcc {
		dm[strings.TrimSpace(dxccName)] = 0
	}
	return &DXCCFilter{dxcc: dm}
}

func (f DXCCFilter) Filter(de *callsign.CallSign) bool {
	_, found := f.dxcc[de.DXCC.DXCCName]
	return found
}
