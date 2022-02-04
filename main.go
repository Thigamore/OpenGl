package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//Define dim for window
const (
	WIDTH  = 500
	HEIGHT = 500
)

func main() {
	//Vertices of a triangle
	triangle := []float32{
		0, 0.5, 0, //Top
		-0.5, -0.5, 0, //Left
		0.5, -0.5, 0, //Right
	}

	//Locks the goroutine to the OS thread, aka software thread
	//Needed when working with OS stuff or C
	runtime.LockOSThread()

	//Preps the window
	window := initGLFW()
	//stops all windows
	defer glfw.Terminate()

	//Preps a new program
	program := initOpenGl()

	//Signal for whether window is closed
	for !window.ShouldClose() {

		//Draws something on the window using shaders at program
		draw(window, program)
	}
}

//Initializes a GLFW window
func initGLFW() *glfw.Window {
	//Checks to see if glfw can run
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	////Set hints for when the window is created
	//Sets the window as resizable
	glfw.WindowHint(glfw.Resizable, glfw.False)
	//Sets the version of opengl
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//Sets compatibility for future
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	//Create the window
	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "Test", nil, nil)
	if err != nil {
		panic(err)
	}
	//Binds window to thread
	window.MakeContextCurrent()

	return window
}

//Initializes OpenGl and returns a program (reference to store shaders)
func initOpenGl() uint32 {
	//Checks to see if openGl can run
	if err := gl.Init(); err != nil {
		panic(err)
	}

	//Print the openGl version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGl Version: ", version)

	//Create a program (Reference to store shaders)
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

//Draws something on the window using shaders at program
func draw(window *glfw.Window, program uint32) {
	//Clears the buffer to preset value
	//Buffer: An array of unformatted memory allocated by the OpenGl Context(GPU)
	//Aka remove anything on screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//Tells OpenGl to use the program reference
	gl.UseProgram(program)

	//Glfw returns any keyboard or mouse events
	glfw.PollEvents()
	//Important cuz buffer swapping
	//Stuff is initially drawn on invisible canvas, then moved to visible canvas
	window.SwapBuffers()
}

//Makes a Vetex Object Array (VAO) which is a set of points used to draw
//Returns pointer to OpenGl vao
func makeVao(points []float32) uint32 {
	//Create vbo to bind vao to
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	//size is 4x cuz 4 bytes in an f32
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
