package main

import (
	"DSWork/Client/Arg"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	fmt.Println("client booting......")
	localIp := net.ParseIP("127.0.0.1")
	Port := 30000
	for i := 0; i < 2; i++ {
		fmt.Println(i)
		go func(Ip net.IP, P int) {
			go BootClient(Ip, P)
		}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名
		Port++
	}

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
	//Client()
}

func BootClient(IP net.IP, Port int) {

	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   IP,
		Port: Port,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}

	var Idata []byte
	Idata = append(Idata, []byte(fmt.Sprintf("%v", rand.Intn(10)))...)
	Idata = append(Idata, []byte(fmt.Sprintf("%v", rand.Intn(10)))...)
	Idata = append(Idata, []byte(fmt.Sprintf("%v", rand.Intn(10)))...)

	I1 := Arg.NewInsertArg("file/check.txt", 2, Idata)
	err = RpcCall("Insert", I1, socket)

	defer socket.Close()
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return
	}
	fmt.Printf("recv data: %v \naddr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}

func Client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP("127.0.0.1")
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   localIp,
		Port: 30000,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}

	// callfunction , arg
	//l := Arg.NewLookArg("file/check.txt", 0, 10)
	//err = RpcCall("LookUp", l, socket)

	Idata1 := []byte{'0', '0', '0'}
	I1 := Arg.NewInsertArg("file/check.txt", 2, Idata1)

	err = RpcCall("Insert", I1, socket)

	defer socket.Close()
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return
	}
	fmt.Printf("recv data: %v \naddr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}

type arg interface {
	StoB() (output []byte, err error)
}

func RpcCall(Func string, Arg arg, socket *net.UDPConn) error {
	var buf bytes.Buffer

	switch Func {
	case "LookUp":
		buf.Write([]byte(fmt.Sprintf("\\CallType LookUp ")))
	case "Insert":
		buf.Write([]byte(fmt.Sprintf("\\CallType Insert ")))
	}

	output, err := Arg.StoB()
	buf.Write(output)

	if err == nil {
		_, err = socket.Write(buf.Bytes()) // 发送数据

		if err != nil {
			fmt.Println("发送数据失败，err:", err)
		}
		return err
	}
	return fmt.Errorf("Argument Error")
}
