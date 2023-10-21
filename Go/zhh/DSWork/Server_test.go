package DSWork

import (
	"DSWork/Server/SysCall"
	"net"
	"testing"
	"time"
)

func Test_server(t *testing.T) {
	localIp := net.ParseIP("127.0.0.1")
	var Port = []int{30000, 30001}

	s := SysCall.NewServer(localIp, Port)
	s.BootServer()

	time.Sleep(100 * time.Second)
	//Server()
}
