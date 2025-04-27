package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"fmt"
)

func ChessboardComponent() common.IComponent {
	var children []common.IComponent
	i := 0
	for i < 64 {
		x := i % 8
		y := i / 8
		isWhite := (x+y)%2 == 0
		color := consts.ColorBlack
		if isWhite {
			color = consts.ColorWhite
		}

		children = append(children, &common.Container{
			Component: common.Component{
				ComponentType: common.TContainer,
				Pos:           common.Position{X: float32(x) * 100, Y: float32(y) * 100},
				Size:          common.Vec2{X: 100, Y: 100},
				ID:            fmt.Sprintf("box%d", i),
				Children:      nil,
			},
			BackgroundColor: color,
		})
		i++
	}

	return &common.Container{
		Component: common.Component{
			ComponentType: common.TContainer,
			Pos:           common.Position{X: 0, Y: 0},
			Size:          common.Vec2{X: 800, Y: 800},
			ID:            "main_container",
			Children:      children,
		},
		BackgroundColor: consts.ColorGray,
	}
}
