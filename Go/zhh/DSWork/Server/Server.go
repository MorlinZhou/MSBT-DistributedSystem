package Server

import (
	"DSWork/Server/SysCall"
	"net"
	"time"
)

func main() {
	//fmt.Println("Server booting......")
	localIp := net.ParseIP("127.0.0.1")
	var Port = []int{30000, 30001}

	s := SysCall.NewServer(localIp, Port)
	s.BootServer()

	time.Sleep(100 * time.Second)

}
