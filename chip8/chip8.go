package chip8

import (
	rand "math/rand"
	log "log"
	os "os"
)

const (
	FONTSET_SIZE = 80
	MEMORY_START_ADDRESS = 0x200
	FONTSET_START_ADDRESS = 0x050
)

var fontset = [FONTSET_SIZE]uint8 {
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
var Gfx[64 * 32]uint8 // video buffer

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
	for i := 0; i < FONTSET_SIZE; i++ {
		memory[FONTSET_START_ADDRESS + i] = fontset[i]
	}

	pc = MEMORY_START_ADDRESS
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

func EmulateCycle() {
	// Process opcodes
}

func DrawFlag() bool {
	return false
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

func randGen() uint8 {
	randByte := uint8(rand.Intn(255))
	return randByte
}