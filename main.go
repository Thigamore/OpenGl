package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Thigamore/OpenGl/src/IndexBuffer"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//Define dim for window
const (
	WIDTH  = 500
	HEIGHT = 500
)

type ShaderSource struct {
	VertexSource   *string
	FragmentSource *string
}

func main() {

	//Vertices of a positions
	positions := []float32{
		-0.5, 0.5, //Top Left
		0.5, 0.5, //Top Right
		0.5, -0.5, // Bottom Right
		-0.5, -0.5, //Bottom Left
	}

	indices := []uint32{
		0, 1, 2, //First Triangle
		2, 3, 0, //Second Triangle
	}

	//Get the shaders
	shaderSrc := ParseShader()

	//Locks the goroutine to the OS thread, aka software thread
	//Needed when working with OS stuff or C
	runtime.LockOSThread()

	//Preps the window
	window := InitGLFW()
	//stops all windows
	defer glfw.Terminate()

	glfw.SwapInterval(1)

	InitOpenGl()

	ibo := IndexBuffer.InitIndexBuffer()

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	shader := CreateShader(shaderSrc.VertexSource, shaderSrc.FragmentSource)
	gl.UseProgram(shader)

	location := gl.GetUniformLocation(shader, gl.Str("u_Color\x00"))
	gl.ProgramUniform4f(shader, location, 0.2, 0.3, 0.8, 1.0)

	r := float32(0.0)
	increment := float32(0.05)

	//Signal for whether window is closed
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		ProcessInput(window)

		gl.ProgramUniform4f(shader, location, r, 0.3, 0.8, 1.0)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		if r > 1 {
			increment = -0.05
		} else if r < 0 {
			increment = 0.05
		}

		r += increment

		window.SwapBuffers()
		glfw.PollEvents()

	}
}

//Initializes a GLFW window
func InitGLFW() *glfw.Window {
	//Checks to see if glfw can run
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// -- Set hints for when the window is created --
	//Sets the window as resizable
	glfw.WindowHint(glfw.Resizable, glfw.False)
	//Sets the version of opengl
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
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

// Initializes OpenGl
func InitOpenGl() {
	//Initializes openGl
	if err := gl.Init(); err != nil {
		panic(err)
	}
	//Print the openGl version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGl Version: ", version)
}

//Deals with input
func ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == 1 {
		window.SetShouldClose(true)
	}
}

func CompileShader(sType uint32, source *string) uint32 {
	id := gl.CreateShader(sType)
	src, free := gl.Strs(*source)

	gl.ShaderSource(id, 1, src, nil)
	free()
	gl.CompileShader(id)

	//Error Handling
	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)
		message := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(id, length, &length, gl.Str(message))
		var sTypeString string
		if sType == gl.VERTEX_SHADER {
			sTypeString = "Vertex"
		} else if sType == gl.FRAGMENT_SHADER {
			sTypeString = "Fragment"
		}
		fmt.Println("Failed to compile " + sTypeString + " shader!")
		fmt.Println(message)
		gl.DeleteProgram(id)
		return 0
	}

	return id
}

//Compiles the shaders
func CreateShader(vertexShader *string, fragmentShader *string) uint32 {
	//Creates a program
	program := gl.CreateProgram()
	vs := CompileShader(gl.VERTEX_SHADER, vertexShader)
	fs := CompileShader(gl.FRAGMENT_SHADER, fragmentShader)

	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program
}

//Parses the shader file
func ParseShader() *ShaderSource {
	//Open file
	file, err := os.Open("res\\shaders\\Basic.shader")
	if err != nil {
		panic("Error parsing shaders: " + err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var line string
	var shaderType int
	shaders := make([]string, 2)
	var (
		VERTEX   = 0
		FRAGMENT = 1
	)

	//Scans through lines to put shaders in a specific order
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "#shader") {
			if strings.Contains(line, "vertex") {
				shaderType = VERTEX
			} else if strings.Contains(line, "fragment") {
				shaderType = FRAGMENT
			}
		} else {
			shaders[shaderType] += line + "\n"
		}
	}
	//Add end line
	for pos := range shaders {
		shaders[pos] += "\x00"
	}

	return &ShaderSource{
		VertexSource:   &shaders[VERTEX],
		FragmentSource: &shaders[FRAGMENT],
	}
}
