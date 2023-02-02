package wsjtx

import (
	"errors"
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

const magic = 0xadbccbda
const schema = 2
const qDataStreamNull = 0xffffffff
const bufLen = 1024

type Server struct {
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
	listening  bool
}

var NotConnectedError = fmt.Errorf("haven't heard from wsjtx yet, don't know where to send commands")

// MakeServer creates a multicast UDP connection to communicate with WSJT-X on the default address
// and port.
func MakeServer(addr string, port uint) (Server, error) {
	//var defaultWsjtxAddr net.IP
	//switch runtime.GOOS {
	//case "windows":
	//	defaultWsjtxAddr = net.ParseIP(localhostAddr)
	//default:
	//	defaultWsjtxAddr = net.ParseIP(multicastAddr)
	//}
	return MakeServerGiven(net.ParseIP(addr), port)
}

func makeMulticastServer(addr *net.UDPAddr) (*net.UDPConn, error) {

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}

	packetConn := ipv4.NewPacketConn(conn)

	if err := packetConn.JoinGroup(nil, addr); err != nil {
		return nil, err
	}

	// test
	if loop, err := packetConn.MulticastLoopback(); err == nil {
		log.Printf("MulticastLoopback status:%v\n", loop)
		if !loop {
			if err := packetConn.SetMulticastLoopback(true); err != nil {
				return nil, err
			}
		}
	}
	return conn, nil
}

// MakeServerGiven creates a UDP connection to communicate with WSJT-X on the given address and
// port.
func MakeServerGiven(ipAddr net.IP, port uint) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ipAddr, port))
	if err != nil {
		return Server{}, err
	}
	var conn *net.UDPConn
	if ipAddr.IsMulticast() {
		conn, err = makeMulticastServer(addr)
	} else {
		log.Printf("%v is normal addrerss", ipAddr)
		conn, err = net.ListenUDP(addr.Network(), addr)
	}
	if err != nil {
		return Server{}, err
	}
	if conn == nil {
		return Server{}, errors.New("wsjtx udp connection not opened")
	}
	return Server{conn, nil, false}, nil
}

//func MakeServerGiven0(ipAddr net.IP, port uint) (Server, error) {
//	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ipAddr, port))
//	if err != nil {
//		return Server{}, err
//	}
//	var conn *net.UDPConn
//	if ipAddr.IsMulticast() {
//		log.Printf("%v is multicast addr", ipAddr)
//		var iface *net.Interface
//		iface, err = net.InterfaceByName("WLAN")
//		if iface != nil {
//			log.Printf("Broadcast iface :%v", iface)
//		}
//		conn, err = net.ListenMulticastUDP(addr.Network(), iface, addr)
//	} else {
//		log.Printf("%v is normal addrerss", ipAddr)
//		conn, err = net.ListenUDP(addr.Network(), addr)
//	}
//	if err != nil {
//		return Server{}, err
//	}
//	if conn == nil {
//		return Server{}, errors.New("wsjtx udp connection not opened")
//	}
//	return Server{conn, nil, false}, nil
//}

func (s *Server) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// ListenToWsjtx listens for messages from WSJT-X. When heard, the messages are parsed and then
// placed in the given message channel. If parsing errors occur, those are reported on the errors
// channel. If a fatal error happens, e.g. the network connection gets closed, the channels are
// closed and the goroutine ends.
func (s *Server) ListenToWsjtx(incomingMessageChannel chan interface{}, errChannel chan error) {
	s.listening = true
	defer close(incomingMessageChannel)
	defer close(errChannel)

	for {

		b := make([]byte, bufLen)
		if s.conn == nil {
			errChannel <- errors.New("wsjtx connection is nil")
			s.listening = false
			return
		}
		length, rAddr, err := s.conn.ReadFromUDP(b)
		if err != nil {
			errChannel <- fmt.Errorf("problem reading from wsjtx: %w", err)
			s.listening = false
			return
		}
		s.remoteAddr = rAddr
		message, err := parseMessage(b, length)
		if err != nil {
			errChannel <- err
		}
		if message != nil {
			incomingMessageChannel <- message
		}

	}
}

// Listening returns whether the ListenToWsjtx goroutine is currently running.
func (s *Server) Listening() bool {
	return s.listening
}

// Heartbeat sends a heartbeat message to WSJT-X.
func (s *Server) Heartbeat(msg HeartbeatMessage) error {
	msgBytes, _ := encodeHeartbeat(msg)
	return s.tryWrite(msgBytes)
}

// Clear sends a message to WSJT-X to clear the band activity window, the RX frequency window, or
// both.
func (s *Server) Clear(msg ClearMessage) error {
	msgBytes, _ := encodeClear(msg)
	return s.tryWrite(msgBytes)
}

// Reply initiates a reply to an earlier decode. The decode message must have started with CQ or
// QRZ.
func (s *Server) Reply(msg ReplyMessage) error {
	msgBytes, _ := encodeReply(msg)
	return s.tryWrite(msgBytes)
}

// Close sends a message to WSJT-X to close the program.
func (s *Server) Close(msg CloseMessage) error {
	msgBytes, _ := encodeClose(msg)
	return s.tryWrite(msgBytes)
}

// Replay sends a message to WSJT-X to replay QSOs in the Band Activity window.
func (s *Server) Replay(msg ReplayMessage) error {
	msgBytes, _ := encodeReplay(msg)
	return s.tryWrite(msgBytes)
}

// HaltTx sends a message to WSJT-X to halt transmission.
func (s *Server) HaltTx(msg HaltTxMessage) error {
	msgBytes, _ := encodeHaltTx(msg)
	return s.tryWrite(msgBytes)
}

// FreeText sends a message to WSJT-X to set the free text of the TX message.
func (s *Server) FreeText(msg FreeTextMessage) error {
	msgBytes, _ := encodeFreeText(msg)
	return s.tryWrite(msgBytes)
}

// Location sends a message to WSJT-X to set this station's Maidenhead grid.
func (s *Server) Location(msg LocationMessage) error {
	msgBytes, _ := encodeLocation(msg)
	return s.tryWrite(msgBytes)
}

// HighlightCallsign sends a message to WSJT-X to set callsign highlighting.
func (s *Server) HighlightCallsign(msg HighlightCallsignMessage) error {
	msgBytes, _ := encodeHighlightCallsign(msg)
	return s.tryWrite(msgBytes)
}

// SwitchConfiguration sends a message to WSJT-X to switch to a different pre-defined configuration.
func (s *Server) SwitchConfiguration(msg SwitchConfigurationMessage) error {
	msgBytes, _ := encodeSwitchConfiguration(msg)
	return s.tryWrite(msgBytes)
}

// Configure sends a message to WSJT-X to change various configuration options.
func (s *Server) Configure(msg ConfigureMessage) error {
	msgBytes, _ := encodeConfigure(msg)
	return s.tryWrite(msgBytes)
}

func (s *Server) TriggerCQ(msg TriggerCQMessage) error {
	msgBytes, _ := encodeTriggerCQ(msg)
	return s.tryWrite(msgBytes)
}

func (s *Server) tryWrite(msgBytes []byte) error {
	if s.remoteAddr == nil {
		return NotConnectedError
	}
	log.Printf("Try to write msg to : %v", s.remoteAddr)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}
