package SysCall

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type Server struct {
	IP   net.IP
	Port []int

	ch    chan *net.UDPAddr
	Sfunc ServerFunc
}

type ServerFunc struct {
	SystemCallFunc string
	FilePath       string
	Offset         int
	Number         int
	Bytes          []byte
}

func (s ServerFunc) print() {
	fmt.Println(s.SystemCallFunc)
	fmt.Println(s.FilePath)
	fmt.Println(s.Offset)
	fmt.Println(s.Number)
	fmt.Println(s.Bytes)
}
func (s ServerFunc) ExecuteCall() ([]byte, error) {
	switch s.SystemCallFunc {
	case "LookUp":
		return s.LookUp()
	case "Insert":
		return s.Insert()
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

	os.Rename("E:/DSWork/Server/file/temp", path)
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
