package examples

import (
	"log"

	"GoUI/common"
	"GoUI/consts"
)

func BuyNowCardComponent() common.IComponent {
	return common.NewContainer().
		SetID("buy_now_card").
		SetBackgroundColor(consts.ColorBlack).
		SetBorderColor(consts.ColorWhite).
		SetBorderWidth(1).
		SetBorderRadius(8).
		SetFlexEnabled(true). // Use Flexbox for children
		SetFlexDirection(common.FlexDirectionColumn).
		SetAlignItems(common.AlignItemsCenter). // Center items horizontally
		SetGap(10).
		// SetPadding(common.EdgeInsets{Top: 15, Bottom: 15, Left: 10, Right: 10}) // Example if padding is added
		// SetFlexBasis(200) // Example: Give card a fixed width if needed
		AddChild(
			common.NewContainer().SetBackgroundColor(consts.ColorGreen).
				SetID("green_rectangle").
				SetSize(common.Vec2{X: 200, Y: 200}).
				SetPosition(common.Position{
					Type: common.PositionTypeRelative,
				}),
		).
		AddChild(
			common.NewText("Green Rectangle").
				SetID("green_rectangle_text").
				SetFontSize(24).
				SetColor(consts.ColorWhite).
				SetPosition(common.Position{
					Type: common.PositionTypeRelative,
				}),
		).
		AddChild(
			common.NewText("$19.99").
				SetID("price_text").
				SetFontSize(16).
				SetColor(consts.ColorWhite).
				SetPosition(common.Position{
					Type: common.PositionTypeRelative,
				}),
		).
		AddChild(
			common.NewButton("Buy Now").
				SetID("buy_button").
				SetOnClick(func() { log.Println("Buy Now Clicked!") }).
				SetBackgroundColor(consts.ColorBlue).
				SetPosition(common.Position{
					Type: common.PositionTypeRelative,
				}),
		).
		SetPosition(common.Position{
			X: 320,
			Y: 330,
		}).
		SetSize(common.Vec2{
			X: 200,
			Y: 280,
		})
}
