package main

import (
	"math/rand"
)

type Instruction uint16

// Instruction maps hold all the instruction set function pointers to avoid a switch statement
// Multiple instruction maps are needed for cases where the MSB is not enough to determine which function should be called
var instructionMap map[Instruction]func(*CHIP8) = map[Instruction]func(*CHIP8){
	0x1000: JP,
	0x2000: CALL,
	0x3000: SEVX,
	0x4000: SNEVXB,
	0x5000: SEVXVY,
	0x6000: LDVX,
	0x7000: ADDVX,
	0x9000: SNEVXVY,
	0xA000: LDI,
	0xB000: JPV0,
	0xC000: RND,
	0xD000: DRW,
}

var instructionMap0x0 map[Instruction]func(*CHIP8) = map[Instruction]func(*CHIP8){
	0x000: CLS,
	0x00E: RET,
}

var instructionMap0x8 map[Instruction]func(*CHIP8) = map[Instruction]func(*CHIP8){
	0x8000: LDVXVY,
	0x8001: ORVXVY,
	0x8002: ANDVXVY,
	0x8003: XORVXVY,
	0x8004: ADDVXVY,
	0x8005: SUBVXVY,
	0x8006: SHRVX,
	0x8007: SUBNVX,
	0x800E: SHLVX,
}

var instructionMap0xF map[Instruction]func(*CHIP8) = map[Instruction]func(*CHIP8){
	0xF007: LDVXDT,
	0xF00A: LDVXK,
	0xF015: LDDT,
	0xF018: LDST,
	0xF01E: ADDI,
	0xF029: LDF,
	0xF033: LDB,
	0xF055: LDIVX,
	0xF065: LDVXI,
}

var instructionMap0xE map[Instruction]func(*CHIP8) = map[Instruction]func(*CHIP8){
	0xE09E: SKP,
	0xE0A1: SKNP,
}

// -- INSTRUCTIONS --

// -- 0xE --
// Skip next instruction if key with the value of Vx is pressed
func SKP(c *CHIP8) {
	if c.Key[c.CPU.V[c.CPU.readX()]] != 0 {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}

// Skip next instruction if key with the value of Vx is not pressed
func SKNP(c *CHIP8) {
	if c.Key[c.CPU.V[c.CPU.readX()]] == 0 {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}

// -- 0xF --
// Wait for a key press, store the value of the key in Vx
func LDVXK(c *CHIP8) {
	for index, element := range c.Key {
		if element > 0 {
			c.CPU.V[c.CPU.readX()] = uint8(index)
			c.CPU.PC += 2
		}
	}
}

// Set sound timer = Vx
func LDST(c *CHIP8) {
	c.CPU.SoundT = c.CPU.V[c.CPU.readX()]
	c.CPU.PC += 2
}

// Set I = location of sprite for digit Vx
func LDF(c *CHIP8) {
	c.CPU.I = uint16(c.CPU.V[c.CPU.readX()] * 0x5)
	c.CPU.PC += 2
}

// Store registers V0 through Vx in memory starting at location I
func LDIVX(c *CHIP8) {
	for i := 0; i <= int(c.CPU.readX()); i++ {
		c.CPU.Mem[i+int(c.CPU.I)] = c.CPU.V[i]
	}
	c.CPU.PC += 2
}

// Read registers V0 through Vx from memory starting at location I
func LDVXI(c *CHIP8) {
	for i := 0; i <= int(c.CPU.readX()); i++ {
		c.CPU.V[i] = c.CPU.Mem[i+int(c.CPU.I)]
	}
	c.CPU.PC += 2
}

// Store BCD representation of Vx in memory locations I, I+1, and I+2
func LDB(c *CHIP8) {
	c.CPU.Mem[c.CPU.I] = c.CPU.V[(c.OC&0x0F00)>>8] / 100
	c.CPU.Mem[c.CPU.I+1] = (c.CPU.V[(c.OC&0x0F00)>>8] / 10) % 10
	c.CPU.Mem[c.CPU.I+2] = (c.CPU.V[(c.OC&0x0F00)>>8] % 100) % 10
	c.CPU.PC += 2
}

// Set Vx = delay timer value
func LDVXDT(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] = c.CPU.DelayT
	c.CPU.PC += 2
}

// Set delay timer = Vx
func LDDT(c *CHIP8) {
	c.CPU.DelayT = c.CPU.V[c.CPU.readX()]
	c.CPU.PC += 2
}

// Set I = I + Vx
func ADDI(c *CHIP8) {
	c.CPU.I = c.CPU.I + uint16(c.CPU.V[c.CPU.readX()])
	c.CPU.PC += 2
}

// -- 0x8 --
// Set Vx = Vx SHL 1
func SHLVX(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()]&0xF == 1 {
		c.CPU.V[0xF] = 1
	} else {
		c.CPU.V[0xF] = 0
	}
	c.CPU.V[c.CPU.readX()] *= 2
	c.CPU.PC += 2
}

// Set Vx = Vy - Vx, set VF = NOT borrow
func SUBNVX(c *CHIP8) {
	if c.CPU.V[c.CPU.readY()] > c.CPU.V[c.CPU.readX()] {
		c.CPU.V[0xF] = 1
	} else {
		c.CPU.V[0xF] = 0
	}
	c.CPU.V[c.CPU.readX()] = c.CPU.V[c.CPU.readY()] - c.CPU.V[c.CPU.readX()]
	c.CPU.PC += 2
}

// Set Vx = Vx SHR 1
func SHRVX(c *CHIP8) {
	if uint8(c.CPU.V[c.CPU.readX()]%2) == 1 {
		c.CPU.V[0xF] = 1
	} else {
		c.CPU.V[0xF] = 0
	}
	c.CPU.V[c.CPU.readX()] /= 2
	c.CPU.PC += 2
}

// Set Vx = Vx - Vy, set VF = NOT borrow
func SUBVXVY(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()] > c.CPU.V[c.CPU.readY()] {
		c.CPU.V[0xF] = 1
	} else {
		c.CPU.V[0xF] = 0
	}
	c.CPU.V[c.CPU.readX()] = c.CPU.V[c.CPU.readX()] - c.CPU.V[c.CPU.readY()]
	c.CPU.PC += 2
}

// Set Vx = Vx + Vy, set VF = carry
func ADDVXVY(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()]+c.CPU.V[c.CPU.readY()] > 255 {
		c.CPU.V[0xF] = 1
	} else {
		c.CPU.V[0xF] = 0
	}
	c.CPU.V[(c.OC&0x0F00)>>8] += c.CPU.V[(c.OC&0x00F0)>>4]
	c.CPU.PC += 2
}

// Set Vx = Vx XOR Vy
func XORVXVY(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] ^= c.CPU.V[c.CPU.readY()]
	c.CPU.PC += 2
}

// Set Vx = Vx AND Vy
func ANDVXVY(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] = c.CPU.V[c.CPU.readX()] & c.CPU.V[c.CPU.readY()]
	c.CPU.PC += 2
}

// Set Vx = Vx OR Vy
func ORVXVY(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] = c.CPU.V[c.CPU.readX()] | c.CPU.V[c.CPU.readY()]
	c.CPU.PC += 2
}

// Set Vx = Vy
func LDVXVY(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] = c.CPU.V[c.CPU.readY()]
	c.CPU.PC += 2
}

// -- 0x0 --
// Return from a subroutine
func RET(c *CHIP8) {
	c.PC = c.Stack[c.SP]
	c.SP--
	c.CPU.PC += 2
}

// Clear the display
func CLS(c *CHIP8) {
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			c.GPU.VRAM[x][y] = 0
		}
	}
	c.CPU.PC += 2
}

// --
// Skip next instruction if Vx = Vy
func SEVXVY(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()] == c.CPU.V[c.CPU.readY()] {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}

// Skip next instruction if Vx != kk
func SNEVXB(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()] != c.CPU.readKK() {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}

// Set Vx = random byte AND kk
func RND(c *CHIP8) {
	rand := uint8(rand.Intn(255))
	c.CPU.V[c.CPU.readX()] = rand & c.CPU.readKK()
	c.CPU.PC += 2
}

// Jump to location nnn
func JP(c *CHIP8) {
	c.CPU.PC = uint16(c.OC & 0x0FFF)
}

// Set I = nnn
func LDI(c *CHIP8) {
	c.CPU.I = uint16(c.OC & 0x0FFF)
	c.CPU.PC += 2
}

// Jump to location nnn + V0
func JPV0(c *CHIP8) {
	c.CPU.PC = uint16(c.OC&0x0FFF) + uint16(c.CPU.V[0])
}

// Call subroutine at nnn
func CALL(c *CHIP8) {
	c.SP++
	c.Stack[c.SP] = c.PC
	c.PC = uint16(c.OC & 0x0FFF)
}

// Set Vx = kk
func LDVX(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] = c.CPU.readKK()
	c.CPU.PC += 2
}

// Set Vx = Vx + kk
func ADDVX(c *CHIP8) {
	c.CPU.V[c.CPU.readX()] += c.CPU.readKK()
	c.CPU.PC += 2
}

// Skip next instruction if Vx != Vy
func SNEVXVY(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()] != c.CPU.V[c.CPU.readY()] {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}

// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision
func DRW(c *CHIP8) {
	var pixel uint8 = 0
	c.CPU.V[0xF] = 0

	for yLine := 0; yLine < int(uint8(c.OC&0x000F)); yLine++ {
		pixel = c.CPU.Mem[c.CPU.I+uint16(yLine)]
		for xLine := 0; xLine < 8; xLine++ {
			if (pixel & (0x80 >> xLine)) != 0 {
				xpos := (c.CPU.V[(c.OC&0x0F00)>>8] + uint8(xLine))
				ypos := uint8(yLine + int(c.CPU.V[(c.OC&0x00F0)>>4]))
				// Test wrapping around the screen, afaik this is not accurate to the original emulation
				if (xpos >= SCREEN_WIDTH) || (ypos >= SCREEN_HEIGHT) {
					continue
				}
				if c.GPU.VRAM[xpos][ypos] == 1 {
					c.CPU.V[0xF] = 1
				}

				c.GPU.VRAM[xpos][ypos] ^= 1

			}
		}
	}
	c.GPU.df = true
	c.CPU.PC += 2
}

// Skip next instruction if Vx = kk
func SEVX(c *CHIP8) {
	if c.CPU.V[c.CPU.readX()] == c.CPU.readKK() {
		c.CPU.PC += 2
	}
	c.CPU.PC += 2
}
