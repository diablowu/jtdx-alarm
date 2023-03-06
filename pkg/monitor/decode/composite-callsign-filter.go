package decode

import "jtdx-alarm/pkg/callsign"

type CompositeFilter struct {
	filters []CallSignFilter
}

func (f CompositeFilter) Filter(de *callsign.CallSign) bool {
	for _, f := range f.filters {
		if f.Filter(de) {
			return true
		}
	}
	return false
}

func NewCompositeFilter(filter ...CallSignFilter) CompositeFilter {
	return CompositeFilter{filters: filter}
}
