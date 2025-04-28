package examples

import (
	"log"

	"GoUI/common"
	"GoUI/consts"
)

func BuyNowCardComponent() common.IComponent {
	return common.NewContainer().
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
				SetSize(common.Vec2{X: 200, Y: 200}).
				SetPosition(common.Position{
					X:    0,
					Y:    0,
					Type: common.PositionTypeRelative,
				}),
		).
		AddChild(
			common.NewText("Green Rectangle").
				SetFontSize(24).
				SetColor(consts.ColorWhite).
				SetPosition(common.Position{
					X:    10,
					Y:    210,
					Type: common.PositionTypeRelative,
				}).
				SetSize(common.Vec2{
					X: 180,
					Y: 30,
				}),
		).
		AddChild(
			common.NewText("$19.99").
				SetFontSize(16).
				SetColor(consts.ColorWhite).
				SetPosition(common.Position{
					X:    10,
					Y:    240,
					Type: common.PositionTypeRelative,
				}).
				SetSize(common.Vec2{
					X: 180,
					Y: 30,
				}),
		).
		AddChild(
			common.NewButton("Buy Now").
				SetID("buy_button").
				SetOnClick(func() { log.Println("Buy Now Clicked!") }).
				SetBackgroundColor(consts.ColorBlue).
				SetPosition(common.Position{
					X:    10,
					Y:    270,
					Type: common.PositionTypeRelative,
				}).
				SetSize(common.Vec2{
					X: 180,
					Y: 28,
				}),
		).
		SetPosition(common.Position{
			X: 320,
			Y: 330,
		}).
		SetSize(common.Vec2{
			X: 200,
			Y: 310,
		})
}
