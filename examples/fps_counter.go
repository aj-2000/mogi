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
		SetColor(consts.ColorRed).
		SetPosition(common.Position{
			X:    10,
			Y:    5,
			Type: common.PositionTypeRelative,
		})

	return common.NewContainer().
		SetBackgroundColor(consts.ColorGreen).
		SetFlexEnabled(true).
		SetJustifyContent(common.JustifyContentCenter).
		SetAlignItems(common.AlignItemsCenter).
		SetFlexBasis(150).
		SetFlexShrink(0).
		AddChild(fpsText).
		SetPosition(common.Position{
			X:    pos.X,
			Y:    pos.Y,
			Type: common.PositionTypeAbsolute,
		}).
		SetSize(common.Vec2{X: 150, Y: 28})
}
