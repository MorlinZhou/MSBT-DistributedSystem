package LocalFunc

import (
	"DSWork/Client/Arg"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func SaveToLocal(Filepath string, data []byte) {
	path := filepath.Join("E:/DSWork/Client", Filepath)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
	path := filepath.Join("E:/DSWork/Client", Filepath)
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
	path := filepath.Join("E:/DSWork/Client", _Arg.FilePath)
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
	path := filepath.Join("E:/DSWork/Client", _Arg.FilePath)
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
