//This exist to test the server without writing a full frontend

package main

import (
	"fmt"
	"net"
	"encoding/json"
	"io"
)

const IPADDR = "localhost:4591"

type message struct {
	Author int `json:"author"`
	Text string `json:"text"`
}

func (m message) String() string {
	return fmt.Sprintf("%s :: %s", m.Author, m.Text)
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
	for {
		fmt.Println("Connecting to ", IPADDR)
		conn, err := net.Dial("tcp", IPADDR)
		checkErr(err)
		go recieve(conn)
		msg := &message{Author: 0, Text: "This is a test"}
		byt, err := json.Marshal(msg)
		checkErr(err)
		conn.Write(byt)
		fmt.Println("Sent :^) \n", msg)
	}
}
