package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

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

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	X := 10
	Y := 10
	Z := 10
	blockConfig := BlockConfig{synapses_sens_radius: 2}
	block := Block{config: blockConfig}.NewBlock(X, Y, Z)
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

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.01, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

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
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.STATIC_DRAW)

	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)

	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.226, 0.226, 0.226, 1.0)

	window.SetCursorPosCallback(cursorPosCallback(&model, &camera, &projection))
	window.SetScrollCallback(scrollCallback(&model, &camera, &projection))

	translationMatrix := mgl32.Translate3D(-0.5, -0.5, -0.5)
	scaleMatrix := mgl32.Scale3D(1/float32(block.x), 1/float32(block.y), 1/float32(block.z))
	model = translationMatrix.Mul4(scaleMatrix)

	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)

		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		gl.BindVertexArray(vao)

		gl.DrawArrays(gl.POINTS, 0, int32(X*Y*Z))

		// Maintenance
		window.SwapBuffers()

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

func cursorPosCallback(model, camera, projection *mgl32.Mat4) glfw.CursorPosCallback {
	var x, y float64
	rotate := false
	return func(w *glfw.Window, xpos float64, ypos float64) {
		if glfw.Press == w.GetMouseButton(glfw.MouseButtonRight) {
			if rotate {
				rotation := mgl32.AnglesToQuat(
					float32((ypos-y)/100),
					float32((xpos-x)/100), 0, 1)
				*model = rotation.Mat4().Mul4(*model)
			} else {
				rotate = true
			}
			x = xpos
			y = ypos
		} else {
			rotate = false
		}
	}
}

func scrollCallback(model, camera, projection *mgl32.Mat4) glfw.ScrollCallback {
	zoom := 3.0
	return func(w *glfw.Window, xoffset float64, yoffset float64) {
		zoom -= yoffset * 0.05
		if zoom < 0 {
			zoom = 0
		}
		*camera = mgl32.LookAtV(mgl32.Vec3{0, 0, float32(zoom)}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	}
}
