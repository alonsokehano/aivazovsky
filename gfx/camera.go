package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	minDistance, maxDistance float32
	distance                 float32
	matrix                   mgl32.Mat4
}

func (camera *Camera) Init() {
	camera.calculate()
}

func (camera *Camera) ZoomOut(distance float32) {
	camera.SetDistance(camera.distance + distance)
}

func (camera *Camera) ZoomIn(distance float32) {
	camera.SetDistance(camera.distance - distance)
}

func (camera *Camera) SetDistance(distance float32) {
	if distance < camera.minDistance {
		camera.distance = camera.minDistance
	} else if distance > camera.maxDistance {
		camera.distance = camera.maxDistance
	} else {
		camera.distance = distance
	}
	camera.calculate()
}

func (camera *Camera) calculate() {
	camera.matrix = mgl32.LookAtV(
		mgl32.Vec3{0, 0, camera.distance},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 1, 0})
}
