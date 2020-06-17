package main

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/l0mkaa/ant-colony/pkg/simulation"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	width, height int
)

func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.Scalef(1, -1, 1)
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
}

const alfa = 0xff

var (
	homeColor          = color.RGBA{255, 224, 0, alfa}
	foodColor          = color.RGBA{117, 184, 200, alfa}
	antColor           = color.RGBA{242, 87, 113, alfa}
	dieAntColor        = color.RGBA{148, 138, 139, alfa}
	pheromoneHomeColor = color.RGBA{199, 103, 214, alfa}
	pheromoneFoodColor = color.RGBA{122, 250, 150, alfa}
)

func render(objects [][][]simulation.Object) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LineWidth(1)

	gc := draw2dgl.NewGraphicContext(width, height)
	gc.BeginPath()

	for _, row := range objects {
		for _, obs := range row {
			for _, o := range obs {
				var c color.RGBA
				size := 1.0
				switch o.GetType() {
				case simulation.ANT:
					if o.IsDead() {
						c = dieAntColor
					} else {
						if o.(*simulation.Ant).CarryingFood {
							c = foodColor
						} else {
							c = antColor
						}
					}
					size = 2.0

				case simulation.FOOD:
					c = foodColor
					size = 5.0
				case simulation.HOME:
					c = homeColor
					/*
						case simulation.PHEROMONEFOOD:
							c = pheromoneFoodColor
							a := alfa *
								o.(*simulation.PheromoneFood).GetPower()
							c.A = uint8(a)
							size = 3.0
						case simulation.PHEROMONEHOME:
							c = pheromoneHomeColor
							a := alfa *
								o.(*simulation.PheromoneHome).GetPower()
							c.A = uint8(a)

							size = 3.0
					*/
				}

				draw(gc, o.GetPosition(), c, size)
			}
		}
	}
	gl.Flush()
}

func draw(gc draw2d.GraphicContext, c simulation.Coordinates, color color.RGBA, size float64) {
	gc.SetFillColor(color)
	x, y := float64(c.X), float64(c.Y)
	draw2dkit.Rectangle(gc, x, y, x+size, y+size)
	gc.Fill()
}

func init() {
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Start")
	sim := simulation.NewSimulation(600, 600,
		simulation.Vars{
			AntCount:           100,
			Lifespan:           1000,
			Sight:              7,
			FoodPheremoneDecay: 0.77,
			HomePheremoneDecay: 0.77},
	)
	ok := glfw.Init()
	if !ok {
		panic("glfw.Init()")
	}
	defer glfw.Terminate()
	width, height = 600, 600
	window, err := glfw.CreateWindow(width, height, "Ant Colony", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release {
			return
		}
		x, y := w.GetCursorPosition()
		fmt.Println(x, y)
		sim.AddFood(simulation.Coordinates{X: int(x), Y: int(y)})
	})

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)

	abort := make(chan bool)

	go func() {
		for !window.ShouldClose() {
			abort <- false
		}
		abort <- true
		close(abort)
	}()

	step, obs := sim.Run(abort)

	for i := range step {
		fmt.Println("Step", i)
		o := <-obs
		render(o)
		window.SwapBuffers()
		glfw.PollEvents()
	}

}
