package main

const (
	WIDTH  = 800.0
	HEIGHT = 800.0

	alpha = 7.0  // Controls the influence of the gradient interpolation
	beta  = 7.0  // Another parameter affecting spatial frequency
	n     = 5    // Number of octaves or layers
	seed  = 7777 // Random seed for noise generation

	mapScale int = 1
)

// shader
var (
	vertexShaderSource = `
	#version 410
	in vec3 vp;
	in vec2 texcoord;
	out vec2 TexCoord;
	void main() {
		gl_Position = vec4(vp, 1.0);
		TexCoord = texcoord;
	}
	`
	fragmentShaderSource = `
	#version 410
	in vec2 TexCoord;
	out vec4 frag_colour;
	uniform sampler2D tex;
	void main() {
		frag_colour = texture(tex, TexCoord);
	}
	`
)
