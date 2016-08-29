package window

import (
	"github.com/alonsokehano/aivazovsky/gfx"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func ScrollCallback(view *gfx.View) glfw.ScrollCallback {
	return func(w *glfw.Window, xoffset float64, yoffset float64) {
		view.Camera.ZoomOut(-float32(yoffset) * 0.05)
	}
}

func CursorPosCallback(view *gfx.View) glfw.CursorPosCallback {
	var x, y float64
	rotate := false
	return func(w *glfw.Window, xpos float64, ypos float64) {
		if glfw.Press == w.GetMouseButton(glfw.MouseButtonRight) {
			if rotate {
				view.Model.Rotate(
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
