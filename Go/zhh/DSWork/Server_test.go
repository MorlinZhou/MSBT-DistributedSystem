package main

import (
	"DSWork/Server/SysCall"
	"net"
	"testing"
	"time"
)

const LocalIp = "127.0.0.1"

func Test_server(t *testing.T) {
	localIp := net.ParseIP(LocalIp)
	var Port = []int{30000, 30001}

	s := SysCall.NewServer(localIp, Port)
	s.BootServer()

	time.Sleep(100 * time.Second)
	//Server()
}
