package examples

import (
	"fmt"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

func FPSCounterComponent(app *mogiApp.App) ui.IComponent {
	fps := fmt.Sprintf("FPS: %.0f", app.GetFPS())
	windowSize := app.GetWindowSize()
	fpsText := app.Text(fps).
		SetID("fps_text").
		SetFontSize(24).
		SetColor(color.Red)

	// TODO: get self size and position
	// TODO: radius vs border radius
	return app.Container().
		SetID("fps_counter").
		SetBackgroundColor(color.Green).
		SetPadding(math.Vec2f32{X: 5, Y: 5}).
		SetBorderRadius(5).
		AddChild(fpsText).
		SetPosition(ui.Position{X: windowSize.X - 120, Y: 20, Type: ui.PositionTypeAbsolute})
}
