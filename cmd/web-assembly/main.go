package main

import (
	"github.com/l0mkaa/ant-colony/pkg/simulation"

	"fmt"
	"image/color"
	"syscall/js"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

func main() {
	done := make(chan struct{})

	var cvs *canvas.Canvas2d
	var width int
	var height int
	width, height = 300, 300

	sim := simulation.NewSimulation(width, height, simulation.Vars{20, 1000, 7, 0.77, 0.77})

	abort := make(chan bool)
	defer close(abort)
	a := false

	cvs, _ = canvas.NewCanvas2d(false)
	cvs.Create(width, height)

	js.Global().Set("stopSimulation", js.FuncOf(
		func(this js.Value, i []js.Value) interface{} {
			fmt.Println("Stop")
			a = true
			return nil
		}))

	js.Global().Set("runSimulation", js.FuncOf(
		func(this js.Value, i []js.Value) interface{} {
			a = false
			go runSimulation(sim, abort, &a, cvs)
			return nil
		}))

	<-done
}

func runSimulation(sim *simulation.Simulation, abort chan bool, a *bool, cvs *canvas.Canvas2d) {
	go func() {
		for !*a {
			abort <- false
		}
		abort <- true
	}()

	step, obs := sim.Run(abort)
	fmt.Println("Start")
	for range step {
		o := <-obs
		cvs.Start(1, func(gc *draw2dimg.GraphicContext) bool {
			return render(gc, o)
		})
	}
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

func render(gc *draw2dimg.GraphicContext, objects []simulation.Object) bool {
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	gc.BeginPath()
	for _, o := range objects {

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
			size = 2.0
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
	gc.Close()
	return true
}

func draw(gc draw2d.GraphicContext, c simulation.Coordinates, color color.RGBA, size float64) {
	gc.SetFillColor(color)
	x, y := float64(c.X), float64(c.Y)
	draw2dkit.Rectangle(gc, x, y, x+size, y+size)
	gc.Fill()
}
