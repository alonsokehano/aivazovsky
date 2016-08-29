package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Projection struct {
	fovy, aspect, near, far float32
	matrix                  mgl32.Mat4
}

func (projection *Projection) Init() {
	projection.matrix = mgl32.Perspective(
		mgl32.DegToRad(projection.fovy),
		projection.aspect,
		projection.near,
		projection.far)
}
