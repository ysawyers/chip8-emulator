package chip8

import (
	fmt "fmt"
	rand "math/rand"
	log "log"
	os "os"
)

const (
	MEMORY_START_ADDRESS = 0x200
	FONTSET_START_ADDRESS = 0x050
)

var fontset = [80]uint8 {
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

var memory[4096]uint8 // 4k memory
var Gfx[(64 * 32) * 15]uint8 // video buffer
var DrawFlag bool

var opcode uint16 // current opcode (instruction)
var stack [16]uint16 // 16 levels of the stack in total
var Keys [16]uint8 // current state of the key

var V[16]uint8 // general cpu registers V0 - VF (16) *VF used as a carry flag for some intructions
var I uint16 // store memory addresses
var pc uint16 // current executing address
var sp uint16 // points to the topmost level of the stack
var delayTimer uint8 // delay timer
var soundTimer uint8 // sound timer

// Initializes CPU
func Initialize() {
	pc = MEMORY_START_ADDRESS // Program Counter starts at 0x200
	opcode = 0 // Reset current opcode
	I = 0 // Reset index register
	sp = 0 // Reset stack pointer
	DrawFlag = false // Reset draw flag to false

	for i := 0; i < len(fontset); i++ {
		memory[FONTSET_START_ADDRESS + i] = fontset[i]
	}
}

// fetches data from c8 file and stores it in memory
func LoadGame(file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	fileStats, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	romData := readNextBytes(file, fileStats.Size())
	for i := 0; i < len(romData); i++ {
		memory[MEMORY_START_ADDRESS + i] = romData[i]
	}
}

func EmulateCycle() { // implement function pointers in the future instead of switch statement
	// Fetch opcode
	opcode = (uint16(memory[pc]) << 8 | uint16(memory[pc + 1]))

	// Decode opcode and execute
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000: // 00E0: Clears the screen
			for i := 0; i < len(Gfx); i++ {
				Gfx[i] = 0
			}
			DrawFlag = true
			pc += 2

		case 0x000E: // 00EE: Returns from a subroutine
			sp--
			pc = stack[sp]
			pc += 2

		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", opcode)
		}
	
	case 0x1000: // 1nnn: Jump to location nnn
		pc = opcode & 0x0FFF

	case 0x2000: // 2nnn: Call subroutine at nnn
		stack[sp] = pc
		sp++
		pc = opcode & 0x0FFF
	
	case 0x3000: // 3xkk: Skip next instruction if Vx = kk
		if uint16(V[(opcode & 0x0F00) >> 8]) == opcode & 0x00FF {
			pc += 4
		} else {
			pc += 2
		}
	
	case 0x4000: // 4xkk: Skip next instruction if Vx != kk
		if uint16(V[(opcode & 0x0F00) >> 8]) != opcode & 0x00FF {
			pc += 4
		} else {
			pc += 2
		}
	
	case 0x5000: // 5xy0: Skip next instruction if Vx = Vy
		if V[(opcode & 0x0F00) >> 8] == V[(opcode & 0x00F0) >> 4] {
			pc += 4
		} else {
			pc += 2
		}
	
	case 0x6000: // 6xkk: Set Vx = kk
		V[(opcode & 0x0F00) >> 8] = uint8(opcode & 0x00FF)
		pc += 2
	
	case 0x7000: // 7xkk: Set Vx = Vx + kk
		V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] + uint8(opcode & 0x00FF)
		pc += 2
	
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000: // 8xy0: Set Vx = Vy
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x00F0) >> 4]
			pc += 2
		
		case 0x0001: // 8xy1: Set Vx = Vx OR Vy
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] | V[(opcode & 0x00F0) >> 4]
			pc += 2
		
		case 0x0002: // 8xy2: Set Vx = Vx AND Vy
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] & V[(opcode & 0x00F0) >> 4]
			pc += 2
		
		case 0x0003: // 8xy3: Set Vx = Vx XOR Vy
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] ^ V[(opcode & 0x00F0) >> 4]
			pc += 2

		case 0x0004: // 8xy4: Set Vx = Vx + Vy, set VF = carry
			if V[(opcode & 0x00F0) >> 4] > 0xFF - V[(opcode & 0x0F00) >> 8] {
				V[0xF] = 1
			} else {
				V[0xF] = 0
			}
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] + V[(opcode & 0x00F0 >> 4)]
			pc += 2
		
		case 0x0005: // 8xy5: Set Vx = Vx - Vy, set VF = NOT borrow
			if V[(opcode & 0x00F0) >> 4] > V[(opcode & 0x0F00) >> 8] {
				V[0xF] = 0
			} else {
				V[0xF] = 1
			}
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] - V[(opcode & 0x00F0) >> 4]
			pc += 2
		
		case 0x0006: // 8xy6: Set Vx = Vx SHR 1
			V[0xF] = V[(opcode & 0x0F00) >> 8] & 0x1
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] >> 1
			pc += 2

		case 0x0007: // 8xy7: Set Vx = Vy - Vx, set VF = NOT borrow
			if V[(opcode & 0x0F00) >> 8] > V[(opcode & 0x00F0) >> 4] {
				V[0xF] = 0
			} else {
				V[0xF] = 1
			}
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x00F0) >> 4] - V[(opcode & 0x0F00) >> 8]
			pc += 2
		
		case 0x000E: // 8xyE: Set Vx = Vx SHL 1
			V[0xF] = V[(opcode & 0x0F00) >> 8] >> 7
			V[(opcode & 0x0F00) >> 8] = V[(opcode & 0x0F00) >> 8] << 1
			pc += 2
		
		default:
			fmt.Printf("Unknown opcode [0x8000]: 0x%X\n", opcode)
		}
	
	case 0x9000: // 9xy0: Skip next instruction if Vx != Vy
		if V[(opcode & 0x0F00) >> 8] != V[(opcode & 0x00F0) >> 4] {
			pc += 4
		} else {
			pc += 2
		}
	
	case 0xA000: // Annn: Set I = nnn
		I = opcode & 0x0FFF
		pc += 2
	
	case 0xB000: // Bnnn: Jump to location nnn + V0
		pc = uint16(V[0x0]) + (opcode & 0x0FFF)
	
	case 0xC000: // Cxkk: Set Vx = random byte AND kk
		V[(opcode & 0x0F00) >> 8] = uint8(rand.Intn(256)) & uint8(opcode & 0x00FF) // error prone #1
		pc += 2
	
	case 0xD000: // Dxyn: Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision
		// x := V[(opcode & 0x0F00) >> 8]
		// y := V[(opcode & 0x00F0) >> 4]
		// h := opcode & 0x000F
		V[0xF] = 0

		// TODO

		DrawFlag = true
		pc += 2
	
	case 0xE000: 
		switch opcode & 0x00FF {
		case 0x009E: // Ex9E: Skip next instruction if key with the value of Vx is pressed
			if Keys[V[(opcode & 0x0F00) >> 8]] == 1 {
				pc += 4
			} else {
				pc += 2
			}
		
		case 0x00A1: // ExA1: Skip next instruction if key with the value of Vx is not pressed
			if Keys[V[(opcode & 0x0F00) >> 8]] == 0 {
				pc += 4
			} else {
				pc += 2
			}
		
		default:
			fmt.Printf("Invalid opcode [0xE000]: 0x%X", opcode)
		}
	
	case 0xF000: // Fx0A: Wait for a key press, store the value of the key in Vx 
		switch opcode & 0x00FF {
		case 0x0007: // Fx07: Set Vx = delay timer value
			V[(opcode & 0x0F00) >> 8] = delayTimer
			pc += 2
		
		case 0x000A: // Fx0A: Wait for a key press, store the value of the key in Vx
			var KeyPress bool = false
			for i := 0; i < len(Keys); i++ {
				if Keys[i] != 0 {
					V[(opcode & 0x0F00) >> 8] = uint8(i)
					KeyPress = true
				}
			}

			if !KeyPress {
				return
			}
		
		case 0x0015: // Fx15: Set delay timer = Vx
			delayTimer = V[(opcode & 0x0F00) >> 8]
			pc += 2
		
		case 0x0018: // Fx18: Set sound timer = Vx 
			soundTimer = V[(opcode & 0x0F00) >> 8]
			pc += 2
		
		case 0x001E: // Fx1E: Set I = I + Vx
			if I + uint16(V[(opcode & 0x0F00) >> 8]) > 0xFFF {
				V[0xF] = 1
			} else {
				V[0xF] = 0
			}
			I += uint16(V[(opcode & 0x0F00) >> 8])
			pc += 2
		
		case 0x0029: // Fx29: Set I = location of sprite for digit Vx
			I = uint16(V[(opcode & 0x0F00) >> 8]) * 0x5
			pc += 2
		
		case 0x0033: // Fx33: Store BCD representation of Vx in memory locations I, I+1, and I+2
			memory[I] = V[(opcode & 0x0F00) >> 8] / 100
			memory[I + 1] = (V[(opcode & 0x0F00) >> 8] / 10) % 10
			memory[I + 2] = (V[(opcode & 0x0F00) >> 8] % 100) / 10
			pc += 2
		
		case 0x0055: // Fx55: Store registers V0 through Vx in memory starting at location I
			for i := 0; i < int((opcode & 0x0F00) >> 8) + 1; i++ {
				memory[uint16(i) + I] = V[i]
			}
			I = ((opcode & 0x0F00) >> 8) + 1
			pc += 2
		
		case 0x0065: // Fx65: Read registers V0 through Vx from memory starting at location I
			for i := 0; i < int((opcode & 0x0F00) >> 8) + 1; i++ {
				V[i] = memory[I + uint16(i)]
			}
			I = ((opcode & 0x0F00) >> 8) + 1
			pc += 2
		
		default:
			fmt.Printf("Unknown opcode [0xF000]: 0x%X\n", opcode)
		}

	default: 
		fmt.Printf("Unknown opcode: 0x%X\n", opcode)
	}

	if delayTimer > 0 {
		delayTimer-- 
	}

	if soundTimer > 0 {
		if soundTimer == 1 {
			fmt.Println("BEEP!")
		}
		soundTimer--
	}
}

// parses n number of bytes in the file
func readNextBytes(file *os.File, n int64) []uint8 {
	bytes := make([]uint8, n)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func CurrentOpcodeDebug() {
	fmt.Printf("Current opcode: 0x%X \n", opcode)
}