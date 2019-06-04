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
	"fmt"
	"io"
	"encoding/json"
	"net"
	"./msg"
)

const PORT = ":4591" //the port at which connections can be made to the server


//Generalized error checker, panics
func checkErr (e error) {
	if e != nil {
		panic(e)
	}
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
	return toRet, nil
}

func displayMsg(input []byte) {
	var m msg.Message
	err := json.Unmarshal(input, &m);
	checkErr(err)
	fmt.Println(m)
}

func main () {
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
