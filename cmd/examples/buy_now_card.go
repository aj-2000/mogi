package examples

import (
	"log"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
	"github.com/aj-2000/mogi/ui"
)

// TODO: how to prevent user to directly use ui.Button etc?

func BuyNowCardComponent(app *mogiApp.App) ui.IComponent {
	return app.Container().
		SetID("buy_now_card").
		SetBackgroundColor(color.Black).
		AddChild(
			app.Container().SetBackgroundColor(color.Green).
				SetID("green_rectangle").
				SetSize(math.Vec2f32{X: 200, Y: 200}),
		).
		AddChild(
			app.Text("Green Rectangle").
				SetID("green_rectangle_text").
				SetFontSize(24).
				SetColor(color.White),
		).
		AddChild(
			app.Text("$19.99").
				SetID("price_text").
				SetFontSize(16).
				SetColor(color.White),
		).
		AddChild(
			app.Button("Buy Now").
				SetID("buy_button").
				SetOnClick(func(_ *ui.Button) { log.Println("Buy Now Clicked!") }).
				SetBackgroundColor(color.Blue),
		).
		SetPosition(ui.Position{
			X:    320,
			Y:    330,
			Type: ui.PositionTypeAbsolute,
		}).
		SetSize(math.Vec2f32{
			X: 200,
			Y: 280,
		})
}
