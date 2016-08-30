package main

import (
	"fmt"
	"github.com/alonsokehano/aivazovsky/gfx"
	"github.com/alonsokehano/aivazovsky/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	// GLFW window preferences
	glfwWindow := window.GLFWWindow{
		Width:  800,
		Height: 600,
		Title:  "Cube",
	}

	err := glfwWindow.Create()
	if err != nil {
		panic(err)
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

	channel := make(chan string)
	go run(&block, channel)

	program, err := gfx.CreateProgram()
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

	glfwWindow.View.Model.Scale(1/float32(block.x), 1/float32(block.y), 1/float32(block.z))
	glfwWindow.View.Model.Translate(-0.5, -0.5, -0.5)

	w.SetCursorPosCallback(window.CursorPosCallback(&glfwWindow))
	w.SetScrollCallback(window.ScrollCallback(&glfwWindow))
	w.SetKeyCallback(window.KeyCallback(&glfwWindow, channel))

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.UniformMatrix4fv(projectionUniform, 1, false, glfwWindow.View.ProjectionUniform())
	gl.UniformMatrix4fv(cameraUniform, 1, false, glfwWindow.View.CameraUniform())
	gl.UniformMatrix4fv(modelUniform, 1, false, glfwWindow.View.ModelUniform())

	for !w.ShouldClose() {

		select {
		case <-channel:
			fmt.Println("tick")
			block.Colors(colors)
			gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.DYNAMIC_DRAW)
		default:
			// Nothing
		}

		/* Clear buffers */
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		/* Bind uniforms */
		gl.UniformMatrix4fv(projectionUniform, 1, false, glfwWindow.View.ProjectionUniform())
		gl.UniformMatrix4fv(cameraUniform, 1, false, glfwWindow.View.CameraUniform())
		gl.UniformMatrix4fv(modelUniform, 1, false, glfwWindow.View.ModelUniform())

		/* Draw points */
		gl.DrawArrays(gl.POINTS, 0, int32(X*Y*Z))

		// Maintenance
		w.SwapBuffers()

		glfw.PollEvents()
	}
}

func run(block *Block, channel chan string) chan string {
	processing := false
	step := 0
	for {
		select {
		case command := <-channel:
			switch command {
			case "toggle":
				if processing {
					fmt.Println("stop")
					processing = false
				} else {
					fmt.Println("start")
					processing = true
					go func() {
						for processing {
							fmt.Println("step", step)
							block.Process()
							step += 1
						}
					}()
				}
			case "start":
				fmt.Println("start")
				processing = true
				go func() {
					for processing {
						fmt.Println("step", step)
						block.Process()
						step += 1
					}
				}()
			case "stop":
				fmt.Println("stop")
				processing = false
			case "step":
				fmt.Println("step", step)
				block.Process()
				step += 1
			}
		}
	}
}
