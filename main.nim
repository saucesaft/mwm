import strutils, strmisc

import os, threadpool

import xlib, x

converter toCint(x: TKeyCode): cint = x.cint

converter int32toCint(x: int32): cint = x.cint

converter int32toCUint(x: int32): cuint = x.cuint

converter toTBool(x: bool): TBool = x.TBool

converter toBool(x: TBool): bool = x.bool

type
  keybind = tuple[mods, keys, cmd: string]
  client = tuple[x, y, w, h: int, win: TWindow]

var
  dummy: keybind = (mods: "", keys: "", cmd: "")
  keybinds = @[dummy]

var

  display:PDisplay

  root:TWindow

  attr:TXWindowAttributes

  start:TXButtonEvent

  ev:TXEvent

# parses the config file
proc parseConfig =

  for line in lines "config.ini":
    case line
    of "[keybinds]":
      continue
    var 
      newline = line.partition("+")
      newMods = newline[0].strip() #Add a new way to parse the extra mod, like shift
      afterline = newline[2].partition("=")
      newKeys = afterline[0].strip()
      newCmd = afterline[2].strip()
      newKeybind: keybind = (mods: newMods, keys: newKeys, cmd: newCmd) # tuple
      
    keybinds.add(newKeybind)

# grab keypress
proc grabKeypress =

  for i in keybinds:
    var newMod: int
    case i.mods
    of "meta": newMod = Mod4Mask
    of "alt": newMod = Mod1Mask
    discard XGrabKey(display, cast[cint](XKeysymToKeycode(display, XStringToKeysym(i.keys))), cast[cuint](newMod), root,
      cast[TBool](true), cast[cint](GrabModeAsync), cast[cint](GrabModeAsync))   
    
# parse the command
proc parseCmd(command: string) =
  #check if spawns something
  var parsedCmd = command.partition("spawn")
  if parsedCmd[1] == "spawn":
    discard spawn execShellCmd(parsedCmd[2])

# parse the keypresses
proc handleKeypress(caca: TXevent) =
  var 
    eve = caca.xkey
    keycode = cast[TKeyCode](eve.keycode)
    keysym = XKeycodeToKeysym(display, keycode, 0)
    casekeysym = XkeysymToString(keysym)
    newMod: cuint

  for i in keybinds:
    case i.mods
    of "meta": newMod = 64
    of "alt": newMod = 8
  
    if casekeysym == i.keys and eve.state == newMod:
      parseCmd(i.cmd)

# procedure which sets up mwm 
proc setup =
  display = XOpenDisplay(nil)
  
  if display == nil:
    quit "wtf, are you seriously trying to run a graphical app without a monitor!?"

  echo("Huh, at least you have a display")  
   
  
  root = DefaultRootWindow(display) # assign the root window to root var

  parseConfig()
  
  discard XGrabButton(display, 1, Mod4Mask, root, # grab mouse button 1
   true, ButtonPressMask, GrabModeAsync, GrabModeAsync, None, None)
  
  discard XGrabButton(display, 3, Mod4Mask, root, # grab mouse button 3
   true, ButtonPressMask, GrabModeAsync, GrabModeAsync, None, None)
    
  grabKeypress()
  
# main procedure which runs it all
proc main =

  setup()
  
  while true:
  
    discard XNextEvent(display,ev.addr)
  
    case ev.theType:

      of ConfigureRequest:
        echo "wanting to get configured"
    
      of MapRequest:
        var 
          wa: TXWindowAttributes
          ev: TXMapRequestEvent = ev.xmaprequest
        echo "lol, its working"
        echo ev.window
  
      of KeyPress: # Will only register the ones from XGrabKey

        handleKeypress(ev)

      of ButtonPress: # Will only register the ones from XGrabKey  
  
        var bev = cast[PXButtonEvent](ev.addr)[]
  
        if not bev.subwindow.addr.isNil:
  
          discard XGrabPointer(display, bev.subwindow, true,
  
                              PointerMotionMask or ButtonReleaseMask, GrabModeAsync,
  
                              GrabModeAsync, None, None, CurrentTime)
  
          discard XGetWindowAttributes(display, bev.subwindow, attr.addr);
  
          start = bev;
  
      of MotionNotify:
  
        var mnev = cast[PXMotionEvent](ev.addr)[]
  
        var bev = cast[PXButtonEvent](ev.addr)[]
  
        while XCheckTypedEvent(display,MotionNotify,ev.addr):
  
          continue
  
        var
  
          xdiff = bev.x_root - start.x_root
  
          ydiff = bev.y_root - start.y_root
  
        discard XMoveResizeWindow(display,mnev.window,
  
                                  attr.x + (if start.button==1: xdiff else: 0),
  
                                  attr.y + (if start.button==1: ydiff else: 0),
  
                                  max(1, attr.width + (if start.button==3: xdiff else: 0)),
  
                                  max(1, attr.height + (if start.button==3: ydiff else: 0)))
                                 
  
      of ButtonRelease:
  
        discard XUngrabPointer(display, CurrentTime)
  
      else: # Ignore unknown events
  
        continue

main()