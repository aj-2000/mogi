package examples

import (
	"GoUI/common"
	"GoUI/consts"
)

func ClayDemoComponent(windowSize common.Vec2) common.IComponent {
	//TODO: use flex layout, padding, margin, grid options, once implemented (Step0: make padding work)
	//TODO: different varaint of same example
	padding := float32(16)
	leftColumnWidth := float32((windowSize.X - padding*3) * 0.3)
	rightColumnWidth := float32((windowSize.X - padding*3) * 0.7)
	columnHeight := float32(windowSize.Y - padding*2)
	tileWidth := leftColumnWidth - padding*2
	tileHeight := float32(55)
	firstTileHeight := tileHeight + 15
	imageWidth := firstTileHeight - 2*padding
	return common.NewContainer().
		SetID("clay_demo_container").
		SetBackgroundColor(consts.ColorWhite()).
		SetSize(windowSize).
		SetPosition(common.Position{Type: common.PositionTypeAbsolute}).
		AddChildren(
			common.NewContainer().
				SetID("left_column").
				SetBackgroundColor(consts.ColorSkin()).
				SetPosition(common.Position{X: padding, Y: padding, Type: common.PositionTypeAbsolute}).
				SetSize(common.Vec2{X: leftColumnWidth, Y: columnHeight}).
				AddChildren(
					common.NewContainer().
						SetID("tile_1").
						SetPosition(common.Position{X: 2 * padding, Y: 2 * padding, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: firstTileHeight}).
						SetBackgroundColor(consts.ColorRed()).
						AddChildren(
							// Stub for a future image
							// TODO: use self size etc to center self
							common.NewContainer().
								SetID("tile_1_image").
								SetSize(common.Vec2{X: imageWidth, Y: imageWidth}).
								SetBackgroundColor(consts.ColorBlue()).
								SetPosition(common.Position{X: 3 * padding, Y: 3 * padding, Type: common.PositionTypeAbsolute}),
							common.NewText("Mogi - UI library").
								SetID("tile_1_text").
								SetFontSize(firstTileHeight-3*padding).
								SetColor(consts.ColorWhite()).
								SetPosition(common.Position{X: 4*padding + imageWidth, Y: 3.5 * padding, Type: common.PositionTypeAbsolute}),
						),
					common.NewContainer().
						SetID("tile_2").
						SetPosition(common.Position{X: 2 * padding, Y: 3*padding + tileHeight + 20, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(consts.ColorOrange()),
					common.NewContainer().
						SetID("tile_3").
						SetPosition(common.Position{X: 2 * padding, Y: 4*padding + 2*tileHeight + 20, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(consts.ColorOrange()),
					common.NewContainer().
						SetID("tile_4").
						SetPosition(common.Position{X: 2 * padding, Y: 5*padding + 3*tileHeight + 20, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(consts.ColorOrange()),
					common.NewContainer().
						SetID("tile_5").
						SetPosition(common.Position{X: 2 * padding, Y: 6*padding + 4*tileHeight + 20, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(consts.ColorOrange()),
					common.NewContainer().
						SetID("tile_6").
						SetPosition(common.Position{X: 2 * padding, Y: 7*padding + 5*tileHeight + 20, Type: common.PositionTypeAbsolute}).
						SetSize(common.Vec2{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(consts.ColorOrange()),
				),
			common.NewContainer().
				SetID("right_column").
				SetBackgroundColor(consts.ColorSkin()).
				SetPosition(common.Position{X: leftColumnWidth + 2*padding, Y: padding, Type: common.PositionTypeAbsolute}).
				SetSize(common.Vec2{X: rightColumnWidth, Y: columnHeight}),
		)
}
