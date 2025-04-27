package examples

import (
	"GoUI/common"
	"GoUI/consts"
)

func BuyNowCardComponent() common.IComponent {
	var children []common.IComponent

	// Create the product image container
	productImage := &common.Container{
		Component: common.Component{
			ComponentType: common.TContainer,
			Pos:           common.Position{X: 10, Y: 10, Type: common.PositionTypeRelative},
			Size:          common.Vec2{X: 180, Y: 200},
			ID:            "product_image",
			Children:      nil,
		},
		BackgroundColor: consts.ColorGreen,
	}

	// Create the product name text
	productName := &common.Text{
		Component: common.Component{
			ComponentType: common.TText,
			Pos:           common.Position{X: 10, Y: 220, Type: common.PositionTypeRelative},
			Size:          common.Vec2{X: 180, Y: 30},
			ID:            "product_name",
			Children:      nil,
		},
		Text:     "Green Rectangle",
		Color:    consts.ColorWhite,
		FontSize: 16,
	}

	productPrice := &common.Text{
		Component: common.Component{
			ComponentType: common.TText,
			Pos:           common.Position{X: 10, Y: 240, Type: common.PositionTypeRelative},
			Size:          common.Vec2{X: 180, Y: 30},
			ID:            "product_price",
			Children:      nil,
		},
		Text:     "$19.99",
		Color:    consts.ColorWhite,
		FontSize: 16,
	}

	buyNowButton := &common.Button{
		Component: common.Component{
			ComponentType: common.TButton,
			Pos:           common.Position{X: 10, Y: 270, Type: common.PositionTypeRelative},
			Size:          common.Vec2{X: 180, Y: 28},
			ID:            "buy_now_button",
			Children:      nil,
		},
		Text: "Buy Now",
		Callback: func() {
			// Handle buy now button click
			println("Buy Now button clicked!")
		},
	}

	children = append(children, productImage, productName, productPrice, buyNowButton)

	return &common.Container{
		Component: common.Component{
			ComponentType: common.TContainer,
			Pos:           common.Position{X: 320, Y: 330},
			Size:          common.Vec2{X: 200, Y: 310},
			ID:            "buy_now_card",
			Children:      children,
		},
		BackgroundColor: consts.ColorBlack,
		BorderColor:     consts.ColorWhite,
		BorderWidth:     2,
		BorderRadius:    10,
	}
}
