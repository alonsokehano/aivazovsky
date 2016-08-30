package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func KeyCallback(w *GLFWWindow, notifier chan string) glfw.KeyCallback {
	return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyS:
				notifier <- "toggle"
			case glfw.KeySpace:
				notifier <- "step"
			default:
				// Nothing
			}
		}
	}
}
