package main

import (
    "fmt"
    "log"
	"net"
	"net/http"
	"io"
)

func main() {

	go func(){
		fmt.Println("Start Pac server..")
		http.ListenAndServe(":7072", http.FileServer(http.Dir(".")))
		fmt.Println("End Pac server..")
	}()
	fmt.Println("Start Sock5 server..")
	var fromport, toport int = 7070, 1086
	fromaddr := fmt.Sprintf("0.0.0.0:%d", fromport)
	toaddr := fmt.Sprintf("127.0.0.1:%d", toport)
	
	fromlistener, err := net.Listen("tcp", fromaddr)

	if err != nil {
		log.Fatal("Unable to listen on: %s, error: %s\n", fromaddr, err.Error())
	}
	defer fromlistener.Close()

	for {
		//fmt.Println("wait for  connect:")
		fromcon, err := fromlistener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a request, error: %s\n", err.Error())
		} else {
			//fmt.Println("new connect:" + fromcon.RemoteAddr().String())
		}

		go handleConnection(fromcon, toaddr)
	}

}

func handleConnection(fromcon net.Conn, toaddr string) {

	defer fromcon.Close()
	toCon, err := net.Dial("tcp", toaddr)
	if err != nil {
		fmt.Printf("can not connect to %s\n", toaddr)
		return
	}
	defer toCon.Close()

	done := make(chan int)
	/*
	var handler = func(r,w net.Conn){
		defer func(){ ch <- 0 }()
		var buffer = make([]byte, 4096)
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
	}*/

	var handler = func(r,w net.Conn){
		defer func(){ done <- 0 }()
		buf := make([]byte, 4096)
		io.CopyBuffer(r, w, buf)
	}

	go handler(toCon,fromcon)
	go handler(fromcon,toCon)
	<- done
}
