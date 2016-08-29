package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	matrix mgl32.Mat4
}

func (model *Model) Init() {
	model.matrix = mgl32.Ident4()
}

func (model *Model) Scale(x, y, z float32) {
	model.matrix = mgl32.Scale3D(x, y, z).Mul4(model.matrix)
}

func (model *Model) Translate(x, y, z float32) {
	model.matrix = mgl32.Translate3D(x, y, z).Mul4(model.matrix)
}

func (model *Model) Rotate(x, y, z float32) {
	model.matrix = mgl32.AnglesToQuat(x, y, z, 1).Mat4().Mul4(model.matrix)
}
