package window

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const VERSION_MAJOR = 3
const VERSION_MINOR = 0

type GLFWWindow struct {
	Title         string
	Width, Height int
	Window        *glfw.Window
}

func (w *GLFWWindow) Create() error {
	if err := glfw.Init(); err != nil {
		return err
	}

	// Lock os thread. Important for pollEvents
	runtime.LockOSThread()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err := glfw.CreateWindow(
		w.Width,
		w.Height,
		w.Title, nil, nil)

	if err != nil {
		return err
	}

	// Activate OpenGL context
	window.MakeContextCurrent()

	w.Window = window

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Print OpenGL version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	return nil
}
