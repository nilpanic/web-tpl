package env

import (
	"net"
	"os"
)

func LocalIP() string {
	ipList := []string{"114.114.114.114:80", "8.8.8.8:80"}
	for _, ip := range ipList {
		conn, err := net.Dial("udp", ip)
		if err != nil {
			continue
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
		return localAddr.IP.String()
	}

	return ""
}

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
