package Arg

//code bt zhh
import (
	"bytes"
	"fmt"
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
