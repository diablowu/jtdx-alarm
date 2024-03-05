package hamalert

import (
	"github.com/gin-gonic/gin"
)

func StartServer(bind string) {
	r := gin.Default()
	r.POST("/alert", alertHandle)

	r.Run(bind)
}

func alertHandle(c *gin.Context) {

}

type HamAlertMessage struct {
	FullCallsign     string `json:"fullCallsign"`
	Callsign         string `json:"callsign"`
	Frequency        string `json:"frequency"`
	Band             string `json:"band"`
	Mode             string `json:"mode"`
	ModeDetail       string `json:"modeDetail"`
	Time             string `json:"time"`
	Dxcc             string `json:"dxcc"`
	HomeDxcc         string `json:"homeDxcc"`
	SpotterDxcc      string `json:"spotterDxcc"`
	CQZone           string `json:"cq"`
	Continent        string `json:"continent"`
	Entity           string `json:"entity"`
	HomeEntity       string `json:"homeEntity"`
	SpotterEntity    string `json:"spotterEntity"`
	Spotter          string `json:"spotter"`
	SpotterCq        string `json:"spotterCq"`
	SpotterContinent string `json:"spotterContinent"`
	RawText          string `json:"rawText"`
	Title            string `json:"title"`
	Comment          string `json:"comment"`
	Source           string `json:"source"`
	Speed            string `json:"speed"`
	Snr              string `json:"snr"`
	TriggerComment   string `json:"triggerComment"`
	Qsl              string `json:"qsl"`
	State            string `json:"state"`
	SpotterState     string `json:"spotterState"`
	IotaGroupRef     string `json:"iotaGroupRef"`
	IotaGroupName    string `json:"iotaGroupName"`
	SummitName       string `json:"summitName"`
	SummitHeight     string `json:"summitHeight"`
	SummitPoints     string `json:"summitPoints"`
	SummitRef        string `json:"summitRef"`
	WwffName         string `json:"wwffName"`
	WwffDivision     string `json:"wwffDivision"`
	WwffRef          string `json:"wwffRef"`
}
