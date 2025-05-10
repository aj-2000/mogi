package examples

import (
	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

// TODO: should we receive Vec2f32 or (x, y) as parameters?

func NestedContainersComponent(app *mogiApp.App) ui.IComponent {
	return app.Container().
		SetID("nested_containers").
		SetBackgroundColor(color.Black).
		SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
		AddChildren(
			app.Container().
				SetID("container_1").
				SetBackgroundColor(color.Blue).
				SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
				AddChildren(
					app.Container().
						SetID("container_1_1").
						SetWidthPercent(50).
						// TODO : fix percent precedence*
						SetBackgroundColor(color.Brown).
						AddChildren(
							app.Container().
								SetID("container_1_1_1").
								SetSize(math.Vec2f32{X: 200, Y: 100}).
								SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
								SetBackgroundColor(color.Red),
							app.Container().
								SetID("container_1_1_2").
								SetSize(math.Vec2f32{X: 300, Y: 400}).
								SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
								SetBackgroundColor(color.Pink),
						),
				),
			app.Container().
				SetID("container_2").
				SetBackgroundColor(color.Green).
				SetPosition(ui.Position{Type: ui.PositionTypeAbsolute}).
				AddChildren(
					app.Container().
						SetID("container_2_1").
						SetBackgroundColor(color.Orange).
						AddChildren(
							app.Container().
								SetID("container_2_1_1").
								SetBackgroundColor(color.Purple).
								SetSize(math.Vec2f32{X: 400, Y: 200}),
							app.Container().
								SetID("container_2_1_2").
								SetBackgroundColor(color.Gray).
								SetSize(math.Vec2f32{X: 100, Y: 300}),
						),
				),
		)
}
