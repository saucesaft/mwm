package main

import (
	"fmt" 
	"log"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	xp "github.com/BurntSushi/xgb/xproto"
)

var keymap [256][]xp.Keysym

type keybind struct {
	keys string
	cmd []string
}

var (
	conn *xgb.Conn // connection to the x server

	testcmd = []string{"test1", "test2"}
	dummy = &keybind{"TESTKEY", testcmd}
//	kb1 = keybind{"mod4","shift","spawn xterm"}
	keybinds = []*keybind{dummy}
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func setup() {
	var err error
	conn, err = xgb.NewConn() // initialize connection
	if err != nil {
		panic(err)
	} else if conn == nil {
		fmt.Println("wtf, are you seriously trying to run a graphical app without a monitor!?")
	}

	info := xp.Setup(conn)
	
	if err := xinerama.Init(conn); err != nil {
		panic(err)
	}

	screen := info.DefaultScreen(conn) // get the screen
	root := screen.Root // assign root window
	
	err = xp.ChangeWindowAttributesChecked(
			conn,
			root,
			xp.CwEventMask,
			[]uint32{
				xp.EventMaskKeyPress|xp.EventMaskKeyRelease|xp.EventMaskButtonPress|
				xp.EventMaskButtonRelease|xp.EventMaskStructureNotify|xp.EventMaskSubstructureRedirect,}).Check()
		if err != nil {
			if _, ok := err.(xp.AccessError); ok {
				fmt.Println("Could not become the WM. Is another WM already running?")
				panic(err)
			}
	}

//	xp.ChangeKeyboardMapping(conn, byte(42), xp.Keycode(65), byte(23), []xp.Keysym{})

	const (
		loKey = 8
		hiKey = 255
	)

	m := xp.GetKeyboardMapping(conn, loKey, hiKey-loKey+1)

	reply, err := m.Reply()
	if err != nil {
		log.Fatal(err)
	}

	if reply == nil {
		log.Fatal("Could not load keyboard map")
	}

	for i := 0; i < hiKey-loKey+1; i++ {
		keymap[loKey+i] = reply.Keysyms[i*int(reply.KeysymsPerKeycode) : (i+1)*int(reply.KeysymsPerKeycode)]
	}

	fmt.Println(keymap)	

	go startIPC() // start ipc listener
	
	fmt.Println("setup done")	
}

func handleKeypress() {
	
}

func main() {
	setup()

	for {
			ev, xerr := conn.WaitForEvent()
			if ev == nil && xerr == nil {
				fmt.Println("Both event and error are nil. Exiting...")
				return
			}
			switch v := ev.(type){ // switch statement for X events

			case xp.KeyPressEvent:
				fmt.Println("Key pressed!", xp.Keysym(v.Detail))
			default:
				fmt.Println(v)
			}
			if xerr != nil {
				fmt.Printf("Error: %s\n", xerr)

			}
		}
	
}
