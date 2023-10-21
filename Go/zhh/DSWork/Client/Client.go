package Client

import (
	"DSWork/Client/LocalFunc"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {

	//fmt.Println("client booting......")
	//localIp := net.ParseIP("127.0.0.1")
	//Port := 30000
	//
	//for i := 0; i < 3; i++ {
	//	index := i + 1
	//	if index > 2 {
	//		index = 1
	//	}
	//	filepath := fmt.Sprintf("file/check%v.txt", index)
	//
	//	fmt.Println(i)
	//	go func(Ip net.IP, P int) {
	//		go BootClient(Ip, P, filepath)
	//	}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名
	//	Port++
	//	if Port == 30002 {
	//		Port = 30000
	//	}
	//}
	//
	//time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	//fmt.Println("client success......")
	////Client()
}

func BootClient(IP net.IP, Port int, _Arg arg) {

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

	//I1 := Arg.NewInsertArg("file/check1.txt", 2, Idata)
	//I1 := Arg.NewInsertArg(filepath, 2, Idata)
	//err = RpcCall("Insert", I1, socket)

	if LocalFunc.CheckLocalFile(_Arg.GetFilepath()) { //本地有数据
		_Arg.LocalCall(socket)
		return
	} else {
		err = RpcCall(_Arg.Type(), _Arg, socket)
	}

	defer socket.Close()
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	data := make([]byte, 4096)

	if _Arg.Type() != "Monitor" {
		n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据

		if err != nil {
			fmt.Println("接收数据失败，err:", err)
			return
		}
		fmt.Println("access file name is", _Arg.GetFilepath())
		//fmt.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", data[:n], remoteAddr, Port, n)

		if _Arg.Type() == "LookUp" {
			i := 0
			for i < n {
				if data[i] == 127 {
					i++
					break
				}
				i++
			}
			fmt.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:i-1]), remoteAddr, Port, n)
			LocalFunc.SaveToLocal(_Arg.GetFilepath(), data[i:n]) //full-chaching
		} else {
			fmt.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:n]), remoteAddr, Port, n)
		}

	} else {
		time.Sleep(500 * time.Millisecond)
		for true {

			n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据

			if err != nil {
				fmt.Println("接收数据失败，err:", err)
				return
			}
			fmt.Println("the file ", _Arg.GetFilepath())
			if n > 30 {
				fmt.Printf("%v\nfile data changed to: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:21]), string(data[21:n]), remoteAddr, Port, n)
			} else {
				fmt.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:n]), remoteAddr, Port, n)
			}

			data = make([]byte, 4096)
		}
	}

}

type arg interface {
	StoB() (output []byte, err error)
	Type() string
	GetFilepath() string
	LocalCall(socket *net.UDPConn) error
}

func RpcCall(Func string, Arg arg, socket *net.UDPConn) error {
	var buf bytes.Buffer

	switch Func {
	case "LookUp":
		buf.Write([]byte(fmt.Sprintf("\\CallType LookUp ")))
	case "Insert":
		buf.Write([]byte(fmt.Sprintf("\\CallType Insert ")))
	case "Monitor":
		buf.Write([]byte(fmt.Sprintf("\\CallType Monitor ")))
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
