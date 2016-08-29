package window

import (
	"fmt"
	"github.com/alonsokehano/aivazovsky/gfx"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func KeyCallback(view *gfx.View) glfw.KeyCallback {
	return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			fmt.Println(key)
		}
	}
}
