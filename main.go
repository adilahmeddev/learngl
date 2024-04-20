package main

import (
	// can you import the go glfw package for me?

	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	glfw.Init()

	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)
	window, err := glfw.CreateWindow(800, 600, "Hello World", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}

	indicies := []uint16{
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
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies)*4, gl.Ptr(indicies), gl.STATIC_DRAW)
	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()

	// copy vertices data into VBO (it needs to be bound first)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(&vertices[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	VShader, err = compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatal(err)
	}

	FShader, err = compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatal(err)
	}

	gl.AttachShader(shaderProgram, VShader)
	gl.AttachShader(shaderProgram, FShader)
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)

	if success == gl.FALSE {
		var infoLog *uint8
		gl.GetProgramInfoLog(shaderProgram, 512, nil, infoLog)
		log.Fatal(infoLog)
	}

	if success == gl.FALSE {
		var logLength int32

		llog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(llog))
		log.Fatal(llog)
	}
	gl.UseProgram(shaderProgram)

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		glfw.PollEvents()

		glfw.WindowHint(glfw.OpenGLProfile, gl.TRIANGLES)

		// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil) // perform draw call
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		// if window.GetKey(glfw.KeyQ) == glfw.Press {
		// 	gl.ClearColor(0.8, 0.5, 0.9, 1.0)
		// } else {
		// 	gl.ClearColor(0.0, 0.30, 0.30, 1.00)
		//
		// }
		// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		window.SwapBuffers()
	}
	gl.DeleteShader(FShader)
	gl.DeleteShader(VShader)
	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.

}

const vertexShader = `
#version 410 core

layout (location = 0) in vec3 aPos;

void main(){
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0f);
}
`

const fragmentShader = `
#version 410 core

out vec4 FragColor;

void main(){
	FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
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

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
