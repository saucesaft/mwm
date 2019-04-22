package main

import (
"fmt" 
//	"log"
//	"os/exec"
	
//	"github.com/BurntSushi/xgb"
	//"github.com/BurntSushi/xgb/xinerama"
	xp "github.com/BurntSushi/xgb/xproto"
)

type Frame struct {
//	clients []client
	x,y int16
	w,h uint16
}

type Client struct {
	window xp.Window
	frame Frame
}

type Workspace struct {
	screen xp.ScreenInfo
	clients []*Client
}

var activeWindow *xp.Window

func (w *Workspace) addWin(win xp.Window) error {
	if err := xp.ConfigureWindowChecked( // check if we can manage the window
		conn,
		win,
		xp.ConfigWindowBorderWidth,
		[]uint32{
			0,
		}).Check(); err != nil {
		return err}
	if err := xp.ChangeWindowAttributesChecked( // know when the window is destroyed
		conn,
		win,
		xp.CwEventMask,
		[]uint32{
			xp.EventMaskStructureNotify |
			xp.EventMaskEnterWindow,
		}).Check(); err != nil {
		return err
	}

	geom, err := xp.GetGeometry(conn, xp.Drawable(win)).Reply() // get the geometry of the window
	if err != nil {
		panic(err)
	}

	f :=  Frame{geom. X, geom.Y, geom.Width, geom.Height} // create a frame with the dimensions of the window
	c := &Client{win, f} // create a client struct with the window inside
	
	w.clients = append(w.clients, c)

//	fmt.Println(geom.Width)

	return nil
}

func (w *Workspace) manage() {
	for _, v := range w.clients {
		geom, err := xp.GetGeometry(conn, xp.Drawable(v.window)).Reply()
		if err != nil {
//			fmt.Println("baddddd")
//			panic(err)
		}

		fmt.Println(geom.Width)

		configerr := xp.ConfigureWindowChecked(conn, v.window, xp.ConfigWindowX|
			xp.ConfigWindowY|xp.ConfigWindowWidth|xp.ConfigWindowHeight,
			[]uint32{uint32(geom.X),uint32(geom.Y),uint32(geom.Width),uint32(geom.Height)}).Check()
		if configerr != nil {
			panic(configerr)
		}
//		fmt.Println(i.window)
	}
}