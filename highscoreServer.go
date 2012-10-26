package main

import (
	"net"
	"os"
	"fmt"
	"strings"
	"github.com/thoj/Go-MySQL-Client-Library"
	"time"
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
		logIt("CONNECTION", "Got new connection")
		
		go newClient(connection)
		
	}

	os.Exit(0)
}

func newClient(connect net.Conn){
	logIt("CONNECTION", "Handling new client")
	var buffer [512]byte

	_, err := connect.Read(buffer[0:])
	if err != nil {
		logError("ERROR", "Error reading from client", err)
		connect.Close()
		return
	}

	commm, _ := parseCommand(string(buffer[0:]))
	_, err2 := connect.Write([]byte(commm))
	if err2 != nil {
		logError("ERROR", "Error writing to client", err2)
		connect.Close()
		return
	}
	connect.Close()
	logIt("CONNECTION", "Closing connection to client")
}

func logError(ertype string, message string, err error){
	fmt.Printf("%s\t%s: %s: %s\n", time.Now().String(), ertype, message, err)
}

func logIt(ertype string, message string){
	fmt.Printf("%s\t%s: %s\n", time.Now().String(), ertype, message)
}

func errorCheck(err error, message string){
	if(err != nil){
		logError("ERROR", message, err)
	}
}

func parseCommand(com string) (string, int){

	dataCon, err := mysql.Connect("tcp", "127.0.0.1:3306", "highscores", "hhss", "highscores")
	errorCheck(err, "Could not connect to MySQL database.")
	var out int
	out = 1
	
	scores := new(mysql.MySQLResponse)

	parts := strings.Split(com, ":")

	switch parts[0]{
		case "user": 
			switch parts[1]{
				case "new":
					logIt("MySQL:Query", "Inserting new user " + parts[2] + " " + parts[3] + ", " + parts[4])
					_, err = dataCon.Query("INSERT INTO users (firstName,lastName,username) VALUES('" + parts[2] + "', " + parts[3] + "', " + parts[4] + ")")
					errorCheck(err, "Could not enter new user into database. USER: " + parts[2])
			}
		case "score":
			switch parts[1]{
				case "new":
					logIt("MySQL:Query", "Inserting new score of " + parts[3] + " for " + parts[2])
					_, err = dataCon.Query("INSERT INTO scores (username,score) VALUES('" + parts[2] + "', " + parts[3] + ")")
					errorCheck(err, "Could not enter score into database. USER: " + parts[2] + "SCORE:" + parts[3])
				case "get":
					logIt("MySQL:Query", "Reading scores for " + parts[2])
					fmt.Println("SELECT * FROM scores WHERE `username` = 'aphelps' ORDER BY score DESC")
					scores, err = dataCon.Query("SELECT * FROM scores WHERE username = " + parts[2] + " ORDER BY score DESC")
					errorCheck(err, "Could not get scores from database. USER: " + parts[2])
					fmt.Println(scores.FetchRowMap())
			}
	}
	dataCon.Quit();

	return com, out
}