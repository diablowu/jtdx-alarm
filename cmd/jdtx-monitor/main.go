package main

import (
	"flag"
	wsjtx "jtdx-alarm"
	"jtdx-alarm/pkg/adif"
	"jtdx-alarm/pkg/city"
	"jtdx-alarm/pkg/monitor"
	"jtdx-alarm/pkg/osx"
	"jtdx-alarm/pkg/qywx"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	bindAddr       = flag.String("bind-addr", "239.255.0.0", "Bind address or Multicast address")
	bindPort       = flag.Uint("bind-port", 2237, "Bind port")
	ctyPath        = flag.String("cty-path", "cty.dat", "CTY file")
	verbose        = flag.Bool("verbose", false, "Verbose mode")
	targetCallSign = flag.String("target-call", "N0CALL", "Callsign of received message")
	filteredDXCC   = flag.String("filtered-dxcc", "BY,JA,HL,BV,YB,UA,UA9", "Filtered DXCC")
	notifiers      = flag.String("notifiers", "log,wx", "Notifier list")
	jtdxLogDir     = flag.String("jtdx-log-dir", filepath.Join(osx.MustUserHomeDir(), "AppData", "Local", "JTDX"), "JTDX logger path")
	useADIFFilter  = flag.Bool("use-adif-filter", true, "Use adif")
)

var (
	DefaultAgentID             = 1000002
	DefaultADIFLogFileName     = "wsjtx_log.adi"
	DefaultADIFRefreshInterval = time.Second * 5
)

var defaultDecodeMessageMonitors *monitor.DecodeMessageMonitors

// Simple driver binary for wsjtx-go library.
func main() {
	initCliFlags()

	qywx.Setup(DefaultAgentID, strings.ToLower(*targetCallSign))
	initBigCTY()

	if *useADIFFilter {
		adif.InitLoggerChecker(filepath.Join(*jtdxLogDir, DefaultADIFLogFileName), DefaultADIFRefreshInterval)
	}

	log.Println("Listening for JTDX...")
	incomingMessageChannel := make(chan interface{}, 5)

	var finalFilter monitor.MessageFilter
	if *useADIFFilter {
		finalFilter = monitor.NewCompositeFilter(monitor.NewDXCCFilter(strings.Split(*filteredDXCC, ",")), monitor.NewADIFFilter())
	} else {
		finalFilter = monitor.NewDXCCFilter(strings.Split(*filteredDXCC, ","))
	}

	defaultDecodeMessageMonitors = monitor.CreateDecodeMessageMonitors(
		monitor.NewDefaultMonitor(finalFilter, strings.Split(*notifiers, ",")))

	wsjtxServer, err := wsjtx.MakeServer(*bindAddr, *bindPort)
	if err != nil {
		log.Fatalf("%v", err)
	}

	errChannel := make(chan error, 5)
	go wsjtxServer.ListenToWsjtx(incomingMessageChannel, errChannel)

	for {
		select {
		case err := <-errChannel:
			log.Printf("error: %v", err)
		case message := <-incomingMessageChannel:
			handleServerMessage(message)
		}
	}
}

func initBigCTY() {
	if err := city.LoadFromCTYData(*ctyPath); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.Println("Success to load cty data")
	}
}

func initCliFlags() {
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}
}

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
	switch message.(type) {
	case wsjtx.HeartbeatMessage:
		log.Println("Heartbeat:", message)
	case wsjtx.StatusMessage:
		log.Println("Other:", reflect.TypeOf(message), message)
	case wsjtx.DecodeMessage:
		defaultDecodeMessageMonitors.Do(message.(wsjtx.DecodeMessage))
	case wsjtx.ClearMessage:
		log.Println("Clear:", message)
	case wsjtx.QsoLoggedMessage:
		log.Println("QSO Logged:", message)
	case wsjtx.CloseMessage:
		log.Println("Close:", message)
	case wsjtx.WSPRDecodeMessage:
		log.Println("WSPR Decode:", message)
	case wsjtx.LoggedAdifMessage:
		log.Println("Logged Adif:", message)
	default:
		log.Println("Other:", reflect.TypeOf(message), message)
	}
}
