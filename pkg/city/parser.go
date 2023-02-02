package city

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var PrefixData = make(map[string]*DXCCEntry)
var CallSignCache = make(map[string]*DXCCEntry)

type DXCCEntry struct {
	City          string
	CQZone        string
	ITUZone       string
	Continent     string
	NorthLatitude string
	WestLongitude string
	UTCOffset     string
	DXCCName      string
}

var unknownDXCC = &DXCCEntry{
	City:          "?",
	CQZone:        "?",
	ITUZone:       "?",
	Continent:     "?",
	NorthLatitude: "?",
	WestLongitude: "?",
	UTCOffset:     "?",
	DXCCName:      "?",
}

func (entry DXCCEntry) String() string {
	return fmt.Sprintf("city:%s,cqzone:%s,ituzone:%s,continent:%s,lat:%s,lang:%s,offset:%s,dxcc-name:%s",
		entry.City, entry.CQZone, entry.ITUZone, entry.Continent, entry.NorthLatitude, entry.WestLongitude, entry.UTCOffset, entry.DXCCName)
}

func FindDXCC(call string) *DXCCEntry {
	if cached, found := CallSignCache[call]; found {
		return cached
	}

	var foundDXCC *DXCCEntry

	if dxcc, found := PrefixData["="+call]; found {
		foundDXCC = dxcc
	} else {
		slashIdx := strings.Index(call, "/")
		if slashIdx < 0 {
			foundDXCC = findDXCCRecursion(call)
		} else {
			foundDXCC = findDXCCRecursion(call[:slashIdx])
		}
	}

	CallSignCache[call] = foundDXCC

	return foundDXCC
}

func findDXCCRecursion(call string) *DXCCEntry {
	callLen := len(call)
	if callLen == 0 {
		return unknownDXCC
	}
	searchPrefix := call[:callLen-1]
	if dxcc, found := PrefixData[searchPrefix]; found {
		return dxcc
	}
	return findDXCCRecursion(searchPrefix)
}

func LoadFromCTYData(ctyDataFilePath string) error {
	PrefixData = make(map[string]*DXCCEntry)

	if fd, err := os.Open(ctyDataFilePath); err == nil && fd != nil {

		scanner := bufio.NewScanner(fd)

		cityDeclareBegin := false
		var currentDXCCEntry *DXCCEntry
		// China:                    24:  44:  AS:   36.00:  -102.00:    -8.0:  BY:
		// /t prefix list
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if !cityDeclareBegin {
				currentDXCCEntry = extractDXCCEntryFromLine(line)
				cityDeclareBegin = true
				continue
			}

			// process prefix list
			if cityDeclareBegin {
				prefixList := strings.Split(line, ",")
				for _, prefix := range prefixList {
					pfx := strings.TrimSpace(prefix)
					if pfx != "" {
						PrefixData[normalizePrefix(pfx)] = currentDXCCEntry
					}
				}

				// prefix list end and this dxcc declared end
				if strings.HasSuffix(line, ";") {
					cityDeclareBegin = false
					currentDXCCEntry = nil
				}
				continue
			}
		}
		return fd.Close()
	} else {
		return err
	}

	return nil
}

// China:                    24:  44:  AS:   36.00:  -102.00:    -8.0:  BY:
func extractDXCCEntryFromLine(line string) *DXCCEntry {
	parts := strings.Split(strings.TrimSpace(line), ":")
	return &DXCCEntry{
		City:          strings.TrimSpace(parts[0]),
		CQZone:        strings.TrimSpace(parts[1]),
		ITUZone:       strings.TrimSpace(parts[2]),
		Continent:     strings.TrimSpace(parts[3]),
		NorthLatitude: strings.TrimSpace(parts[4]),
		WestLongitude: strings.TrimSpace(parts[5]),
		UTCOffset:     strings.TrimSpace(parts[6]),
		DXCCName:      strings.TrimSpace(parts[7]),
	}
}

/* Following rules:

(#)	Override CQ Zone
[#]	Override ITU Zone
<#/#>	Override latitude/longitude
{aa}	Override Continent
~#~	Override local time offset from GMT

end with ;
*/
func normalizePrefix(prefix string) string {
	rs := removePairTag(prefix, "(", ")")
	rs = removePairTag(rs, "[", "]")
	rs = removePairTag(rs, "<", ">")
	rs = removePairTag(rs, "{", "}")
	rs = removePairTag(rs, "~", "~")
	rs = strings.TrimSuffix(rs, ";")
	return rs
}

// prefix<BEGIN>34<END>
func removePairTag(raw, beginTag, endTag string) string {
	beginIdx := strings.Index(raw, beginTag)

	if beginIdx <= 0 {
		return raw
	}
	pre := raw[0:beginIdx]
	rest := raw[beginIdx+1:]

	endIdx := strings.Index(rest, endTag)
	if beginIdx <= 0 {
		return raw
	}
	return pre + rest[endIdx+1:]
}
