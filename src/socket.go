package main

import (
    "fmt"
    "net"
    "encoding/gob"
	"os/exec"
	"log"
)

type P struct {
    M string
}

func runCommand(cmd *P) {
	hey := exec.Command(cmd.M)
	if err := hey.Run(); err != nil {
		log.Fatal(err)
	}	
}

func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &P{}
    dec.Decode(p)
//    runCommand(p)
	fmt.Println(p.M)
    conn.Close()
}

func startIPC() {
    fmt.Println("ipc started");
   ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        // handle error
    }
    for {
        conn, err := ln.Accept() // this blocks until connection or error
        if err != nil {
            // handle error
            continue
        }
        go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
    }
}