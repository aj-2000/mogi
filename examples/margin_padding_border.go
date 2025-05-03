package examples

import (
	"mogi/common"
	"mogi/consts"
)

func ExampleMarginPaddingBorder() common.IComponent {

	return common.NewContainer().
		SetID("example_mpb_main_container").
		SetBackgroundColor(consts.ColorWhite()).
		AddChildren(
			common.NewContainer().
				SetID("container_1").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorBlue()),
			common.NewContainer().
				SetID("container_2").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorGreen()),
			common.NewContainer().
				SetID("container_3").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorCyan()),
			common.NewContainer().
				SetID("container_4").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorMagenta()),
			common.NewContainer().
				SetID("container_5").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorRed()),
			common.NewContainer().
				SetID("container_6").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorSkin()),
			common.NewContainer().
				SetID("container_7").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorYellow()),
			common.NewContainer().
				SetID("container_8").
				SetSize(common.Vec2{X: 100, Y: 100}).
				SetBackgroundColor(consts.ColorPurple()),
		)
}
