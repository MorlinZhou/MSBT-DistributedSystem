package SysCall

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	IP   net.IP
	Port []int

	mu        sync.Mutex
	InsertNum map[string]int
	Sfunc     []ServerFunc //每个port接受并执行一个func
}

type ServerFunc struct {
	SystemCallFunc string
	FilePath       string
	Offset         int
	Number         int
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

func NewServer(IP net.IP, Port []int) *Server {
	s := &Server{
		IP:   IP,
		Port: Port,

		InsertNum: make(map[string]int),
		Sfunc:     make([]ServerFunc, len(Port)),
	}

	return s
}

func (s Server) print() {
	fmt.Println(s.Sfunc[0].SystemCallFunc)
	fmt.Println(s.Sfunc[0].FilePath)
	fmt.Println(s.Sfunc[0].Offset)
	fmt.Println(s.Sfunc[0].Number)
	fmt.Println(s.Sfunc[0].Bytes)
}

func (s Server) BootServer() {

	//fmt.Println("Server booting......")

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

				ReturnMsg, errMsg := s.ExecuteCall(i) //第i个Port执行操作

				if errMsg != nil {
					listen.WriteToUDP([]byte(errMsg.Error()), addr)
				}
				_, err = listen.WriteToUDP(ReturnMsg, addr) // 发送数据
				fmt.Printf("data:%v addr:%v count:%v\n", string(ReturnMsg), addr, n)
				if err != nil {
					fmt.Println("write to udp failed, err:", err)
					continue
				}
				s.Sfunc[i].Clear()//清空结构体，防止filepath参数对后续产生影响

			}
		}(i, P)

	}

	time.Sleep(100 * time.Second)
}

func (s *Server) ExecuteCall(index int) (B []byte, e error) {
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
		return B, e
	}
	return []byte{}, nil
}
func (s ServerFunc) LookUp() ([]byte, error) {
	fmt.Printf("Open File Path is %v\n", s.FilePath)
	path := filepath.Join("E:/DSWork/Server", s.FilePath)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read fail")
	}
	defer f.Close()
	output := make([]byte, s.Number)
	f.ReadAt(output, int64(s.Offset))
	return output, err
}

func (s ServerFunc) Insert() ([]byte, error) {
	//_, filename := filepath.Split(s.FilePath)

	fmt.Printf("Open File Path is %v\n", s.FilePath)
	path := filepath.Join("E:/DSWork/Server", s.FilePath)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("read file fail")
		return nil, err
	}
	tempf, err := os.OpenFile("E:/DSWork/Server/file/temp", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

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

	err = os.Rename("E:/DSWork/Server/file/temp", path)
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
