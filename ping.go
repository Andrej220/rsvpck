package main

import (
	"net"
	"os"
	"runtime"
	"time"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Ping(host string, attempts int, timeout time.Duration) bool {
	if attempts <= 0 {
		attempts = 1
	}
	
	ip, err := resolveIPv4(host)
	if err != nil {
		return false
	}

	network := chooseNetwork() // "ip4:icmp" on Unix, "udp4" on Windows
	conn, err := openICMP(network)
	if err != nil {
		return false
	}
	defer conn.Close()

	dest := makeDestAddr(network, ip)
	id := os.Getpid() & 0xffff

	for seq := 1; seq <= attempts; seq++ {
		packet, err := buildEcho(id, seq, []byte("PING"))
		
		
		if err != nil {
			continue
		}
		if err := sendEcho(conn, dest, packet); err != nil {
			continue
		}
		msg, err := readICMPReply(conn, timeout)
		if err != nil {
			continue
		}
		if isMatchingEchoReply(msg, id, seq) {
			return true
		}
	}
	return false
}

func chooseNetwork() string {
	if runtime.GOOS == "windows" {
		return "udp4"
	}
	return "ip4:icmp"
}

func openICMP(network string) (net.PacketConn, error) {
	return icmp.ListenPacket(network, "0.0.0.0") // Bind to all interfaces
}

func makeDestAddr(network string, ip *net.IPAddr) net.Addr {
	if network == "udp4" {
		return &net.UDPAddr{IP: ip.IP}
	}
	return ip // "ip4:icmp"
}

func buildEcho(id, seq int, payload []byte) ([]byte, error) {
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: payload,
		},
	}
	return m.Marshal(nil)
}

func sendEcho(c net.PacketConn, dst net.Addr, packet []byte) error {
	_, err := c.WriteTo(packet, dst)
	return err
}

func readICMPReply(c net.PacketConn, timeout time.Duration) (*icmp.Message, error) {
	_ = c.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1500)
	n, _, err := c.ReadFrom(buf)
	if err != nil {
		return nil, err
	}
	return icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buf[:n])
}

func isMatchingEchoReply(r *icmp.Message, id, seq int) bool {
	if r.Type != ipv4.ICMPTypeEchoReply {
		return false
	}
	echo, ok := r.Body.(*icmp.Echo)
	if !ok {
		return false
	}

	return echo.ID == id // && echo.Seq == seq
}

func resolveIPv4(host string) (*net.IPAddr, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, &net.DNSError{Err: "not IPv4", Name: host}
		}
		return &net.IPAddr{IP: ip4}, nil
	}
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, err
	}
	if addr.IP.To4() == nil {
		return nil, &net.DNSError{Err: "not IPv4", Name: host}
	}
	return addr, nil
}