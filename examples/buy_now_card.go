package examples

import (
	"GoUI/common"
	"GoUI/consts"
)

func BuyNowCardComponent() common.IComponent {
	productImage := common.NewContainer(common.ContainerOptions{
		ID:              "product_image",
		Position:        common.Position{X: 10, Y: 10, Type: common.PositionTypeRelative},
		Size:            common.Vec2{X: 180, Y: 200},
		BackgroundColor: consts.ColorGreen,
	})

	productName := common.NewText(common.TextOptions{
		ID:       "product_name",
		Content:  "Green Rectangle",
		FontSize: 16,
		Color:    consts.ColorWhite,
		Position: common.Position{X: 10, Y: 220, Type: common.PositionTypeRelative},
		Size:     common.Vec2{X: 180, Y: 30},
	})

	productPrice := common.NewText(common.TextOptions{
		ID:       "product_price",
		Content:  "$19.99",
		FontSize: 16,
		Color:    consts.ColorWhite,
		Position: common.Position{X: 10, Y: 240, Type: common.PositionTypeRelative},
		Size:     common.Vec2{X: 180, Y: 30},
	})

	buyNowButton := common.NewButton(common.ButtonOptions{
		ID:       "buy_now_button",
		Label:    "Buy Now",
		Callback: func() { println("Buy Now button clicked!") },
		Position: common.Position{X: 10, Y: 270, Type: common.PositionTypeRelative},
		Size:     common.Vec2{X: 180, Y: 28},
	})

	children := []common.IComponent{
		productImage,
		productName,
		productPrice,
		buyNowButton,
	}

	return common.NewContainer(common.ContainerOptions{
		ID:              "buy_now_card",
		Position:        common.Position{X: 320, Y: 330},
		Size:            common.Vec2{X: 200, Y: 310},
		BackgroundColor: consts.ColorBlack,
		BorderColor:     consts.ColorWhite,
		BorderWidth:     2,
		BorderRadius:    10,
		Children:        children,
	})
}
