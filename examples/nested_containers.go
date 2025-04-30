package examples

import (
	"GoUI/common"
	"GoUI/consts"
)

func NestedContainersComponent() common.IComponent {
	return common.NewContainer().
		SetID("nested_containers").
		SetBackgroundColor(consts.ColorBlack()).
		SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
		AddChildren(
			common.NewContainer().
				SetID("container_1").
				SetBackgroundColor(consts.ColorBlue()).
				SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
				AddChildren(
					common.NewContainer().
						SetID("container_1_1").
						SetBackgroundColor(consts.ColorBrown()).
						AddChildren(
							common.NewContainer().
								SetID("container_1_1_1").
								SetSize(common.Vec2{X: 200, Y: 100}).
								SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
								SetBackgroundColor(consts.ColorRed()),
							common.NewContainer().
								SetID("container_1_1_2").
								SetSize(common.Vec2{X: 300, Y: 400}).
								SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
								SetBackgroundColor(consts.ColorPink()),
						),
				),
			common.NewContainer().
				SetID("container_2").
				SetBackgroundColor(consts.ColorGreen()).
				SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
				AddChildren(
					common.NewContainer().
						SetID("container_2_1").
						SetBackgroundColor(consts.ColorOrange()).
						AddChildren(
							common.NewContainer().
								SetID("container_2_1_1").
								SetBackgroundColor(consts.ColorPurple()).
								SetSize(common.Vec2{X: 400, Y: 200}),
							common.NewContainer().
								SetID("container_2_1_2").
								SetBackgroundColor(consts.ColorGray()).
								SetSize(common.Vec2{X: 100, Y: 300}),
						),
				),
		)
}
