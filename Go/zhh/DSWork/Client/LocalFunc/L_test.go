package LocalFunc

import (
	"DSWork/Client/Arg"
	"fmt"
	"testing"
)

func TestSaveToLocal(t *testing.T) {
	SaveToLocal("file/check3.txt", []byte("abcdefghijk"))
}

func TestReadFromFile(t *testing.T) {
	_Arg := Arg.NewLookArg("file/check1.txt", 2, 2)
	ReadFromFile(_Arg)
}

func TestInsertFile(t *testing.T) {
	_Arg := Arg.NewInsertArg("file/check1.txt", 2, []byte{'a', 'b'})
	InsertFile(_Arg)
}

func TestCheckLocalFile(t *testing.T) {
	fmt.Println(CheckLocalFile("file/check1.txt"))
	fmt.Println(CheckLocalFile("file/check2.txt"))
}
