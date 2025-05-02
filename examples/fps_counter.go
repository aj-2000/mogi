package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"fmt"
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
		AddChild(fpsText).
		SetPosition(common.Position{X: pos.X, Y: pos.Y, Type: common.PositionTypeAbsolute})
}
