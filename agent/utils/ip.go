package utils

import (
	"fmt"
	"net"
)

var (
	// 本地IP
	LocalIpArray []string
)

// 获取本地所有ip
func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("get local ip failed, err:", err))
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				LocalIpArray = append(LocalIpArray, ipnet.IP.String())
			}
		}
	}
}
