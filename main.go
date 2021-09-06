package main

import (
	fmt "fmt"
	// log "log"
	time "time"
	strings "strings"

	c8 "github.com/BigBellyBigDreams/chip8-emulator/chip8"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH = 64
	HEIGHT = 32
	VIDEO_SCALE = 15
)

var filePath string
var windowTitle string

func main() {
	fmt.Print("Enter game path: ")
	fmt.Scanf("%s", &filePath) 

	c8.Initialize() // Initialize CPU
	c8.LoadGame(filePath) // Load Game Into Memory

	tempArr := strings.Split(filePath, "/")
	windowTitle = tempArr[(len(tempArr) - 1)] 

	rl.InitWindow(WIDTH * VIDEO_SCALE, HEIGHT * VIDEO_SCALE, windowTitle)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		c8.EmulateCycle() // Decoding & Executing opcodes
		// c8.CurrentOpcodeDebug() 

		if c8.DrawFlag { // Draw If Flag Is Detected
			rl.ClearBackground(rl.Black)
			drawGraphics()
			rl.EndDrawing()
		}
		processInput()

		time.Sleep((1000 / 120) * time.Millisecond) // 60 Hz (60 iterations per second)
	}
	rl.CloseWindow()
}

func drawGraphics() {
	for j := 0; j < len(c8.Gfx); j ++ {
		for i := 0; i < len(c8.Gfx[j]); i++ {
			x := i * VIDEO_SCALE
			y := j * VIDEO_SCALE
			if c8.Gfx[j][i] != 0 {	
				rl.DrawRectangle(int32(x), int32(y), VIDEO_SCALE, VIDEO_SCALE, rl.White)
			} else {
				rl.DrawRectangle(int32(x), int32(y), VIDEO_SCALE, VIDEO_SCALE, rl.Black)
			}
		}
	}
	c8.DrawFlag = false
}

func processInput() {
	switch true {
		case rl.IsKeyDown(259): // Backspace 
			rl.CloseWindow()
		case rl.IsKeyDown(88): // X
			c8.Keys[0x0] = 1
		case rl.IsKeyDown(49): // 1
			c8.Keys[0x1] = 1
		case rl.IsKeyDown(50): // 2
			c8.Keys[0x2] = 1
		case rl.IsKeyDown(51): // 3
			c8.Keys[0x3] = 1
		case rl.IsKeyDown(81): // Q
			c8.Keys[0x4] = 1
		case rl.IsKeyDown(87): // W
			c8.Keys[0x5] = 1
		case rl.IsKeyDown(69): // E
			c8.Keys[0x6] = 1
		case rl.IsKeyDown(65): // A
			c8.Keys[0x7] = 1
		case rl.IsKeyDown(83): // S
			c8.Keys[0x8] = 1
		case rl.IsKeyDown(68): // D
			c8.Keys[0x9] = 1
		case rl.IsKeyDown(90): // Z
			c8.Keys[0xA] = 1
		case rl.IsKeyDown(67): // C
			c8.Keys[0xB] = 1
		case rl.IsKeyDown(52): // 4
			c8.Keys[0xC] = 1
		case rl.IsKeyDown(82): // R
			c8.Keys[0xD] = 1
		case rl.IsKeyDown(70): // F
			c8.Keys[0xE] = 1
		case rl.IsKeyDown(86): // V
			c8.Keys[0xF] = 1
		}
	
		switch true {
		case rl.IsKeyReleased(88):
			c8.Keys[0x0] = 0
		case rl.IsKeyReleased(49):
			c8.Keys[0x1] = 0
		case rl.IsKeyReleased(50):
			c8.Keys[0x2] = 0
		case rl.IsKeyReleased(51):
			c8.Keys[0x3] = 0
		case rl.IsKeyReleased(81):
			c8.Keys[0x4] = 0
		case rl.IsKeyReleased(87):
			c8.Keys[0x5] = 0
		case rl.IsKeyReleased(69):
			c8.Keys[0x6] = 0
		case rl.IsKeyReleased(65):
			c8.Keys[0x7] = 0
		case rl.IsKeyReleased(83):
			c8.Keys[0x8] = 0
		case rl.IsKeyReleased(68):
			c8.Keys[0x9] = 0
		case rl.IsKeyReleased(90):
			c8.Keys[0xA] = 0
		case rl.IsKeyReleased(67):
			c8.Keys[0xB] = 0
		case rl.IsKeyReleased(52):
			c8.Keys[0xC] = 0
		case rl.IsKeyReleased(82):
			c8.Keys[0xD] = 0
		case rl.IsKeyReleased(70):
			c8.Keys[0xE] = 0
		case rl.IsKeyReleased(86):
			c8.Keys[0xF] = 0
		}
}