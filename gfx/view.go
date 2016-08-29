package gfx

const MIN_DISTANCE = 0.01
const MAX_DISTANCE = 10.0
const DEFAULT_DISTANCE = 3.0
const FOVY = 45.0

type View struct {
	Model      Model
	Camera     Camera
	Projection Projection
}

func CreateView(width, height int) View {
	model := Model{}

	camera := Camera{
		minDistance: MIN_DISTANCE,
		maxDistance: MAX_DISTANCE,
		distance:    DEFAULT_DISTANCE,
	}

	projection := Projection{
		fovy:   FOVY,
		aspect: float32(width) / float32(height),
		near:   MIN_DISTANCE,
		far:    MAX_DISTANCE,
	}

	camera.Init()
	projection.Init()
	model.Init()

	return View{
		Model:      model,
		Camera:     camera,
		Projection: projection,
	}
}

func (v *View) ModelUniform() *float32 {
	return &v.Model.matrix[0]
}

func (v *View) CameraUniform() *float32 {
	return &v.Camera.matrix[0]
}

func (v *View) ProjectionUniform() *float32 {
	return &v.Projection.matrix[0]
}
