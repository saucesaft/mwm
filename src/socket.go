package main

import (
    "fmt"
    "net"
    "encoding/gob"
	"os/exec"
	"log"
	s"strings"
)

type P struct { // special struct for sending things
    M string
}

func addKey(cmd *P) { // adds keybinds to struct
	complete := s.Split(cmd.M, "/")
	newKeybind := &keybind{complete[0], complete[1:]}
	keybinds = append(keybinds, newKeybind)
}

func runCommand(cmd *P) { // runs shell commands
	hey := exec.Command(cmd.M)
	if err := hey.Run(); err != nil {
		log.Fatal(err)
	}	
}

func handleConnection(conn net.Conn) { // ran every time needs to send something to mwm
    dec := gob.NewDecoder(conn)
    p := &P{}
    dec.Decode(p)
	addKey(p)
    conn.Close()
}

func startIPC() { // start the ipc listener func
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