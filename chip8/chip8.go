package chip8

/*
==============
opcode processor
==============
*/

import (
	fmt "fmt"
	log "log"
	os "os"
)

const (
	CHIP8_WIDTH = 64
	CHIP8_HEIGHT = 32
	DISPLAY_SCALE = 15
	GFX_SIZE = (CHIP8_WIDTH * CHIP8_HEIGHT) * DISPLAY_SCALE
)

var characters = []uint8 {
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
var gfx[GFX_SIZE]uint8 // display

var opcode uint16 // current opcode (instruction)
var stack [16]uint16 // 16 levels of the stack in total
var key [16]uint8 // current state of the key

var V[16]uint8 // general cpu registers V0 - VF (16) *VF used as a carry flag for some intructions
var I uint16 // store memory addresses
var pc uint16 // current executing address
var sp uint16 // points to the topmost level of the stack
var delay_timer uint8 // delay timer
var sound_timer uint8 // sound timer

// Initializes CPU
func Initialize() {
	const START_ADDRESS int = 0x050
	for i := 0; i < len(characters); i++ {
		memory[START_ADDRESS + i] = characters[i]
	}
}

// fetches data from c8 file and stores it in memory
func LoadGame(file_path string) {
	const START_ADDRESS int = 0x200

	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	file_len, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	ROM_data := readNextBytes(file, file_len.Size())
	for i := 0; i < len(ROM_data); i++ {
		memory[START_ADDRESS + i] = ROM_data[i]
	}

	fmt.Println(memory)
}

func EmulateCycle() {
	// TODO
}

func DrawFlag() bool {
	return true
}

func SetKeys() {
	// TODO
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