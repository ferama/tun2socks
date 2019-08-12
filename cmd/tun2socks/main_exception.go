// +build exception

package main

import (
	"flag"
	"net"
	"strings"

	"github.com/xjasonlyu/tun2socks/common/log"
	"github.com/xjasonlyu/tun2socks/core"
	"github.com/xjasonlyu/tun2socks/proxy/exception"
	"github.com/xjasonlyu/tun2socks/proxy/socks"
)

func init() {
	args.ExceptionApps = flag.String("exceptionApps", "", "A list of exception apps separated by commas")
	args.ExceptionSendThrough = flag.String("exceptionSendThrough", "192.168.1.101:0", "Exception send through address")

	registerHandlerCreator("socks_exception", func() {
		// Verify proxy server address.
		proxyAddr, err := net.ResolveTCPAddr("tcp", *args.ProxyServer)
		if err != nil {
			log.Fatalf("invalid proxy server address: %v", err)
		}
		proxyHost := proxyAddr.IP.String()
		proxyPort := uint16(proxyAddr.Port)

		proxyTCPHandler := socks.NewTCPHandler(proxyHost, proxyPort, fakeDns, sessionStater)
		proxyUDPHandler := socks.NewUDPHandler(proxyHost, proxyPort, *args.UdpTimeout, fakeDns, sessionStater)

		sendThrough, err := net.ResolveTCPAddr("tcp", *args.ExceptionSendThrough)
		if err != nil {
			log.Fatalf("invalid exception send through address: %v", err)
		}
		apps := strings.Split(*args.ExceptionApps, ",")
		tcpHandler := exception.NewTCPHandler(proxyTCPHandler, apps, sendThrough)
		udpHandler := exception.NewUDPHandler(proxyUDPHandler, apps, sendThrough, *args.UdpTimeout)

		core.RegisterTCPConnHandler(tcpHandler)
		core.RegisterUDPConnHandler(udpHandler)
	})
}
