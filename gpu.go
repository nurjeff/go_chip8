package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 64
	SCREEN_HEIGHT = 32
)

type GPU struct {
	VRAM      [64][32]uint8 // GFX memory
	df        bool          // Drawflag
	ssFac     int32         // Supersampling factor
	surface   *sdl.Surface  // Drawing surface
	running   bool          // Used for SDL polling
	sdlBuffer [16]uint8     // ^
	window    *sdl.Window   // SDL window
}

// Keymap holds the key mapping from Chip8 hex keypad -> SDL input
var keyMap map[sdl.Keycode]int = map[sdl.Keycode]int{
	sdl.K_1: 0x1,
	sdl.K_2: 0x2,
	sdl.K_3: 0x3,
	sdl.K_4: 0xC,
	sdl.K_q: 0x4,
	sdl.K_w: 0x5,
	sdl.K_e: 0x6,
	sdl.K_r: 0xD,
	sdl.K_a: 0x7,
	sdl.K_s: 0x8,
	sdl.K_d: 0x9,
	sdl.K_f: 0xE,
	sdl.K_z: 0xA,
	sdl.K_x: 0x0,
	sdl.K_c: 0xB,
	sdl.K_v: 0xF,
}

// Setup SDL window and start polling input
func (g *GPU) Initialize() *sdl.Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	w, err := sdl.CreateWindow("CHIP8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH*g.ssFac, SCREEN_HEIGHT*g.ssFac, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	s, err := w.GetSurface()
	if err != nil {
		panic(err)
	}
	g.surface = s

	go g.fillSDLBuffer()

	g.window = w
	return w
}

// Hold input in a buffer, since sdl polling is much faster than CHIP8 clock
func (g *GPU) fillSDLBuffer() {
	g.running = true
	for g.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				g.running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.RELEASED {
					g.sdlBuffer[keyMap[t.Keysym.Sym]] = 0
				}
				if t.State != sdl.RELEASED {
					g.sdlBuffer[keyMap[t.Keysym.Sym]] = 1
				}
			}
		}
	}
}

// Draw the current VRAM to the SDL window
func (g *GPU) drawGfx() {
	g.surface.FillRect(nil, 10)
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			if g.VRAM[x][y] > 0 {
				rect := sdl.Rect{X: int32(x * int(g.ssFac)), Y: int32(y * int(g.ssFac)), W: 1 * g.ssFac, H: 1 * g.ssFac}
				g.surface.FillRect(&rect, 0xffff0000)

			} else {
				rect := sdl.Rect{X: int32(x * int(g.ssFac)), Y: int32(y * int(g.ssFac)), W: 1 * g.ssFac, H: 1 * g.ssFac}
				g.surface.FillRect(&rect, 0x00000000)
			}
		}
	}

	g.window.UpdateSurface()
}
