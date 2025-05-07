package examples

import (
	"strconv"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

// ChessboardComponent creates an 8x8 chessboard using manually positioned squares.
func ChessboardComponent(app *mogiApp.App) ui.IComponent {
	boardSize := float32(800.0)
	squareSize := boardSize / 8.0
	children := make([]ui.IComponent, 64) // Pre-allocate slice capacity

	for i := range children {
		x := i % 8
		y := i / 8
		isWhite := (x+y)%2 == 0

		bgColor := color.Black
		if isWhite {
			bgColor = color.White
		}

		if i == 2 {
			bgColor = color.Red
		}

		if i == 8 {
			bgColor = color.Orange
		}

		if i == 63 {
			bgColor = color.Green
		}

		children[i] = app.Container().
			SetID("chess_square_" + strconv.Itoa(i)).
			SetBackgroundColor(bgColor).
			SetSize(math.Vec2f32{
				X: squareSize,
				Y: squareSize,
			}).
			SetPosition(
				ui.Position{
					Type: ui.PositionTypeRelative,
				},
			)
	}

	return app.Container().
		SetID("chessboard_container").
		SetSize(math.Vec2f32{X: boardSize, Y: boardSize}).
		AddChildren(children...)
}
