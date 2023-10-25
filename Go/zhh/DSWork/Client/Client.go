package Client

import (
	"DSWork/Client/Arg"
	"DSWork/Client/LocalFunc"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

func BootClientWithFresh(IP net.IP, Port int, _Arg arg, freshinter int) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   IP,
		Port: Port,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	LookArg := Arg.NewLookArg(_Arg.GetFilepath(), 0, 1)
	go func(socket *net.UDPConn) {
		for true {
			time.Sleep(time.Duration(freshinter) * time.Second)
			data := make([]byte, 1024)
			err = RpcCall(LookArg.Type(), LookArg, socket)
			if err != nil {
				fmt.Println("发送数据失败，err:", err)
				return
			}
			n, _, err := socket.ReadFromUDP(data)
			if err != nil {
				fmt.Println("接收数据失败，err:", err)
				return
			}
			i := 0
			for i < n {
				if data[i] == 127 { //判断结束符，设定结束符为delete字符 ascii码为127
					i++
					break
				}
				i++
			}
			log.Println("data is:", string(data[i:n]))
			LocalFunc.SaveToLocal(_Arg.GetFilepath(), data[i:n])
		}
	}(socket)

	if _Arg.Type() != "Monitor" && LocalFunc.CheckLocalFile(_Arg.GetFilepath()) { //本地有数据
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
				if data[i] == 127 { //判断结束符，设定结束符为delete字符 ascii码为127
					i++
					break
				}
				i++
			}
			fmt.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:i-1]), remoteAddr, Port, n-1)
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
	time.Sleep(50 * time.Second)
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

	if _Arg.Type() != "Monitor" && LocalFunc.CheckLocalFile(_Arg.GetFilepath()) { //本地有数据
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
				log.Printf("recv data: %v \naddr:%v port:%v \ncount:%v\n\n", string(data[:n]), remoteAddr, Port, n)
			}

			if string(data[:n]) == "monitor done!" {
				log.Println("exist monitor")
				break
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
	case "Append":
		buf.Write([]byte(fmt.Sprintf("\\CallType Append ")))
	case "Search":
		buf.Write([]byte(fmt.Sprintf("\\CallType Search ")))
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
