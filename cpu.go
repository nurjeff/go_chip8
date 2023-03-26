package main

import (
	"fmt"
)

type CPU struct {
	PC     uint16      // Program counter
	OC     Instruction // Current c.OC
	I      uint16      // Index register
	Mem    [4096]uint8 // Memory
	V      [16]uint8   // Registers
	DelayT uint8       // Delay timer
	SoundT uint8       // Sound timer
	Stack  [16]uint16  // Stack
	SP     uint16      // Stack pointer
	Key    [16]uint8   // Hex Keypad input
}

// Copy fontset in the memory at position 0 - n
func (c *CPU) loadFontset() {
	copy(c.Mem[:], fontSet)
}

// Fetch the next OpCode
func (c *CPU) fetchOpcode() {
	c.OC = Instruction(uint16(c.Mem[c.PC])<<8 | uint16(c.Mem[c.PC+1]))
}

// Timers should be decreased by 1 per clock
// When sound timer reaches 0 in the current clock, play beep.
func (c *CPU) updateTimers() {
	if c.DelayT > 0 {
		c.DelayT--
	}
	if c.SoundT > 0 {
		if c.SoundT == 1 {
			fmt.Print("BEEP\a\n")
		}
		c.SoundT--
	}
}

// Read x from the current OpCode - e.G. 0xFx00
func (c *CPU) readX() Instruction {
	return c.OC & 0x0F00 >> 8
}

// Read y from the current OpCode - e.G. 0xF0y0
func (c *CPU) readY() Instruction {
	return c.OC & 0x00F0 >> 4
}

// Read y from the current OpCode - e.G. 0xF0kk
func (c *CPU) readKK() uint8 {
	return uint8(c.OC & 0x00FF)
}

// Fontset holds the sprite information for the supported runes
var fontSet []uint8 = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}
