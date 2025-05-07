package examples

import (
	"strconv"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
	"github.com/aj-2000/mogi/ui"
)

func ExampleMarginPaddingBorder(app *mogiApp.App) ui.IComponent {
	numBoxes := 64
	boxSize := float32(50.0)
	boxes := make([]ui.IComponent, numBoxes)
	for i := range boxes {
		boxes[i] = app.Container().
			SetID("box_" + strconv.Itoa(i+1)).
			SetSize(math.Vec2f32{X: boxSize, Y: boxSize}).
			SetBackgroundColor(color.RGBA{
				R: float32(i) / float32(numBoxes),
				G: float32(i) / float32(numBoxes),
				B: float32(i) / float32(numBoxes),
				A: 1.0,
			}).
			SetMargin(math.Vec2f32{X: 5, Y: 5}).
			AddChild(
				app.Text("Box " + strconv.Itoa(i+1)).
					SetID("text_" + strconv.Itoa(i+1)).
					SetFontSize(12).
					SetColor(color.Red))
	}

	return app.Container().
		SetID("example_mpb_main_container").
		SetBackgroundColor(color.Orange).
		AddChildren(
			boxes...,
		)
}
