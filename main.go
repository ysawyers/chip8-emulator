package main

import (
	fmt "fmt"
	runtime "runtime"
	log "log"

	c8 "github.com/BigBellyBigDreams/chip8-emulator/chip8"
	gl "github.com/go-gl/gl/v2.1/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
)

const (
	CHIP8_WIDTH = 64
	CHIP8_HEIGHT = 32
	DISPLAY_SCALE = 15
)

func main() {
	runtime.LockOSThread() // glfw requires everything to run on a single thread

	window := initGlfw()
	defer glfw.Terminate()
	// program := initOpenGL()

	c8.Initialize() // initialize CPU
	c8.LoadGame("/Users/yonden/Desktop/chip8-emulator/roms/PONG") // temporary for testing

	for !window.ShouldClose() {
		glfw.PollEvents() // figure out what this does**

		c8.EmulateCycle() // process opcodes

		if c8.DrawFlag() {
			drawGraphics()
			window.SwapBuffers() // display visual changes from back buffer to front buffer
		}

		window.SetKeyCallback(processInput) // listen to keyboard events on the window and process them
	}
}

// initializes glfw window & properties
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(CHIP8_WIDTH * DISPLAY_SCALE, CHIP8_HEIGHT * DISPLAY_SCALE, "CHIP-8 EMULATOR", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	return window
}

// initializes openGL & creates program
func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
            panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    program := gl.CreateProgram()
    gl.LinkProgram(program)
    return program
}

func drawGraphics() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// add pixel information to video buffer
	// map that to back buffer
}

func processInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == 1 {
		switch key {
		case 259:
			fmt.Println("Quit Program")
		case 88:
			c8.Keys[0] = 1
		case 49:
			c8.Keys[1] = 1
		case 50:
			c8.Keys[2] = 1
		case 51:
			c8.Keys[3] = 1
		case 81:
			c8.Keys[4] = 1
		case 87:
			c8.Keys[5] = 1
		case 69:
			c8.Keys[6] = 1
		case 65:
			c8.Keys[7] = 1
		case 83:
			c8.Keys[8] = 1
		case 68:
			c8.Keys[9] = 1
		case 90:
			c8.Keys[0xA] = 1
		case 67:
			c8.Keys[0xB] = 1
		case 52:
			c8.Keys[0xC] = 1
		case 82:
			c8.Keys[0xD] = 1
		case 70:
			c8.Keys[0xE] = 1
		case 86:
			c8.Keys[0xF] = 1
		}
	} else if action == 0 {
		switch key {
		case 259:
			fmt.Println("Quit Program")
		case 88:
			c8.Keys[0] = 0
		case 49:
			c8.Keys[1] = 0
		case 50:
			c8.Keys[2] = 0
		case 51:
			c8.Keys[3] = 0
		case 81:
			c8.Keys[4] = 0
		case 87:
			c8.Keys[5] = 0
		case 69:
			c8.Keys[6] = 0
		case 65:
			c8.Keys[7] = 0
		case 83:
			c8.Keys[8] = 0
		case 68:
			c8.Keys[9] = 0
		case 90:
			c8.Keys[0xA] = 0
		case 67:
			c8.Keys[0xB] = 0
		case 52:
			c8.Keys[0xC] = 0
		case 82:
			c8.Keys[0xD] = 0
		case 70:
			c8.Keys[0xE] = 0
		case 86:
			c8.Keys[0xF] = 0
		}
	}
}