package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

type TerrainCell struct {
	altitude float64
	color    [4]float64
}

func (cell *TerrainCell) calcColor() {
	if cell.altitude < 0 {
		cell.color = [4]float64{0, 0, -mapRange(cell.altitude, 0, -1, 255, 0), 0}
	}
	if cell.altitude >= 0 && cell.altitude < 0.3 {
		cell.color = [4]float64{0, -mapRange(cell.altitude, 0, 0.3, 100, 200), 0, 0}
	}
	if cell.altitude >= 0.3 {
		cell.color = [4]float64{-mapRange(cell.altitude, 0.3, 1, 100, 255), -mapRange(cell.altitude, 0.3, 1, 100, 255), 0, 0}
	}
}

func mapRange(value, fromL, fromH, toL, toH float64) float64 {
	return (value-fromL)*(toH-toL)/(fromH-fromL) + toL
}

func setUp() [][]TerrainCell {
	mapWidth := WIDTH / mapScale
	mapHeight := HEIGHT / mapScale

	pn := perlin.NewPerlin(alpha, beta, n, seed)
	t := make([][]TerrainCell, mapHeight)
	for y := 0; y < mapHeight; y++ {
		t[y] = make([]TerrainCell, mapWidth)
		for x := 0; x < mapWidth; x++ {
			t[y][x].altitude = pn.Noise2D(float64(y)*0.015, float64(x)*0.015)
			t[y][x].calcColor()
		}
	}
	return t
}

var textureID uint32

func createTextureData(t [][]TerrainCell) {
	mapWidth := len(t[0])
	mapHeight := len(t)

	pixelData := make([]byte, 0)
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			for i := 0; i < 4; i++ {
				pixelData = append(pixelData, byte(t[y][x].color[i]))
			}
		}
	}
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(mapWidth), int32(mapHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixelData))
}
func render() {
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	aspectRatio := float32(WIDTH) / float32(HEIGHT)

	vertices := []float32{
		-aspectRatio, -1.0, 0.0,
		aspectRatio, -1.0, 0.0,
		aspectRatio, 1.0, 0.0,
		-aspectRatio, 1.0, 0.0,
	}

	texCoords := []float32{
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
	}

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.Ptr(vertices))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.Ptr(texCoords))
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
}
func run(window *glfw.Window) {
	generatedT := setUp()
	// for a, v := range generatedT {
	// 	for b, x := range v {
	// 		fmt.Printf("[%v][%v]: %v, %v | ", a, b, x.altitude, x.color)
	// 	}
	// }
	createTextureData(generatedT)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "Terrain", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)
		return
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	shaderProgram, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		fmt.Printf("Error creating shader program: %s\n", err)
		return
	}
	gl.UseProgram(shaderProgram)

	//run
	run(window)
}

func createShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := createShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := createShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}
