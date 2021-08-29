package chip8

/*
==============
opcode processor
==============
*/

import (
	f "fmt"
)

var opcode uint16
var memory [4096]byte
var V[16]byte
var I uint16
var pc uint16

func chip8() {
	f.Println("Hello World")
}