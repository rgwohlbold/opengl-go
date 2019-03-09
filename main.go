package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 640
const windowHeight = 480

func main() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "random window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.DEBUG_OUTPUT)

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
	}

	indices := []uint32{
		0, 1, 2,
		1, 2, 3,
	}

	texCoords := []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	}

	model, err := CreateModelFromData(vertices, indices, texCoords)
	// model, err := CreateModelFromFile("bunny.obj")
	if err != nil {
		panic(err)
	}
	err = model.AddTexture("brick.jpeg", true)
	if err != nil {
		panic(err)
	}
	defer model.Delete()

	entity := Entity{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, 1.0, &model}

	program, err := CreateProgramFromFiles("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}
	defer program.Delete()

	viewMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 1000.0)
	program.Use()
	program.LoadUniformMatrix("viewMatrix", viewMatrix)
	program.Unuse()

	var red float32

	for !window.ShouldClose() {
		gl.ClearColor(1.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		red += 0.01
		if red > 1 {
			red = 0.0
		}

		program.Use()
		program.LoadUniformFloat("red", red)
		entity.Load(&program)
		model.Draw()
		program.Unuse()

		window.SwapBuffers()
		glfw.PollEvents()

		e := gl.GetError()
		if e != gl.NO_ERROR {
			fmt.Println("a gl error has occured")
		}
	}
}
