package Arg

//code bt zhh
import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
)

type LookArg struct {
	FilePath string
	Offset   int
	Number   int
}

func NewLookArg(filepath string, offset int, number int) (l *LookArg) {
	l = &LookArg{
		FilePath: filepath,
		Offset:   offset,
		Number:   number,
	}
	return
}

func (l *LookArg) StoB() (output []byte, err error) {
	var buf bytes.Buffer

	if len(l.FilePath) != 0 {
		buf.Write([]byte(fmt.Sprintf("\\FilePath %v ", l.FilePath)))
	} else {
		return []byte{}, fmt.Errorf("intput FilePath is null")
	}
	buf.Write([]byte(fmt.Sprintf("\\Offset %v ", l.Offset)))
	buf.Write([]byte(fmt.Sprintf("\\Number %v", l.Number)))
	return buf.Bytes(), nil
}

func (l *LookArg) Type() string {
	return "LookUp"
}

func (l *LookArg) GetFilepath() string {
	return l.FilePath
}

func (l *LookArg) LocalCall(socket *net.UDPConn) error {
	fmt.Printf("Open File Path is %v\n", l.FilePath)
	path := filepath.Join("E:/DSWork/Client", l.FilePath)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read fail")
		return err
	}
	defer f.Close()
	output := make([]byte, l.Number)
	_, err = f.ReadAt(output, int64(l.Offset))
	if err != nil {
		fmt.Println("read fail")
		return err
	}
	fmt.Printf("Read From Local\ndata is %v\n", string(output))
	return err
}

type InsertArg struct {
	FilePath string
	Offset   int
	Bytes    []byte
}

func NewInsertArg(filepath string, offset int, bytes []byte) (I *InsertArg) {
	I = &InsertArg{
		FilePath: filepath,
		Offset:   offset,
		Bytes:    bytes,
	}
	return
}

func (I *InsertArg) StoB() (output []byte, err error) {
	var buf bytes.Buffer

	if len(I.FilePath) != 0 {
		buf.Write([]byte(fmt.Sprintf("\\FilePath %v ", I.FilePath)))
	} else {
		return []byte{}, fmt.Errorf("intput FilePath is null")
	}
	buf.Write([]byte(fmt.Sprintf("\\Offset %v ", I.Offset)))
	buf.Write([]byte(fmt.Sprintf("\\Bytes %v", I.Bytes)))
	return buf.Bytes(), nil
}

func (I *InsertArg) Type() string {
	return "Insert"
}
func (I *InsertArg) GetFilepath() string {
	return I.FilePath
}
func (I *InsertArg) LocalCall(socket *net.UDPConn) error {

	path := filepath.Join("E:/DSWork/Client", I.FilePath)
	fmt.Printf("Open File Path is %v\n", path)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("read file fail")
		return err
	}
	//创建临时文件
	tempname := fmt.Sprintf("E:/DSWork/Client/file/temp%v", rand.Int())
	tempf, err := os.OpenFile(tempname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("create tempfile fail")
		return err
	}
	writebytes := make([]byte, 1024)

	n, err := f.Read(writebytes)

	if err != nil {
		fmt.Println("read file fail")
		return err
	}

	_, err = tempf.Write(writebytes[:I.Offset])
	if err != nil {
		fmt.Println("write file fail")
		return err
	}
	_, err = tempf.Write(I.Bytes)
	if err != nil {
		fmt.Println("write file fail")
		return err
	}
	_, err = tempf.Write(writebytes[I.Offset:n])
	if err != nil {
		fmt.Println("write file fail")
		return err
	}

	f.Close()
	tempf.Close()

	err = os.Rename(tempname, path)
	fmt.Println(path)

	if err != nil {
		fmt.Println(err)
		fmt.Println("rename error")
	}

	//fmt.Println("插入执行完毕")
	check, err := os.Open(path)
	defer check.Close()
	read := make([]byte, 1024)
	n, err = check.Read(read)
	data := read[:n]
	fmt.Printf("the content of local file %v has changed\ndata is %v\n", I.FilePath, string(data))

	//等待服务器响应
	data = make([]byte, 4096)
	var buf bytes.Buffer
	buf.Write([]byte(fmt.Sprintf("\\CallType Insert ")))
	output, err := I.StoB()
	buf.Write(output)
	if err == nil {
		_, err = socket.Write(buf.Bytes()) // 发送数据

		if err != nil {
			fmt.Println("发送数据失败，err:", err)
		}
		return err
	}
	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据

	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return err
	}
	fmt.Printf("require insert successfully, recv data: %v \naddr:%v count:%v\n\n", string(data[:n]), remoteAddr, n)

	return fmt.Errorf("Argument Error")
}

type MonitorArg struct {
	Filepath        string
	MonitorInterval int //单位秒
}

func NewMonitorArg(Filepath string, MonitorInterval int) (M *MonitorArg) {
	M = &MonitorArg{
		Filepath:        Filepath,
		MonitorInterval: MonitorInterval,
	}
	return
}

func (M *MonitorArg) StoB() (output []byte, err error) {
	var buf bytes.Buffer

	if len(M.Filepath) != 0 {
		buf.Write([]byte(fmt.Sprintf("\\FilePath %v ", M.Filepath)))
	} else {
		return []byte{}, fmt.Errorf("intput FilePath is null")
	}
	buf.Write([]byte(fmt.Sprintf("\\Monitor %v ", M.MonitorInterval)))
	return buf.Bytes(), nil
}

func (M *MonitorArg) Type() string {
	return "Monitor"
}
func (M *MonitorArg) GetFilepath() string {
	return M.Filepath
}
func (M *MonitorArg) LocalCall(socket *net.UDPConn) error {
	return nil
}
