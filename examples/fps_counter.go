package examples

import (
	"fmt"
	"mogi/common"
	"mogi/consts"
)

func FPSCounterComponent(pos common.Vec2, avgFps float32) common.IComponent {
	fpsTextContent := fmt.Sprintf("Avg. FPS: %.0f", avgFps)

	fpsText := common.NewText(fpsTextContent).
		SetID("fps_counter_text").
		SetFontSize(24).
		SetColor(consts.ColorRed())

	return common.NewContainer().
		SetID("fps_counter").
		SetBackgroundColor(consts.ColorGreen()).
		SetSize(common.Vec2{
			X: 180,
			Y: 35,
		}).
		AddChild(fpsText).
		SetPosition(common.Position{X: pos.X, Y: pos.Y, Type: common.PositionTypeAbsolute})
}
