package LocalFunc

import (
	"DSWork/Client/Arg"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const LocalFilePath = "E:/DSWork/Client"

func SaveToLocal(Filepath string, data []byte) {
	path := filepath.Join(LocalFilePath, Filepath)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) //分别对应创建，只写，清空
	defer f.Close()
	if err != nil {
		fmt.Println("create file fail")
		return
	}
	_, err = f.Write(data)
	if err != nil {
		fmt.Println("write file fail")
		return
	}
	fmt.Println("create file successfully")

}

func CheckLocalFile(Filepath string) bool {
	path := filepath.Join(LocalFilePath, Filepath)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}

func ReadFromFile(_Arg *Arg.LookArg) {
	fmt.Printf("Open File Path is %v\n", _Arg.FilePath)
	path := filepath.Join(LocalFilePath, _Arg.FilePath)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read fail")
	}
	defer f.Close()
	output := make([]byte, _Arg.Number)
	_, err = f.ReadAt(output, int64(_Arg.Offset))
	fmt.Printf("Read From Local\ndata is %v\n", string(output))
}

func InsertFile(_Arg *Arg.InsertArg) error {
	//_, filename := filepath.Split(s.FilePath)

	fmt.Printf("Open File Path is %v\n", _Arg.FilePath)
	path := filepath.Join(LocalFilePath, _Arg.FilePath)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("read file fail")
		return err
	}
	//创建临时文件
	tempname := fmt.Sprintf("%v/file/temp%v", LocalFilePath, rand.Int())
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

	_, err = tempf.Write(writebytes[:_Arg.Offset])
	if err != nil {
		fmt.Println("write file fail")
		return err
	}
	_, err = tempf.Write(_Arg.Bytes)
	if err != nil {
		fmt.Println("write file fail")
		return err
	}
	_, err = tempf.Write(writebytes[_Arg.Offset:n])
	if err != nil {
		fmt.Println("write file fail")
		return err
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
	data := read[:n]
	fmt.Printf("the content of local file %v has changed\ndata is %v\n", _Arg.FilePath, string(data))
	return nil
}

func Search(_Arg *Arg.SearchArg) (offset int) {
	path := filepath.Join(LocalFilePath, _Arg.FilePath)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println("search file fail")
		return
	}
	data := make([]byte, 4096)
	n, err := f.Read(data)
	fmt.Println(string(data[:n]))
	offset = strings.Index(string(data[:n]), string(_Arg.Bytes)) //回车算两个字符
	return
}
