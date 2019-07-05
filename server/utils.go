package main

import (
	"log"
	"net"
	"os"
)

func getInternalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Oops:" + err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	os.Exit(0)
	return ""
}
