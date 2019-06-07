//A server for chatting and file sharing between various friends. Does not require cloud computing or a large company to run.
/*
//STARTED - 5/23/19
//LAST TOUCHED - 6/3/19
FEATURE LIST
5/23/19 >> Send and receive messages in json format through tcp ports
(STARTED) >> Save messages in a database
TODO >> Secure and sign messages and user accounts using RSA
TODO >> Multimedia embedding into a message
TODO >> Files saved to the server are kept in a persistent repository
TODO >> Servers can be private and require users to be invited
TODO >> Users can have various permissions and roles
TODO >> Servers can push updates about activity to members who aren't currently connected
--------At this point start working on the front end-------------------------------------
*/

package main

import (
	"os"
	"fmt"
	"io"
	"encoding/json"
	"net"
	"./msg"
	"./db"
	"strconv"
)

const PORT = ":4591" //the port at which connections can be made to the server


//Generalized error checker, panics
func checkErr (e error) {
	if e != nil {
		panic(e)
	}
}

func runTests() {
	//test messages
	tstmsg := msg.Message{Author: 0, Text: "This is a test"}
	fmt.Println(tstmsg)
	fmt.Println(tstmsg.ToJson())
	jsn, err := tstmsg.ToJson()
	checkErr(err)
	fmt.Println(string(jsn))
	newmsg := msg.Message{Author: 1, Text: "Ree"}
	jsn, err = newmsg.ToJson()
	checkErr(err)
	fmt.Println(tstmsg.FromJson(jsn))

	//test database
	err = db.AddMsg(tstmsg)
	checkErr(err)
	for i:=0; i<19; i++ {
		tstmsg := msg.Message{Author:2 , Text: "This is message " + strconv.Itoa(i)}
		err = db.AddMsg(tstmsg)
		checkErr(err)
	}
	var holder []msg.Message
	db.Last20(&holder)
	fmt.Println(holder)
}

func handleConn(conn net.Conn) ([]byte, error){
	//defer conn.Close()
	fmt.Println("Got a connection!")
	buffer := make([]byte, 8)
	toRet := make([]byte, 0, 256)
	for {
		bts, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			fmt.Println("bytes written from buffer: ", bts)
			toRet = append(toRet, buffer[:bts]...)
		}
	}
	fmt.Println("Bytes writern: ", len(toRet), " Raw Data: \n-----------------------------------------\n", string(toRet))
	conn.Close()
	var m msg.Message
	m.FromJson(toRet)
	db.AddMsg(m)
	return toRet, nil
}

func displayMsg(input []byte) {
	var m msg.Message
	err := json.Unmarshal(input, &m);
	checkErr(err)
	fmt.Println(m)
}

func main () {
	args := os.Args

	if args[1] == "-t" {
		runTests()
	}

	fmt.Println("Starting server")
	ln, err := net.Listen("tcp", PORT)
	checkErr(err)
	fmt.Println("Listening on port ", PORT)
	for {
		conn, err := ln.Accept()
		checkErr(err)
		go handleConn(conn)
	}
}
