package decode

import (
	"jtdx-alarm/pkg/callsign"
	"strings"
)

type BlacklistDXCCCallSignFilter struct {
	blacklist map[string]int
}

func NewBlacklistDXCCCallSignFilter(dxcc []string) *BlacklistDXCCCallSignFilter {
	dm := make(map[string]int, len(dxcc))
	for _, dxccName := range dxcc {
		dm[strings.TrimSpace(dxccName)] = 0
	}
	return &BlacklistDXCCCallSignFilter{blacklist: dm}
}

func (f BlacklistDXCCCallSignFilter) Filter(de *callsign.CallSign) bool {
	_, found := f.blacklist[de.DXCC.DXCCName]
	return found
}
