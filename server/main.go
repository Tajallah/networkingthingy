package main

import (
	"io"
	"fmt"
	"net"
	"time"
	"bufio"
	//"encoding/json"
)

//CONSTANTS
const PORT = "localhost:4591"

//Global Variables
var (
	awaitedConns = make(chan net.Conn)
	connections = make(map[string]bool) // This is a list of the network addresses of clients and whether or not they're live
	publishStack = make(chan []byte)
	cleanerSwitch = true
)

//generalized error checker
func checkErr (e error) {
	if e != nil {
		fmt.Println(e) //for now we'll just print the error
	}
}

//this exists so that we can check if a connection is in our list of connections
/*func connIter (c net.Conn) bool {
	fmt.Println("---CONN ITER---")
	for conn, _ := range (connections) {
		if conn == c {
			return true
		}
	}
	return false
}*/

/*func connRemove (c net.Conn) string {
	fmt.Println("---REMOVING---")
	for i, conn := range(connections) {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			return "---CONNECTION TO " + c.RemoteAddr().String() + " LOST---"
		}
	}
	return "---TRIED TO REMOVE " + c.RemoteAddr().String() + " BUT IT CANNOT BE FOUND IN connections---"
}*/

//sends out a message on the publish stack to every open connection
func broadcast (msg []byte) {
	fmt.Println("---BROADCASTING---")
	for conn, active := range (connections) {
		if active == true{
			c, err := net.Dial("tcp", conn)
			defer c.Close()
			checkErr(err)
			if err != nil {
				c.Write(msg)
				fmt.Println("---SENT ", string(msg), " to ", conn, "---")
			} else {
				connections[conn] = false
			}
		}
	}
}

//Check to make sure that clients are still alive
func touch () {
	fmt.Println("---TOUCHING---")
	for conn, active := range(connections) {
		if active == true {
			c, err := net.Dial("tcp", conn)
			checkErr(err)
			if err != nil {
				c.Write([]byte("🆗"))
			} else {
				connections[conn] = false
			}
		}
	}
}

//clean false connections
func cleaner () error{
	if cleanerSwitch != false {
		return nil
	}
	for conn, active := range(connections) {
		if active == false {
			delete(connections, conn)
		}
	}
	cleanerSwitch = false
	return nil
}

func cleanerTimer() {
	timing := 3 * time.Second
	for {
		cleanerSwitch = true
		time.Sleep(timing)
	}
}

//Handling an incoming connection
func handleConn (conn net.Conn) ([]byte, error) {
	fmt.Println("---HANDLING CONNECTION---")
	addr := conn.RemoteAddr().String()
	fmt.Println("---GOT A CONNECTION FROM ", addr, "---")
	defer conn.Close()
	buffer := make([]byte, 256)
	fnl := make([]byte, 16)
	rdr := bufio.NewReader(conn)
	for {
		bytesnum, err := rdr.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return nil, err
			}
		} else if bytesnum != 0{
			fmt.Println("--recieved ", bytesnum, " bytes --")
			fnl = append(fnl, buffer...)
		} else {
			break
		}
	}
	fmt.Println(string(fnl))
	publishStack <- fnl
	fmt.Println(fnl)
	connections[addr] = true
	return fnl, nil
}

func awaitConns (ln net.Listener) {
	for {
		conn, err := ln.Accept()
		checkErr(err)
		defer conn.Close()
		deadline := 5 * time.Second
		conn.SetDeadline(time.Now().Add(deadline))
		awaitedConns <- conn
	}
}

func main () {
	fmt.Println("---STARTING SERVER---")
	go cleanerTimer()
	var msg []byte
	var conn net.Conn
	ln, err :=  net.Listen("tcp", PORT)
	checkErr(err)
	go awaitConns(ln)
	for {
		msg = nil
		conn = nil
		select{
		case conn = <-awaitedConns:
			handleConn(conn)
		case msg = <-publishStack:
			broadcast(msg)
		}
		if len(connections) > 0 {
			touch()
			cleaner()
		}
	}
}
