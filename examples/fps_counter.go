package examples

import (
	"fmt"
	"mogi/common"
	"mogi/consts"
)

func FPSCounterComponent(pos common.Vec2, getFps func() float32) common.IComponent {
	fps := fmt.Sprintf("FPS: %.0f", getFps())
	fpsText := common.NewText(fps).
		SetID("").
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
