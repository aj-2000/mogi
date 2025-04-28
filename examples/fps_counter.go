package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"fmt"
)

func FPSCounterComponent(x float32, y float32, avgFps float32) common.IComponent {
	return common.NewContainer(common.ContainerOptions{
		BackgroundColor: consts.ColorGreen,
		Position:        common.Position{X: x, Y: y},
		ID:              "fps_counter",
		Size:            common.Vec2{X: 200, Y: 28},
		Children: []common.IComponent{
			common.NewText(common.TextOptions{
				ID:       "fps_counter_text",
				Content:  fmt.Sprintf("Avg. FPS: %.0f", avgFps),
				FontSize: 24,
				Color:    consts.ColorRed,
				Position: common.Position{X: 10, Y: 0, Type: common.PositionTypeRelative},
				Size:     common.Vec2{X: 180, Y: 30},
			}),
		},
	})
}
