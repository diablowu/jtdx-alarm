package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"jtdx-alarm/pkg/adif"
	"jtdx-alarm/pkg/city"
	"jtdx-alarm/pkg/monitor/decode"
	"jtdx-alarm/pkg/osx"
	"jtdx-alarm/pkg/qywx"
	"jtdx-alarm/pkg/wsjtx"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	bindAddr       = flag.String("bind-addr", "239.255.0.0", "接收JTDX UDP消息的地址，可以是组播地址或者本地地址")
	bindPort       = flag.Uint("bind-port", 2237, "接收JTDX UDP消息的端口")
	ctyPath        = flag.String("cty-path", "cty.dat", "Big CTY文件地址")
	targetCallSign = flag.String("target-call", "N0CALL", "接收消息的CALL")
	filteredDXCC   = flag.String("filtered-dxcc", "BY,JA,HL,BV,YB,UA,UA9", "被过滤掉的DXCC")
	notifiers      = flag.String("notifiers", "log,wx", "启用的通知器")
	jtdxLogDir     = flag.String("jtdx-log-dir", filepath.Join(osx.MustUserHomeDir(), "AppData", "Local", "JTDX"), "JTDX日志目录，一般不用修改。如果你有多个JTDX安装目录，需要配置此选项。")
	useJTDXLog     = flag.Bool("use-jtdx-log", true, "是否使用JTDX的日志过滤，如果启用，将会定时读取JTDX下的日志文件，仅仅针对没有记录的DXCC进行通知")
	logLevel       = flag.String("log-level", "info", "日志输出级别:panic,fatal,error,warn,info,debug,trace")

	qywxAgentID = flag.Int("qywx-agent-id", 1000002, "企业微信应用ID")
	qywxCorpId  = flag.String("qywx-corp-id", "wx861828161a3f015c", "企业微信企业ID")
	qywxSecret  = flag.String("qywx-secret", "q7EFHBUKk-S1pNBWD0pDXuYjDzahLZ2VaxQ7QfBrYeU", "企业微信应用密钥")

	verbose = flag.Bool("verbose", false, "是否输出详细日志")
)

var (
	DefaultADIFLogFileName     = "wsjtx_log.adi"
	DefaultADIFRefreshInterval = time.Minute * 5
)

var defaultDecodeMesageMonitor *decode.DecodeMessageMonitors

// Simple driver binary for wsjtx-go library.
func main() {

	flag.Parse()

	initLog()
	initQYWX()
	initBigCTY()

	log.Infof("Using JTDX log: %v", *useJTDXLog)
	if *useJTDXLog {
		adif.InitLoggerChecker(filepath.Join(*jtdxLogDir, DefaultADIFLogFileName), DefaultADIFRefreshInterval)
	} else {
		log.Infoln("JTDX logger checker was disabled.")
	}

	log.Infof("Use address %v:%d to receive wsjtx message", *bindAddr, *bindPort)

	defaultDecodeMesageMonitor = initDecodeMessageMonitors()
	wsjtxServer, err := wsjtx.MakeServerGiven(net.ParseIP(*bindAddr), *bindPort)
	if err != nil {
		log.Fatalf("Failed to create udp notify server, %v", err)
	}

	incomingMessageChannel := make(chan interface{}, 5)
	errChannel := make(chan error, 5)

	go wsjtxServer.Listen(incomingMessageChannel, errChannel)

	for {
		select {
		case err := <-errChannel:
			log.Printf("error: %v", err)
		case message := <-incomingMessageChannel:
			HandleServerMessage(message)
		}
	}

}

func initDecodeMessageMonitors() *decode.DecodeMessageMonitors {
	var finalFilter decode.CallSignFilter
	if *useJTDXLog {
		finalFilter = decode.NewCompositeFilter(decode.NewBlacklistDXCCCallSignFilter(strings.Split(*filteredDXCC, ",")), decode.NewADIFFilter())
	} else {
		finalFilter = decode.NewBlacklistDXCCCallSignFilter(strings.Split(*filteredDXCC, ","))
	}

	return decode.CreateDecodeMessageMonitors(
		decode.NewDefaultMonitor(finalFilter, strings.Split(*notifiers, ",")))
}

func initQYWX() {
	log.Infoln("Init qywx...")
	qywx.Setup(*qywxAgentID, *qywxCorpId, *qywxSecret, strings.ToLower(*targetCallSign))
}

func initBigCTY() {
	log.Infoln("Init bigcty...")
	if err := city.LoadFromCTYData(*ctyPath); err != nil {
		log.Fatalf("Failed to load cty file %v", err)
	} else {
		log.Infoln("Success to load cty data")
	}
}

func initLog() {

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(*verbose)
	if lvl, err := log.ParseLevel(*logLevel); err == nil {
		log.SetLevel(lvl)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func HandleServerMessage(message interface{}) {
	switch message.(type) {
	case wsjtx.HeartbeatMessage:
		id := message.(wsjtx.HeartbeatMessage).Id
		log.Debugf("Heartbeat(%s): %v", id, message)
		//s.targetName = &id
	case wsjtx.StatusMessage:
		log.Debugf("Status: %s %v", reflect.TypeOf(message), message)
	case wsjtx.DecodeMessage:
		defaultDecodeMesageMonitor.Do(message.(wsjtx.DecodeMessage))
	case wsjtx.ClearMessage:
		log.Debugf("Clear: %v", message)
	case wsjtx.QsoLoggedMessage:
		log.Debugf("QSO Logged: %v", message)
	case wsjtx.CloseMessage:
		log.Debugf("Close: %v", message)
	case wsjtx.WSPRDecodeMessage:
		log.Debugf("WSPR Decode: %v", message)
	case wsjtx.LoggedAdifMessage:
		log.Debugf("Logged Adif: %v", message)
	default:
		log.Debugf("Other: %s %v", reflect.TypeOf(message), message)
	}
}
