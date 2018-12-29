package ipp

import (
	"net"
)

func GetCurrentMachineIps() []string {
	var ips = make([]string, 0)

	address, err := net.InterfaceAddrs()
	if err != nil {
		Logger.Print(err)
		return ips
	}

	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Logger.Printf("found current machine ip %s", ipnet.IP.String())
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips
}