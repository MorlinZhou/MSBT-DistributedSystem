package SysCall

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

//var prepath = "E:/DSWork/Server"
//var pretmppath = "E:/DSWork/Server/file/temp"

type Server struct {
	IP   net.IP
	Port []int

	mu sync.Mutex

	Waitlist  []WaitList
	IsUpdate  map[string]bool
	InsertNum map[string]int
	Sfunc     []ServerFunc //每个port接受并执行一个func
}

type WaitList struct { //用于监听操作 monitor
	FilePath string
	Raddress *net.UDPAddr
	Listen   *net.UDPConn
	Interval int
	Port     int
	Nowtime  time.Time
	used     bool
}

type ServerFunc struct {
	SystemCallFunc string
	FilePath       string
	Offset         int
	Number         int
	MonitorInter   int
	Bytes          []byte
	mu             sync.Mutex
}

func (s *ServerFunc) Clear() {
	s.FilePath = ""
	s.SystemCallFunc = ""
	s.Offset = 0
	s.Number = 0
	s.Bytes = nil
}
func (w *WaitList) Clear() {
	w.FilePath = ""
	w.Raddress = nil
	w.Listen = nil
	w.Interval = 0
	w.Port = 0
	w.Nowtime = time.Now()
	w.used = false
}

func NewServer(IP net.IP, Port []int) *Server {
	s := &Server{
		IP:       IP,
		Port:     Port,
		IsUpdate: make(map[string]bool),

		InsertNum: make(map[string]int),
		Waitlist:  make([]WaitList, 20), //monitor队列长度
		Sfunc:     make([]ServerFunc, len(Port)),
	}

	return s
}

func (s *Server) print() {
	fmt.Println(s.Sfunc[0].SystemCallFunc)
	fmt.Println(s.Sfunc[0].FilePath)
	fmt.Println(s.Sfunc[0].Offset)
	fmt.Println(s.Sfunc[0].Number)
	fmt.Println(s.Sfunc[0].Bytes)
}

func (s *Server) BootServer() {

	//fmt.Println("Server booting......")
	go s.StartMonitor()

	for i, P := range s.Port {
		go func(i int, P int) {
			fmt.Printf("Server sets up successfully,Port is %v\n", P)
			listen, err := net.ListenUDP("udp", &net.UDPAddr{
				IP:   s.IP,
				Port: P,
			})
			if err != nil {
				fmt.Println("listen failed, err:", err)
				return
			}
			defer listen.Close()
			for {
				var data [1024]byte
				n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
				if err != nil {
					fmt.Println("read udp failed, err:", err)
					continue
				}
				fmt.Printf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)
				BtoS(data[:n], &s.Sfunc[i]) //赋值给S

				ReturnMsg, errMsg := s.ExecuteCall(i, addr, listen) //第i个Port执行操作

				if errMsg != nil {
					listen.WriteToUDP([]byte(errMsg.Error()), addr)
				}
				_, err = listen.WriteToUDP(ReturnMsg, addr) // 发送数据
				fmt.Printf("access filename is:%v \ndata:%v addr:%v count:%v\n", filepath.Base(s.Sfunc[i].FilePath), string(ReturnMsg), addr, n)
				if err != nil {
					fmt.Println("write to udp failed, err:", err)
					continue
				}

				if s.Sfunc[i].SystemCallFunc == "Monitor" {
					for j := 0; j < len(s.Waitlist); j++ {
						for s.Waitlist[j].Raddress == addr && s.Waitlist[j].used == true { //判断是否有monitor
							time.Sleep(time.Second)
							if s.Waitlist[j].used == false {
								listen.WriteToUDP([]byte("monitor done!"), addr)
							}
						}
					}
				}
				s.Sfunc[i].Clear() //清空结构体，防止filepath参数对后续产生影响

			}
		}(i, P)

	}

	//主程序提前结束时，可能会导致已经在运行的 goroutine 不再具有上下文，并且可能无法正常完成,但是goroutine是独立于主程序的
	time.Sleep(100 * time.Second)
}

func (s *Server) ExecuteCall(index int, addr *net.UDPAddr, listen *net.UDPConn) (B []byte, e error) {
	switch s.Sfunc[index].SystemCallFunc {
	case "LookUp":
		return s.Sfunc[index].LookUp()
	case "Insert":
		//实现顺序互斥执行的重点
		//s.mu.Lock()
		//B, e = s.Sfunc[index].Insert()
		//s.mu.Unlock()
		//return B, e
		for i := 0; i < len(s.Port); i++ {
			if s.Sfunc[i].FilePath == s.Sfunc[index].FilePath {
				s.Sfunc[i].mu.Lock() //如果有新的请求，访问同样的filepath，会发现锁被占用了，并wait
				defer s.Sfunc[i].mu.Unlock()
			}
		}
		B, e = s.Sfunc[index].Insert()
		s.IsUpdate[s.Sfunc[index].FilePath] = true
		go s.callback(s.Sfunc[index].FilePath, B)
		return B, e
	case "Monitor":
		B, e = s.Monitor(index, addr, listen)
		return
	}
	return []byte{}, nil
}

func (s *ServerFunc) LookUp() ([]byte, error) {
	fmt.Printf("Open File Path is %v\n", s.FilePath)
	path := filepath.Join("E:/DSWork/Server", s.FilePath)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read fail")
	}
	defer f.Close()
	output := make([]byte, s.Number)
	_, err = f.ReadAt(output, int64(s.Offset))
	if err != nil {
		return nil, err
	}

	//return output,err // no full-caching

	fulldata := make([]byte, 1024) //全部数据
	n, err := f.Read(fulldata)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.Write(output)
	deleteChar := rune(0x007F)
	buf.WriteRune(deleteChar)
	//buf.WriteString("\n")
	buf.Write(fulldata[:n])
	return buf.Bytes(), err
}

func (s *ServerFunc) Insert() ([]byte, error) {
	//_, filename := filepath.Split(s.FilePath)

	fmt.Printf("Open File Path is %v\n", s.FilePath)
	path := filepath.Join("E:/DSWork/Server", s.FilePath)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("read file fail")
		return nil, err
	}
	//创建临时文件
	tempname := fmt.Sprintf("E:/DSWork/Server/file/temp%v", rand.Int())
	tempf, err := os.OpenFile(tempname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("create tempfile fail")
		return nil, err
	}
	writebytes := make([]byte, 1024)

	n, err := f.Read(writebytes)

	if err != nil {
		fmt.Println("read file fail")
		return nil, err
	}

	_, err = tempf.Write(writebytes[:s.Offset])
	if err != nil {
		fmt.Println("write file fail")
		return nil, err
	}
	_, err = tempf.Write(s.Bytes)
	if err != nil {
		fmt.Println("write file fail")
		return nil, err
	}
	_, err = tempf.Write(writebytes[s.Offset:n])
	if err != nil {
		fmt.Println("write file fail")
		return nil, err
	}

	f.Close()
	tempf.Close()

	err = os.Rename(tempname, path)
	if err != nil {
		fmt.Println("rename error")
	}

	//fmt.Println("插入执行完毕")
	check, err := os.Open(path)
	defer check.Close()
	read := make([]byte, 1024)
	n, err = check.Read(read)
	return read[:n], nil
}

func (s *Server) Monitor(index int, addr *net.UDPAddr, listen *net.UDPConn) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := "succeed,interval is "

	//使用range迭代时得到的是原始数据的副本，而不是原始数据本身,对副本的修改不会影响原始数据本身
	for i := 0; i < len(s.Waitlist); i++ { //不能用range s.Waitlist ,range是值传递
		if s.Waitlist[i].used == true {
			continue
		} else {
			s.Waitlist[i].Port = s.Port[index]
			s.Waitlist[i].FilePath = s.Sfunc[index].FilePath
			s.Waitlist[i].Raddress = addr
			s.Waitlist[i].Listen = listen
			s.Waitlist[i].Interval = s.Sfunc[index].MonitorInter
			s.Waitlist[i].used = true
			s.Waitlist[i].Nowtime = time.Now() //记录当前时间
			out = out + fmt.Sprintf("%v", s.Waitlist[i].Interval)
			break
		}
	}
	if out == "monitor set successfully,interval is " {
		out = "the waiting list is full!"
		return []byte(out), nil
	}

	out = out + " second!"
	return []byte(out), nil
}

func (s *Server) StartMonitor() {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	for true {
		for i := 0; i < len(s.Waitlist); i++ {
			if s.Waitlist[i].used == false {
				continue
			} else {
				NowTime := time.Now()
				spendtime := NowTime.Sub(s.Waitlist[i].Nowtime) //计算过去的时间
				if spendtime >= time.Duration(s.Waitlist[i].Interval)*time.Second {
					//fmt.Println("clear")
					s.Waitlist[i].Clear() //大于监听时间则清空
				}
			}
		}

		time.Sleep(500 * time.Millisecond) //每0.5秒验证一次
	}
}

// 如果文件被修改调用callback函数给waitlist发消息
func (s *Server) callback(Filepath string, data []byte) {
	for _, m := range s.Waitlist {
		if m.used == true && m.FilePath == Filepath {

			Listen := m.Listen
			ReturnMsg := []byte("the file has changed!")
			ReturnMsg = append(ReturnMsg, data...)

			time.Sleep(10 * time.Millisecond)
			_, err := Listen.WriteToUDP(ReturnMsg, m.Raddress) // 发送数据

			if err != nil {
				fmt.Println("write to udp failed, err:", err)
				continue
			}

			fmt.Printf("send message successfully, the address is %v\n", m.Raddress)
		}
	}
	s.IsUpdate[Filepath] = false
}

func BtoS(data []byte, S *ServerFunc) (err error) {
	var s string
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			switch s {
			case "\\CallType":
				i++
				out, index := Read(data[i:])
				S.SystemCallFunc = string(out)
				fmt.Println(S.SystemCallFunc)
				i += index
			case "\\FilePath":
				i++
				out, index := Read(data[i:])
				S.FilePath = string(out)
				i += index
			case "\\Offset":
				i++
				out, index := Read(data[i:])
				S.Offset, err = strconv.Atoi(string(out))
				i += index
			case "\\Number":
				i++
				out, index := Read(data[i:])
				S.Number, err = strconv.Atoi(string(out))
				i += index
			case "\\Bytes":
				i++
				out, index := ReadBytes(data[i:])
				S.Bytes = out
				i += index
				//fmt.Println(out)
			case "\\Monitor":
				i++
				out, index := Read(data[i:])
				S.MonitorInter, err = strconv.Atoi(string(out))
				i += index
			}
			s = ""
		} else {
			s += string(data[i])
		}
	}
	return nil
}
func ReadBytes(data []byte) ([]byte, int) {
	var bytes []byte
	var s string
	i := 1
	for i < len(data) {
		if data[i] == ' ' {
			num, _ := strconv.Atoi(s)
			Byte := byte(num)
			bytes = append(bytes, Byte)
			s = ""
		} else if data[i] == ']' {
			num, _ := strconv.Atoi(s)
			Byte := byte(num)
			bytes = append(bytes, Byte)
			return bytes, i
		} else {
			s += string(data[i])
		}
		i++
	}
	fmt.Println(bytes)
	return bytes, i
}
func Read(data []byte) ([]byte, int) {
	var s string
	i := 0
	for i < len(data) {
		if data[i] == ' ' {
			return []byte(s), i
		} else {
			s += string(data[i])
		}
		i++
	}
	return []byte(s), i
}
