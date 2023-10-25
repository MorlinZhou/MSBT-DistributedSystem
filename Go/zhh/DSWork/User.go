package main

import (
	"DSWork/Client"
	"DSWork/Client/Arg"
	"fmt"
	"log"
	"net"
)

func main() {
	var input_interval string
	var input_operation string
	var input_file_path string
	var Offset, Number int
	var Bytes []byte
	localIp := net.ParseIP("127.0.0.1")
	Port := 30000

	fmt.Print("Please set up a valid freshness interval before start (in seconds). ")
	_, err1 := fmt.Scanf("%s", &input_interval)
	if err1 != nil {
		fmt.Println("invalid interval:", err1)
		return
	}
	fmt.Println("", input_interval)
	fmt.Scanf("%s", new(string))

	fmt.Print("Choose operation (LOOKUP/INSERT/MONITOR/SEARCH/APPEND): ")

	_, err2 := fmt.Scanf("%s", &input_operation)
	if err2 != nil {
		fmt.Println("invalid operation:", err2)
		return
	}
	validOperations := map[string]bool{
		"LOOKUP":  true,
		"INSERT":  true,
		"MONITOR": true,
		"SEARCH":  true,
		"APPEND":  true,
	}
	if !validOperations[input_operation] {
		fmt.Println("INVALID OPERATION:", input_operation)
		return
	}
	fmt.Println("", input_operation)
	fmt.Scanf("%s", new(string))

	fmt.Print("Enter file path: such as file/check1.txt: ")
	_, err3 := fmt.Scanf("%s", &input_file_path)
	if err3 != nil {
		fmt.Println("invalid file_path:", err3)
		return
	}
	fmt.Scanf("%s", new(string))
	fmt.Println("", input_file_path)

	switch input_operation {
	case "LOOKUP":
		fmt.Println("Performing LOOKUP operation.")
		fmt.Println("Please input offset and number (such as 2,2): ")
		_, err_1 := fmt.Scanf("%d,%d", &Offset, &Number)
		if err_1 != nil {
			fmt.Println("invalid offset or number:", err_1)
			return
		}
		fmt.Scanf("%s", new(string))

		_Arg := Arg.NewLookArg(input_file_path, Offset, Number)

		Client.BootClient(localIp, Port, _Arg)

	case "INSERT":
		var s string
		fmt.Println("Performing INSERT operation.")
		fmt.Println("Please input offset and Bytes (such as 2,string): ")
		_, err_2 := fmt.Scanf("%d,%s", &Offset, &s)
		Bytes = []byte(s)
		if err_2 != nil {
			fmt.Println("invalid offset or Bytes:", err_2)
			return
		}
		fmt.Scanf("%s", new(string))

		_Arg := Arg.NewInsertArg(input_file_path, Offset, Bytes)
		Client.BootClient(localIp, Port, _Arg)

	case "MONITOR":
		var MonitorInterval int
		fmt.Println("Performing MONITOR operation.")
		fmt.Println("Please input MonitorInterval int (in seconds): ")
		_, err_3 := fmt.Scanf("%d", &MonitorInterval)
		if err_3 != nil {
			fmt.Println("invalid MonitorInterval:", err_3)
			return
		}
		fmt.Scanf("%s", new(string))
		_Arg := Arg.NewMonitorArg(input_file_path, MonitorInterval)

		log.Println("client send monitor msg")
		Client.BootClient(localIp, Port, _Arg)

	case "SEARCH":
		var s string
		fmt.Println("Performing SEARCH operation.")
		fmt.Println("Please input Bytes (such as string): ")
		_, err_4 := fmt.Scanf("%s", &s)
		Bytes = []byte(s)
		if err_4 != nil {
			fmt.Println("invalid Bytes:", err_4)
			return
		}
		fmt.Scanf("%s", new(string))

		_Arg := Arg.NewSearchArg(input_file_path, Bytes)
		Client.BootClient(localIp, Port, _Arg)

	case "APPEND":
		var s string
		fmt.Println("Performing APPEND operation.")
		fmt.Println("Please input Bytes (such as string): ")
		_, err_5 := fmt.Scanf("%s", &s)
		Bytes = []byte(s)
		if err_5 != nil {
			fmt.Println("invalid Bytes:", err_5)
			return
		}
		fmt.Scanf("%s", new(string))

		_Arg := Arg.NewAppendArg(input_file_path, Bytes)
		Client.BootClient(localIp, Port, _Arg)

	}
}
