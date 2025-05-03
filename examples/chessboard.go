package examples

import (
	"mogi/common"
	"mogi/consts"
	"strconv"
)

// ChessboardComponent creates an 8x8 chessboard using manually positioned squares.
func ChessboardComponent() common.IComponent {
	boardSize := float32(800.0)
	squareSize := boardSize / 8.0
	children := make([]common.IComponent, 64) // Pre-allocate slice capacity

	for i := range children {
		x := i % 8
		y := i / 8
		isWhite := (x+y)%2 == 0

		bgColor := consts.ColorBlack
		if isWhite {
			bgColor = consts.ColorWhite
		}

		if i == 2 {
			bgColor = consts.ColorRed // Highlight the square at index 2 (for example)
		}

		if i == 8 {
			bgColor = consts.ColorOrange // Highlight the square at index 2 (for example)
		}

		if i == 63 {
			bgColor = consts.ColorGreen // Highlight the square at index 2 (for example)
		}

		children[i] = common.NewContainer(). // No ID needed for constructor
							SetID("chess_square_" + strconv.Itoa(i)). // Set optional ID
							SetBackgroundColor(bgColor()).
							SetSize(common.Vec2{ // Explicit size
				X: squareSize,
				Y: squareSize,
			}).
			SetPosition(
				common.Position{
					Type: common.PositionTypeRelative,
				},
			)
	}

	// The main container holding the board squares
	return common.NewContainer().
		SetID("chessboard_container"). // Optional ID for the board itself
		// BackgroundColor defaults to transparent, no need to set if squares cover it
		// BorderColor defaults to transparent
		SetSize(common.Vec2{X: boardSize, Y: boardSize}). // Explicit size for the main board container
		AddChildren(children...)
}
