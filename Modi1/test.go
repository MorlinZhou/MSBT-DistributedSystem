package main

type rpc struct {
	Type int
	Name string

	Offset int
}

func main() {
	//bytes := []byte(string("\\CallType LookUp \\FilePath file/check.txt \\Offset 0 \\Number 10"))
	//r := rpc{
	//	Type:   10,
	//	Name:   "1",
	//	Offset: 111110,
	//}
	//value := reflect.ValueOf(r)
	//fmt.Println(value)
	//fmt.Println(value.Field(1))
	//fmt.Println(value.NumField())
}
