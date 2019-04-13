package main

import (
	"fmt"
	"strings"
	"encoding/gob"
	"github.com/devfacet/gocmd"
	"net"
	"log"
)

type P struct {
    M string
}

func main() {
	fmt.Println("start client")
	
	flags := struct {
		Help      bool `short:"h" long:"help" description:"Display usage" global:"true"`
		Version   bool `short:"v" long:"version" description:"Display version"`
		VersionEx bool `long:"vv" description:"Display version (extended)"`
		Bind      struct {
			Settings bool `settings:"true" allow-unknown-arg:"true"`
		} `command:"bind" description:"Print arguments"`
	}{}	

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("connection error", err)
	}	
	
	encoder := gob.NewEncoder(conn)
	
	// Echo command
	gocmd.HandleFlag("Bind", func(cmd *gocmd.Cmd, args []string) error {
		p := &P{strings.Join(cmd.FlagArgs("Bind")[1:], "/")}
		encoder.Encode(p)
		conn.Close()
		
		return nil
	})

	// Init the app
	gocmd.New(gocmd.Options{
		Name:        "hydro",
		Version:     "0.0.0",
		Description: "cli client for mwm",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})


	conn.Close()
	fmt.Println("done")
}
