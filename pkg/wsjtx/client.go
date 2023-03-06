package wsjtx

import (
	"fmt"
	"net"
)

type Client struct {
	conn *net.UDPConn
}

func MakeClient(port uint) (*Client, error) {
	if remoteUDPAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
		if conn, err := net.DialUDP("udp4", nil, remoteUDPAddr); err == nil {
			return &Client{conn: conn}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

}

func (c *Client) TriggerCQ(msg TriggerCQMessage) error {
	msgBytes, _ := encodeTriggerCQ(msg)
	return c.tryWrite(msgBytes)
}

func (c *Client) tryWrite(msgBytes []byte) error {
	_, err := c.conn.Write(msgBytes)
	return err
}
