package main

/*
==============
MAIN GAME LOOP
==============
*/

import (
	"fmt"
	"runtime"
	"github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread() // GLFW requires everything to run on a single thread

	window := initGlfw()
	defer glfw.Terminate()
	
	program := initOpenGL()
	setupInput()

	for !window.ShouldClose() {
		draw(window, program)
	}
}

// initializes glfw window & properties
func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	width, height := 500, 500
	window, err := glfw.CreateWindow(width, height, "CHIP-8 EMULATOR", nil, nil)
	if err != nil {
		panic(err)
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

    prog := gl.CreateProgram()
    gl.LinkProgram(prog)
    return prog
}

// draws initial game state
func draw(window *glfw.Window, program uint32) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(program)
    
    glfw.PollEvents() // checks for keyboard or mouse inputs
    window.SwapBuffers() // GLFW lib uses double buffering
}

func setupInput() {
	// TODO
}