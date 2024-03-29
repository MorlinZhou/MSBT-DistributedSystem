package main

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
	//Server()
}

//func Server() {
//	var s SysCall.Server
//	fmt.Println("Server booting......")
//	localIp := net.ParseIP("127.0.0.1")
//	listen, err := net.ListenUDP("udp", &net.UDPAddr{
//		IP:   localIp,
//		Port: 30000,
//	})
//	if err != nil {
//		fmt.Println("listen failed, err:", err)
//		return
//	}
//	defer listen.Close()
//	for {
//		var data [1024]byte
//		n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
//		if err != nil {
//			fmt.Println("read udp failed, err:", err)
//			continue
//		}
//		fmt.Printf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)
//		SysCall.BtoS(data[:n], &s.Sfunc)
//		ReturnMsg, errMsg := s.ExecuteCall(addr)
//		if errMsg != nil {
//			listen.WriteToUDP([]byte(errMsg.Error()), addr)
//		}
//		_, err = listen.WriteToUDP(ReturnMsg, addr) // 发送数据
//		fmt.Printf("data:%v addr:%v count:%v\n", string(ReturnMsg), addr, n)
//		if err != nil {
//			fmt.Println("write to udp failed, err:", err)
//			continue
//		}
//	}
//}
