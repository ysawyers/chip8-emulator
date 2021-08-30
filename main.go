package main

/*
==============
MAIN GAME LOOP
==============
*/

import (
	fmt "fmt"
	runtime "runtime"
	log "log"

	chip8 "github.com/BigBellyBigDreams/chip8-emulator/chip8"
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
	program := initOpenGL()

	chip8.Initialize()
	chip8.LoadGame("/Users/yonden/Desktop/chip8-emulator/roms/PONG") // temporary for testing

	for !window.ShouldClose() {
		draw(window, program)

		chip8.EmulateCycle()

		if chip8.DrawFlag() {
			drawGraphics()
		}

		chip8.SetKeys()
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

// draws initial game state
func draw(window *glfw.Window, program uint32) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(program)
    
    glfw.PollEvents() // checks for keyboard or mouse inputs
    window.SwapBuffers() // GLFW lib uses double buffering
}

func drawGraphics() {
	// TODO
}