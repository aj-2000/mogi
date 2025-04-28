package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"strconv"
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

		children = append(children, common.NewContainer(common.ContainerOptions{
			BackgroundColor: color,
			BorderColor:     consts.ColorBlack,
			BorderWidth:     1,
			BorderRadius:    0,
			Position: common.Position{
				X:    float32(x) * 100,
				Y:    float32(y) * 100,
				Type: common.PositionTypeRelative,
			},
			ID: "chess_square_" + strconv.Itoa(i),
			Size: common.Vec2{
				X: 100,
				Y: 100,
			},
		}))
		i++
	}

	return common.NewContainer(common.ContainerOptions{
		BackgroundColor: consts.ColorGray,
		BorderColor:     consts.ColorBlack,
		BorderWidth:     0,
		BorderRadius:    0,
		Position: common.Position{
			X:    0,
			Y:    0,
			Type: common.PositionTypeRelative,
		},
		ID:       "main_container",
		Size:     common.Vec2{X: 800, Y: 800},
		Children: children,
	})
}
