package main

import (
	"fmt" 
	"log"
	"os/exec"
	
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	xp "github.com/BurntSushi/xgb/xproto"
)

var conn *xgb.Conn // connection to the x server
var screen xp.ScreenInfo // principal screen
var root  xp.Window // root window

var (
	atomWMProtocols xp.Atom
	atomWMTakeFocus xp.Atom
)

var defaultw = &Workspace{screen: screen} // create the first workspace
var workspaces = []*Workspace{defaultw} // create a list of workspaces

func init() {
	log.SetFlags(log.Lshortfile)
}

/// X Event Handlers ///

func handleConfigureRequest(e xp.ConfigureRequestEvent) {
	fmt.Println("Configure request baby")
	wc := xp.ConfigureNotifyEvent{
		Event: e.Window, Window: e.Window, AboveSibling: 0, X: e.X, Y: e.Y,
		Width: e.Width, Height: e.Height, BorderWidth:0, OverrideRedirect: false,
	}
	xp.SendEventChecked(conn, false, e.Window, xp.EventMaskStructureNotify, string(wc.Bytes()))

	//xp.ConfigureWindow(conn, e.Window, e.value_mask, wc)
}

func handleMapRequest(e xp.MapRequestEvent) {
	fmt.Println("Map request baby")
	w := workspaces[0]
	w.addWin(e.Window)
	w.manage()
	xp.MapWindow(conn, e.Window)
}

func handleDestroyNotify(e xp.Window) {
	fmt.Println("Window destroyed!")
	index := -1
	for _, w := range workspaces{
		for i, v := range w.clients {
			if v.window == e {
				index = i
				break
			}
		}
		if index != -1 {
			for _, w := range workspaces {
				w.clients[index] = w.clients[len(w.clients)-1]
				w.clients = w.clients[:len(w.clients)-1]
			}
		}
	}
}

func handlerEnterNotify( e xp.EnterNotifyEvent) {
	activeWindow = &e.Event
	prop, err := xp.GetProperty(conn, false, e.Event, atomWMProtocols, xp.GetPropertyTypeAny, 0, 64).Reply()
	focused := false
	if err  == nil {
		TakeFocusPropLoop:
		for x := prop.Value; len(x) >= 4; x = x[4:] {
			switch xp.Atom(uint32(x[0]) | uint32(x[1])<<8 | uint32(x[2])<<16 | uint32(x[3])<<24) {
			case atomWMTakeFocus:
				xp.SendEventChecked(conn, false, e.Event, xp.EventMaskNoEvent, string(xp.ClientMessageEvent{
					Format: 32, Window: *activeWindow, Type: atomWMProtocols, Data: xp.ClientMessageDataUnionData32New([]uint32{
						uint32(atomWMTakeFocus), uint32(e.Time), 0, 0, 0,}),
				}.Bytes())).Check()
			focused = true
			break TakeFocusPropLoop
			}
		}
	}
	if !focused {
		if _, err := xp.SetInputFocusChecked(conn, xp.InputFocusPointerRoot, e.Event, e.Time).Reply(); err != nil {
//			fmt.Println("baddddd")
//			panic(err)
		}
	}
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

	fmt.Printf("%T\n", root)

	atomWMProtocols = getAtom("WM_PROTOCOLS")
	
	err = xp.ChangeWindowAttributesChecked( // listen for specific events
			conn,
			root,
			xp.CwEventMask,
			[]uint32{
				xp.EventMaskButtonPress|xp.EventMaskButtonRelease|xp.EventMaskStructureNotify|xp.EventMaskSubstructureRedirect,}).Check()
		if err != nil {
			if _, ok := err.(xp.AccessError); ok {
				fmt.Println("Could not become the WM. Is another WM already running?")
				panic(err)
			}
	}

	go startIPC() // start ipc listener
	
	fmt.Println("setup done")	
}

func getAtom (name string) xp.Atom {
	rply, err := xp.InternAtom(conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		fmt.Println(err)
	}
	if rply == nil {
		return 0
	}
	return rply.Atom
}

func main() {
	setup()

	cmd := exec.Command("/bin/bash", "/home/eduarch/repos/mwm/examples/mwmrc")
	fmt.Println("running config...")
	go cmd.Run()

	for {
			ev, xerr := conn.WaitForEvent()
			if ev == nil && xerr == nil {
				fmt.Println("Both event and error are nil. Exiting...")
				return
			}
			switch v := ev.(type){ // switch statement for X events

			case xp.ConfigureRequestEvent: // window wants to have a size and position asssigned
				handleConfigureRequest(v)
			case xp.MapRequestEvent: // window wants to be shown
				handleMapRequest(v)
			case xp.DestroyNotifyEvent:
				handleDestroyNotify(v.Window)
				for _, w := range workspaces{
					w.manage()
				}
				_, err := xp.SetInputFocusChecked(conn, xp.InputFocusPointerRoot, root, xp.TimeCurrentTime).Reply()
				if err != nil {
					fmt.Println(err)
				}
			case xp.EnterNotifyEvent:
				handlerEnterNotify(v)
			default:
				fmt.Println(ev)
			}
			if xerr != nil {
				fmt.Printf("Error: %s\n", xerr)

			}
		}
	
}