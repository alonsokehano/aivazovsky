package main

import (
	"fmt"
	"github.com/alonsokehano/aivazovsky/gfx"
	"github.com/alonsokehano/aivazovsky/window"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const vertexShaderSource = `
#version 300 es

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 position;
in vec3 color;
out vec3 fragmentColor;

void main() {
	gl_Position = projection * camera * model * vec4(position, 1);
  gl_PointSize = 10.0 - gl_Position.z;
  fragmentColor = color;
}
` + "\x00"

const fragmentShaderSource = `
#version 300 es

in lowp vec3 fragmentColor;
out lowp vec3 outColor;

void main()
{
    outColor = fragmentColor;
}
` + "\x00"

func main() {

	// GLFW window preferences
	glfwWindow := window.GLFWWindow{
		Width:  800,
		Height: 600,
		Title:  "Cube",
	}

	err := glfwWindow.Create()
	if err != nil {
		log.Fatalln("Failed to create GLFW window:", err)
	}

	w := glfwWindow.Window

	X := 10
	Y := 10
	Z := 10
	blockConfig := BlockConfig{
		synapses_sens_radius: 2,
		synapses_threshold:   50.0,
		spiking_speed:        10.0,
		relaxation_speed:     10.0,
		relaxation_threshold: 10.0,
	}
	block := Block{x: X, y: Y, z: Z, config: blockConfig}
	block.Initialize()
	vertices := make([]float32, X*Y*Z*3)
	block.Vertices(vertices)
	block.CreatePattern(5, 5, 5, 2, 0.3)
	colors := make([]float32, X*Y*Z*3)
	block.Colors(colors)

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vertexbuffer uint32
	gl.GenBuffers(1, &vertexbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	posAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	var colorbuffer uint32
	gl.GenBuffers(1, &colorbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.DYNAMIC_DRAW)

	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)

	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.226, 0.226, 0.226, 1.0)

	view := gfx.CreateView(glfwWindow.Width, glfwWindow.Height)
	view.Model.Scale(1/float32(block.x), 1/float32(block.y), 1/float32(block.z))
	view.Model.Translate(-0.5, -0.5, -0.5)

	w.SetCursorPosCallback(window.CursorPosCallback(&view))
	w.SetScrollCallback(window.ScrollCallback(&view))

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.UniformMatrix4fv(projectionUniform, 1, false, view.ProjectionUniform())
	gl.UniformMatrix4fv(cameraUniform, 1, false, view.CameraUniform())
	gl.UniformMatrix4fv(modelUniform, 1, false, view.ModelUniform())

	ch := make(chan int)
	go block.Run(ch)

	for !w.ShouldClose() {

		select {
		case <-ch:
			fmt.Println("tick")
			block.Colors(colors)
			gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.DYNAMIC_DRAW)
		default:
			// fmt.Println("    .")
			// time.Sleep(50 * time.Millisecond)
		}

		/* Clear buffers */
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		/* Choose rendering program */
		gl.UseProgram(program)

		/* Bind uniforms */
		gl.UniformMatrix4fv(projectionUniform, 1, false, view.ProjectionUniform())
		gl.UniformMatrix4fv(cameraUniform, 1, false, view.CameraUniform())
		gl.UniformMatrix4fv(modelUniform, 1, false, view.ModelUniform())

		gl.BindVertexArray(vao)

		/* Draw points */
		gl.DrawArrays(gl.POINTS, 0, int32(X*Y*Z))

		// Maintenance
		w.SwapBuffers()

		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

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
