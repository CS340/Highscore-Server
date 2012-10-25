package main

import (
	"net"
	"os"
	"fmt"
	//"io/ioutil"
	"time"
	//"strconv"
)

func main() {

	addr, err := net.ResolveTCPAddr("ip4", ":4848")
	errorCheck(err, "Problem resolving TCP address")
	
	listen, err := net.ListenTCP("tcp", addr)
	errorCheck(err, "TCP listening error")

	for{
		connection, err := listen.Accept()
		if(err != nil){
			continue
		}
		
		//connection.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
		newClient(connection)
		connection.Close()
	}

	os.Exit(0)
}

func newClient(connect net.Conn){

	var buffer [512]byte

	for{
		n, err := connect.Read(buffer[0:])
		if err != nil {
			return
		}

		fmt.Println(string(buffer[0:]))
		_, err2 := connect.Write(buffer[0:n])
		if err2 != nil {
			return
		}
	}
}

func logIt(ertype string, message string, err error){
	fmt.Printf("%s\t%s: %s: %s", time.Now().String(), ertype, message, err)
}

func errorCheck(err error, message string){
	if(err != nil){
		logIt("ERROR", message, err)
	}
}