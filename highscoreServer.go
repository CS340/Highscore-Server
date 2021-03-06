package main

import (
	"net"
	"os"
	"fmt"
	"strings"
	"github.com/thoj/go-mysqlpure"
	"time"
	"bytes"
)

func main() {
	
	logIt("SETUP", "Starting...")
	
	addr, err := net.ResolveTCPAddr("ip4", ":4848")
	errorCheck(err, "Problem resolving TCP address")
	
	listen, err := net.ListenTCP("tcp", addr)
	errorCheck(err, "TCP listening error")
	
	logIt("SETUP", "Ready.")

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

	commm := parseCommand(string(buffer[0:]))
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

func parseCommand(com string) (string){
	com = string(bytes.TrimRight([]byte(com), string(byte(0))))
	com = strings.Replace(com, "\n", "", -1)
	var response string
	fmt.Println("COMMAND:"+com)
	dataCon, err := mysql.Connect("tcp", "127.0.0.1:3306", "hhss", "highscores", "hhss")
	errorCheck(err, "Could not connect to MySQL database.")
	
	scores := new(mysql.MySQLResponse)
	checker := new(mysql.MySQLResponse)

	parts := strings.Split(com, ":")
	switch parts[0]{
		case "user": 
			switch parts[1]{
				case "new":
					checker, err = dataCon.Query("SELECT username FROM users WHERE username='" + parts[2] + "';")
					if len(checker.FetchRowMap()) == 0{
						logIt("QUERY", "Inserting new user " + parts[2] + " " + parts[4] + ", " + parts[3])
						_, err = dataCon.Query("INSERT INTO users (firstName,lastName,username) VALUES('" + parts[3] + "', '" + parts[4] + "', '" + parts[2] + "')")
						if(err != nil){
							errorCheck(err, "Could not enter new user into database. USER: " + parts[2])
							response = "user:new:failure (reason unknown)\n"
						} else{
							response = "user:new:success\n"
						}

					} else {
						response = "user:exists\n";
					}
			}
		case "score":
			switch parts[1]{
				case "new":
					checker, err = dataCon.Query("SELECT username FROM users WHERE username='" + parts[2] + "';")
					if len(checker.FetchRowMap()) > 0 {
						logIt("QUERY", "Inserting new score of " + parts[3] + " for " + parts[2])
						_, err = dataCon.Query("INSERT INTO scores (username,score) VALUES('" + parts[2] + "', " + parts[3] + ")")
						errorCheck(err, "Could not enter score into database. USER: " + parts[2] + "SCORE:" + parts[3])
						response = "score:added\n"
					} else {
						response = "score:user does not exist\n";
					}
				case "get":
					parts[2] = string(bytes.TrimRight([]byte(parts[2]), string(byte(0))))
					if parts[2] == "all" { 
						logIt("QUERY", "Reading scores for " + parts[2])
						scores, err = dataCon.Query("SELECT * FROM scores ORDER BY score DESC")
						errorCheck(err, "Could not get scores from database. USER: " + parts[2])
						response = "score:all:"
						i := 0
						for row := scores.FetchRowMap(); row != nil && i < 10; row = scores.FetchRowMap() {
								response += row["username"] + "," + row["score"] + ";"
								i += 1
						}
						response += "\n"
						fmt.Println(response)
					} else {
						logIt("QUERY", "Reading scores for " + parts[2])
						scores, err = dataCon.Query("SELECT * FROM scores ORDER BY score DESC")
						errorCheck(err, "Could not get scores from database. USER: " + parts[2])
						//fmt.Println(len("ll"))
						response = "score:" + parts[2] + ":"
						i := 0
						for row := scores.FetchRowMap(); row != nil && i < 10; row = scores.FetchRowMap() {
							if row["username"] == parts[2] {
								response += row["username"] + "," + row["score"] + ";"
								i += 1
							}
						}
						response += "\n"
						fmt.Println(response)
					}
			}
	}
	dataCon.Quit();

	return response
}
