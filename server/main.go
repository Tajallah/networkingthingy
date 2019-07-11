//A server for chatting and file sharing between various friends. Does not require cloud computing or a large company to run.
/*
//STARTED - 5/23/1159
//LAST TOUCHED - 7/119
FEATURE LIST
5/23/19 >> Send and receive messages in json format through tcp ports
(STARTED) >> Save messages in a database
	*save singular messages
	-Let's keep this light. No Database is actually needed. We'll justy store the posts as loose json files 
	-serve saved messages on request
	-when a new message is commited to db it must be autoserved to all clients
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

//this guy will send all messages from the publish stack to all clients later
func publish (publishStack chan []byte, clients []net.Conn) error {
	toWrite := <-publishStack
	for _, conn := range(clients) {
		bts, err := conn.Write(toWrite)
		if err != nil {
			return err
		} else {
			fmt.Println("Wrote ", bts,"to connection at ", conn)
		}
	}
	return nil
}

func handleConn(conn net.Conn, publishStack chan []byte) error {
	defer conn.Close()
	fmt.Println("Got a connection!")
	buffer := make([]byte, 8)
	toRet := make([]byte, 0, 256)
	for {
		bts, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			fmt.Println("bytes written from buffer: ", bts)
			toRet = append(toRet, buffer[:bts]...)
		}
	}
	//fmt.Println("Bytes writern: ", len(toRet), " Raw Data: \n-----------------------------------------\n", string(toRet))
	//conn.Close()
	publishStack <- toRet
	var m msg.Message
	m.FromJson(toRet)
	db.AddMsg(m)
	return nil

func displayMsg(input []byte) {
	var m msg.Message
	err := json.Unmarshal(input, &m);
	checkErr(err)
	fmt.Println(m)
}

func main () {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "-t" {
			runTests()
		}
	}

	var clients []net.Conn
	publishStack := make(chan []byte)

	fmt.Println("Starting server")
	ln, err := net.Listen("tcp", PORT)
	checkErr(err)
	fmt.Println("Listening on port ", PORT)
	for {
		conn, err := ln.Accept()
		checkErr(err)
		clients = append(clients, conn)
		for _, cli := range(clients) {
			go handleConn(cli, publishStack)
		}
		go publish(publishStack, clients)
	}
}
