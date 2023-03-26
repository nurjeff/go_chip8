package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Supersampling scales up the screen res, since 64x32 is a little tiny
const SUPERSAMPLING = 16

// Boot up, read a rom and start emulating
func main() {
	c8 := CHIP8{}
	defer c8.Boot(SUPERSAMPLING).Destroy()
	defer sdl.Quit()
	c8.LoadROM("./roms/PONG2.ch8")
	for {
		c8.emulateCycle()
	}
}
