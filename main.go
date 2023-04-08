package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	wc   sync.WaitGroup
	from string
	to   string
)

func Usage() {
	flag.StringVar(&from, "f", "", "From IP:PORT")
	flag.StringVar(&to, "t", "", "To IP:PORT")
	flag.Parse()
	if from == "" || to == "" {
		Head()
		os.Exit(1)
	}

}

func Head() {
	fmt.Println(`
-------------------------------------------------------------------------
Example:goport -f 127.0.0.1:8000 -t 127.0.0.1:9000
-------------------------------------------------------------------------																							  
	`)
}

func handleConnection(r, w net.Conn) {
	defer r.Close()
	defer w.Close()
	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}
		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}

func main() {
	Usage()
	fromaddr := fmt.Sprintf("%s", from)
	toaddr := fmt.Sprintf("%s", to)

	fromListener, err := net.Listen("tcp", fromaddr)
	if err != nil {
		log.Fatal("Unable to listen on: %s, error: %s\n", fromaddr, err.Error())
	}

	for {
		fromconn, err := fromListener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a request, error: %s\n", err.Error())
		} else {
			fmt.Println("new connect:" + fromconn.RemoteAddr().String())
		}

		toconn, err := net.Dial("tcp", toaddr)
		if err != nil {
			fmt.Printf("can not connect to %s\n", toaddr)
			continue
		}
		go handleConnection(fromconn, toconn)
		go handleConnection(toconn, fromconn)

	}

}
