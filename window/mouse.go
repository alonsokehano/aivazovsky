package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func ScrollCallback(window *Window) glfw.ScrollCallback {
	return func(w *glfw.Window, xoffset float64, yoffset float64) {
		window.View.Camera.ZoomOut(-float32(yoffset) * 0.05)
	}
}

func CursorPosCallback(window *Window) glfw.CursorPosCallback {
	var x, y float64
	rotate := false
	return func(w *glfw.Window, xpos float64, ypos float64) {
		if glfw.Press == w.GetMouseButton(glfw.MouseButtonRight) {
			if rotate {
				window.View.Model.Rotate(
					float32((ypos-y)/100),
					float32((xpos-x)/100),
					0)
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
