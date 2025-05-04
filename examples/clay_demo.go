package examples

import (
	mogiApp "mogi/app"
	"mogi/color"
	"mogi/math"
	"mogi/ui"
)

func ClayDemoComponent(app *mogiApp.App) ui.IComponent {
	//TODO: use flex layout, padding, margin, grid options, once implemented (Step0: make padding work)
	//TODO: different varaint of same example
	windowSize := app.GetWindowSize()
	padding := float32(16)
	leftColumnWidth := float32((windowSize.X - padding*3) * 0.3)
	rightColumnWidth := float32((windowSize.X - padding*3) * 0.7)
	columnHeight := float32(windowSize.Y - padding*2)
	tileWidth := leftColumnWidth - padding*2
	tileHeight := float32(55)
	firstTileHeight := tileHeight + 15
	imageWidth := firstTileHeight - 2*padding
	mogiWidth := float32(512)
	mogiHeight := float32(512)
	return app.Container().
		SetID("clay_demo_container").
		SetBackgroundColor(color.White).
		SetSize(windowSize).
		SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
		AddChildren(
			app.Container().
				SetID("left_column").
				SetBackgroundColor(color.Skin).
				SetPosition(ui.Position{X: padding, Y: padding, Type: ui.PositionTypeAbsolute}).
				SetSize(math.Vec2f32{X: leftColumnWidth, Y: columnHeight}).
				AddChildren(
					app.Container().
						SetID("tile_1").
						SetPosition(ui.Position{X: 2 * padding, Y: 2 * padding, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: firstTileHeight}).
						SetBackgroundColor(color.Red).
						AddChildren(
							// Stub for a future image
							// TODO: use self size etc to center self
							app.Image("mogi.png").
								SetID("tile_1_image").
								SetSize(math.Vec2f32{X: imageWidth, Y: imageWidth}).
								SetPosition(ui.Position{X: 3 * padding, Y: 3 * padding, Type: ui.PositionTypeAbsolute}),
							app.Text("Mogi - UI library").
								SetID("tile_1_text").
								SetFontSize(firstTileHeight-3*padding).
								SetColor(color.White).
								SetPosition(ui.Position{X: 4*padding + imageWidth, Y: 3.5 * padding, Type: ui.PositionTypeAbsolute}),
						),
					app.Container().
						SetID("tile_2").
						SetPosition(ui.Position{X: 2 * padding, Y: 3*padding + tileHeight + 20, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(color.Orange),
					app.Container().
						SetID("tile_3").
						SetPosition(ui.Position{X: 2 * padding, Y: 4*padding + 2*tileHeight + 20, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(color.Orange),
					app.Container().
						SetID("tile_4").
						SetPosition(ui.Position{X: 2 * padding, Y: 5*padding + 3*tileHeight + 20, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(color.Orange),
					app.Container().
						SetID("tile_5").
						SetPosition(ui.Position{X: 2 * padding, Y: 6*padding + 4*tileHeight + 20, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(color.Orange),
					app.Container().
						SetID("tile_6").
						SetPosition(ui.Position{X: 2 * padding, Y: 7*padding + 5*tileHeight + 20, Type: ui.PositionTypeAbsolute}).
						SetSize(math.Vec2f32{X: tileWidth, Y: tileHeight}).
						SetBackgroundColor(color.Orange),
				),
			app.Container().
				SetID("right_column").
				SetBackgroundColor(color.Skin).
				SetPosition(ui.Position{X: leftColumnWidth + 2*padding, Y: padding, Type: ui.PositionTypeAbsolute}).
				SetSize(math.Vec2f32{X: rightColumnWidth, Y: columnHeight}).
				AddChildren(
					app.Image("mogi.png").SetID("image_1").
						SetSize(math.Vec2f32{X: mogiWidth, Y: mogiHeight}).
						SetPosition(ui.Position{X: 2*padding + leftColumnWidth + (rightColumnWidth / 2) - mogiWidth/2, Y: padding + columnHeight/2 - mogiHeight/2, Type: ui.PositionTypeAbsolute}),
				),
		)
}
