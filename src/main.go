package main

import (
	// can you import the go glfw package for me?

	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	glfw.Init()

	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)
	window, err := glfw.CreateWindow(800, 600, "Hello World", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		log.Fatalln(err)
	}
	vertices := []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}
	indicies := []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	var (
		VAO           uint32
		VBO           uint32
		EBO           uint32
		VShader       uint32
		FShader       uint32
		shaderProgram uint32
	)

	shaderProgram = gl.CreateProgram()

	VShader, err = compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("VSHADER BAD")
		log.Fatal(err)
		fmt.Println("VSHADER BAD")
	}

	FShader, err = compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("FSHADEF BAD")
		log.Fatal(err)
		fmt.Println("FSHADEF BAD")
	}

	gl.AttachShader(shaderProgram, VShader)
	gl.AttachShader(shaderProgram, FShader)
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)

	if success == gl.FALSE {
		log.Println("bad")
		var logLength int32

		llog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(llog))
		log.Fatal(llog)
	}

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &EBO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies)*4, gl.Ptr(indicies), gl.STATIC_DRAW)
	// copy vertices data into VBO (it needs to be bound first)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	glfw.WindowHint(glfw.OpenGLProfile, gl.TRIANGLES)
	for !window.ShouldClose() {
		// Do OpenGL stuff.

		gl.ClearColor(0.5, 0.5, 0.5, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		time := glfw.GetTime()
		green := (float32(math.Sin(time)) / float32(2)) + float32(0.5)
		fmt.Println(green)
		vertexColorLocation := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor\000"))
		gl.UseProgram(shaderProgram)
		gl.Uniform4f(vertexColorLocation, 0.0, green, 0, 1)
		// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.Ptr(nil))

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteBuffers(1, &EBO)
	gl.DeleteShader(FShader)
	gl.DeleteShader(VShader)
}

const vertexShader = `#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
	gl_Position = vec4(aPos, 1.0);
}
`

const fragmentShader = `#version 330 core
out vec4 FragColor;

uniform vec4 ourColor;

void main()
{
	FragColor = ourColor;
}
`

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		fmt.Println("shader type:")
		fmt.Println(shader)
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
