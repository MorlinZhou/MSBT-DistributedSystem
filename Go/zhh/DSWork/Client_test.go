package DSWork

import (
	"DSWork/Client"
	"DSWork/Client/Arg"
	"fmt"
	"net"
	"testing"
	"time"
)

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

func test1_client() {
	fmt.Println("client booting......")
	localIp := net.ParseIP("127.0.0.1")
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
	localIp := net.ParseIP("127.0.0.1")
	Port := 30000

	for i := 0; i < 2; i++ {
		index := i + 1
		if index > 2 {
			index = 1
		}

		filepath := fmt.Sprintf("file/check%v.txt", index)
		offset := 2
		bytes := []byte("123456")

		_Arg := Arg.NewInsertArg(filepath, offset, bytes)

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
	localIp := net.ParseIP("127.0.0.1")
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
