package ssdp

import (
	"net"
	"time"
)

var (
	sendAddrIPv4 = "239.255.255.250:1900"
	recvAddrIPv4 = "224.0.0.0:1900"
	ssdpAddrIPv4 *net.UDPAddr
)

func init() {
	var err error
	ssdpAddrIPv4, err = net.ResolveUDPAddr("udp4", sendAddrIPv4)
	if err != nil {
		panic(err)
	}
}

type packetHandler func(net.Addr, []byte) (bool, error)

func readPackets(conn *net.UDPConn, timeout time.Duration, h packetHandler) error {
	buf := make([]byte, 65535)
	conn.SetReadBuffer(len(buf))
	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				return nil
			}
			return err
		}

		ret, err := h(addr, buf[:n])
		if ret || err != nil {
			return err
		}
	}
}

func sendTo(to *net.UDPAddr, data []byte) (int, error) {
	conn, err := net.DialUDP("udp4", nil, to)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	n, err := conn.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
