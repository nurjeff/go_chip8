package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type CHIP8 struct {
	CPU
	GPU
}

// Setup CPU and GPU
func (c *CHIP8) Boot(superSampling uint32) *sdl.Window {
	c.CPU = CPU{PC: 0x200}
	c.GPU = GPU{ssFac: int32(superSampling)}

	c.CPU.loadFontset()
	return c.GPU.Initialize()
}

// Read a binary rom from the specified path and copy it to the memory, starting at index 512
func (c *CHIP8) LoadROM(path string) {
	fmt.Println("Loading rom:", path)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	stats, err := f.Stat()
	if err != nil {
		panic(err)
	}
	var size int64 = stats.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(f)
	_, err = buffer.Read(bytes)
	if err != nil {
		panic(err)
	}
	for index, byte := range bytes {
		c.CPU.Mem[index+512] = byte
	}
}

// Emulate a clock cycle
func (c *CHIP8) emulateCycle() {
	for i := 0; i < len(c.CPU.Key); i++ {
		c.CPU.Key[i] = c.GPU.sdlBuffer[i]
	}

	c.CPU.fetchOpcode()
	c.executeOpcode()
	c.CPU.updateTimers()
	c.GPU.drawGfx()

	// CPU Clock, this is not really specified anywhere as far as i can tell
	time.Sleep(time.Microsecond * (16667 / 1000))
}

// Fetch current CPU OpCode and call correct function pointer
func (c *CHIP8) executeOpcode() {
	switch c.CPU.OC & 0xF000 {

	case 0x000:
		instructionMap0x0[c.CPU.OC&0x000F](c)

	case (0x8000):
		instructionMap0x8[c.CPU.OC&0x800F](c)

	case (0xF000):
		instructionMap0xF[c.OC&0xF0FF](c)

	case (0xE000):
		instructionMap0xE[c.OC&0xE0FF](c)

	default:
		instructionMap[c.CPU.OC&0xF000](c)
	}
}
