package adif

import (
	goadiflog "github.com/Eminlin/GoADIFLog"
	"jtdx-alarm/pkg/city"
	"log"
	"time"
)

type ADIFCache map[string]struct{}

var defaultADIFCache *ADIFCache

func (ac ADIFCache) addDXCC(dxcc string) {
	ac[dxcc] = struct{}{}
}

func (ac ADIFCache) doCheck(dxcc string) (exist bool) {
	_, exist = ac[dxcc]
	return
}

func DoCheck(dxcc string) bool {
	return defaultADIFCache.doCheck(dxcc)
}

func InitLoggerChecker(adifPath string, refreshInterval time.Duration) {

	if cache, err := loadADIF(adifPath); err != nil {
		log.Fatalf("Failed to read adif: %s", err)
	} else {
		defaultADIFCache = cache
	}

	go func() {
		to := time.NewTimer(refreshInterval)
		<-to.C
		refreshTicker := time.NewTicker(refreshInterval)

		for {
			select {
			case <-refreshTicker.C:
				{
					if cache, err := loadADIF(adifPath); err != nil {
						log.Printf("Failed to read adif, %s", err)
					} else {
						defaultADIFCache = cache
					}
				}
			}
		}
	}()
}

func loadADIF(adifPath string) (*ADIFCache, error) {

	if logContent, err := goadiflog.Parse(adifPath); err != nil {
		return nil, err
	} else {
		cache := ADIFCache{}
		for _, one := range logContent {
			dxcc := city.FindDXCC(one.Call)
			//log.Printf("call:%s, dxcc:%s", one.Call, dxcc.DXCCName)
			cache.addDXCC(dxcc.DXCCName)
		}

		log.Printf("ADIF loaded. DXCC count: %d", len(cache))
		log.Printf("ADIF loaded. DXCC list: %v", cache)

		return &cache, nil
	}
}
