//This exist to test the server without writing a full frontend

package main

import (
	"os"
	"fmt"
	"net"
	"encoding/json"
	"io"
)

const PORT = "4591"

type message struct {
	Author int `json:"author"`
	Text string `json:"text"`
}

func (m message) String() string {
	return fmt.Sprintf("%s :: %s", m.Author, m.Text)
}

func mkIP () string, error{
	var e error
	if len(os.Args) > 1 {
		if os.Args[1] == "-p" {
			if len(os.Args) > 3 {
				ip := os.Args[2]
				return ip + ":" + PORT, nil
			} else {
				fmt.Println("Are you stupid? You didn't give an ip address after -p")
				return nil, 
			}
		} else {
			fmt.Println("invalid argument ", os.Args[1])
		}
	}
	return "localhost:" + PORT
}

func checkErr (e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func recieve (conn net.Conn) error {
	defer conn.Close()
	fmt.Println("Got a message!")
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
	fmt.Println(string(toRet))
	return nil
}

func main () {
	ipAddr := mkIP
	for {
		fmt.Println("Connecting to ", )
		conn, err := net.Dial("tcp", )
		checkErr(err)
		go recieve(conn)
		msg := &message{Author: 0, Text: "This is a test"}
		byt, err := json.Marshal(msg)
		checkErr(err)
		conn.Write(byt)
		fmt.Println("Sent :^) \n", msg)
	}
}
