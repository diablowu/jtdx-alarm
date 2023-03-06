package decode

import (
	"jtdx-alarm/pkg/callsign"
)

type CallSignFilter interface {
	Filter(de *callsign.CallSign) bool
}
