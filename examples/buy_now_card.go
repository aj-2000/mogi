package examples

import (
	"log"

	"GoUI/common"
	"GoUI/consts"
)

func BuyNowCardComponent() common.IComponent {
	return common.NewContainer().
		SetID("buy_now_card").
		SetBackgroundColor(consts.ColorBlack()).
		AddChild(
			common.NewContainer().SetBackgroundColor(consts.ColorGreen()).
				SetID("green_rectangle").
				SetSize(common.Vec2{X: 200, Y: 200}),
		).
		AddChild(
			common.NewText("Green Rectangle").
				SetID("green_rectangle_text").
				SetFontSize(24).
				SetColor(consts.ColorWhite()),
		).
		AddChild(
			common.NewText("$19.99").
				SetID("price_text").
				SetFontSize(16).
				SetColor(consts.ColorWhite()),
		).
		AddChild(
			common.NewButton("Buy Now").
				SetID("buy_button").
				SetOnClick(func() { log.Println("Buy Now Clicked!") }).
				SetBackgroundColor(consts.ColorBlue()),
		).
		SetPosition(common.Position{
			X:    320,
			Y:    330,
			Type: common.PositionTypeAbsolute,
		}).
		SetSize(common.Vec2{
			X: 200,
			Y: 280,
		})
}
