package main

import (
	"DSWork/Client"
	"DSWork/Client/Arg"
	"fmt"
	"net"
	"testing"
	"time"
)

const ServerIp = "172.20.10.7" //"127.0.0.1"

// 测试功能一，多端口读数据
func TestRead(t *testing.T) {
	fmt.Println("execute remote read function")
	go test1_client()
	time.Sleep(100 * time.Second)

}

func TestInsert(t *testing.T) {
	fmt.Println("execute remote Insert function")
	go test2_client()
	time.Sleep(100 * time.Second)
}

func TestMonitor(t *testing.T) {
	fmt.Println("execute Monitor function")
	go test3_client()
	time.Sleep(100 * time.Second)
}

func TestInsert2(t *testing.T) {
	fmt.Println("execute local Insert function")
	go test2_client()
	time.Sleep(100 * time.Second)
}

func TestAppend(t *testing.T) {
	fmt.Println("execute local Append function")
	go test4_client()
	time.Sleep(100 * time.Second)
}

func TestSearch(t *testing.T) {
	fmt.Println("execute local Search function")
	go test5_client()
	time.Sleep(100 * time.Second)
}

func TestFresh(t *testing.T) {
	fmt.Println("execute Fresh function")
	go test_fresh()
	time.Sleep(100 * time.Second)
}

func test1_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	for i := 0; i < 2; i++ {
		index := i + 1
		if index > 2 {
			index = 1
		}

		filepath := fmt.Sprintf("file/check%v.txt", index)
		offset := 2
		number := 2

		_Arg := Arg.NewLookArg(filepath, offset, number)

		fmt.Println(i)
		go func(Ip net.IP, P int) {
			go Client.BootClient(Ip, P, _Arg)
		}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名
		Port++
		if Port == 30002 {
			Port = 30000
		}
	}

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}

func test2_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	for i := 0; i < 3; i++ {
		index := i + 1
		if index > 2 {
			index = 1
		}

		filepath := fmt.Sprintf("file/check%v.txt", index)
		offset := 2
		bytes := []byte("123456")

		_Arg := Arg.NewInsertArg(filepath, offset, bytes[i:])

		fmt.Println(i)
		go func(Ip net.IP, P int) {
			go Client.BootClient(Ip, P, _Arg)
		}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名
		Port++
		if Port == 30002 {
			Port = 30000
		}
	}

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}

func test3_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	filepath := fmt.Sprintf("file/check1.txt")
	monitor := 3

	_ArgM := Arg.NewMonitorArg(filepath, monitor)

	fmt.Println("require monitor")
	go func(Ip net.IP, P int) {
		go Client.BootClient(Ip, P, _ArgM)
	}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名

	Port++
	time.Sleep(time.Second)
	offset := 2
	bytes := []byte("123")
	_ArgI := Arg.NewInsertArg(filepath, offset, bytes)

	fmt.Println("require insert")
	go func(Ip net.IP, P int) {
		go Client.BootClient(Ip, P, _ArgI)
	}(localIp, Port) //func后面是自定义变量，可以改变量名字，后面括号里面是传入的对应变量名

	time.Sleep(10 * time.Second) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}

func test4_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	filepath := fmt.Sprintf("file/check1.txt")
	bytes := []byte("1212a")

	_Arg := Arg.NewAppendArg(filepath, bytes)

	Client.BootClient(localIp, Port, _Arg)

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}

func test5_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	filepath := fmt.Sprintf("file/check1.txt")
	bytes := []byte("HF")

	_Arg := Arg.NewSearchArg(filepath, bytes)

	Client.BootClient(localIp, Port, _Arg)

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}

func test_fresh() {
	fmt.Println("client booting......")
	localIp := net.ParseIP(ServerIp)
	Port := 30000

	filepath := fmt.Sprintf("file/check1.txt")
	bytes := []byte("HF")

	_Arg := Arg.NewSearchArg(filepath, bytes)

	Client.BootClientWithFresh(localIp, Port, _Arg, 5)

	time.Sleep(500 * time.Millisecond) //需要等待，要不然变量没传进去就释放了
	fmt.Println("client success......")
}
