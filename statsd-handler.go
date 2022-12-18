package main

import (
	"fmt"
	"net"

	"sideproject/dogdata/common"
)

func serveStatsd() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 7125})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		common.CheckErr(err, "listener.ReadFromUDP()")
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}
