package gfx

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

const VertexShader = `
#version 300 es

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 position;
in vec3 color;
out vec3 fragmentColor;

void main() {
	gl_Position = projection * camera * model * vec4(position, 1);
  gl_PointSize = 10.0 - gl_Position.z;
  fragmentColor = color;
}
` + "\x00"

const FragmentShader = `
#version 300 es

in lowp vec3 fragmentColor;
out lowp vec3 outColor;

void main()
{
    outColor = fragmentColor;
}
` + "\x00"

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}
