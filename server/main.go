package main

import (
	"io"
	"os"
	"fmt"
	"net"
	"time"
	"bufio"
	//"encoding/json"
)

//CONSTANTS
const PORT = "localhost:4591"
const SILOS = "datapile/silos"
const MEDIA = "datapile/media"
const LOCKED = "datapile/media/locked"

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

//Checks if a path and/or file exists and makes it if it doesn't. Only used for critical structure
func checkCritical (path string) error{
	fmt.Println("Checking directory ", path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.MkdirAll(path, 0666)
		if err != nil {
			fmt.Println("criticalCheck of ", path, " failed.")
			return err
		}
		fmt.Println("Check OK!")
		reutrn nil
	}
}

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

//Ran on each startup. makes sure if critical structure is valid
func preStartup () {
	fmt.Println("PRE-STARTUP BEHAVIOR BEGINNING")
	var paths := [SILOS, MEDIA, LOCKED]
	for _, str := range(paths) {
		err := ckeckCritical(str)
		if err != nil {
			panic(err)
		}
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

//Awaiting connections
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

//Main loop
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
