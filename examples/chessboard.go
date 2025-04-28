package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"strconv"
)

// ChessboardComponent creates an 8x8 chessboard using manually positioned squares.
func ChessboardComponent() common.IComponent {
	boardSize := float32(800.0)
	squareSize := boardSize / 8.0
	children := make([]common.IComponent, 0, 64) // Pre-allocate slice capacity

	for i := 0; i < 64; i++ {
		x := i % 8
		y := i / 8
		isWhite := (x+y)%2 == 0

		bgColor := consts.ColorBlack
		if isWhite {
			bgColor = consts.ColorWhite
		}

		square := common.NewContainer(). // No ID needed for constructor
							SetID("chess_square_" + strconv.Itoa(i)). // Set optional ID
							SetBackgroundColor(bgColor).
							SetBorderColor(consts.ColorBlack). // Explicit border color
							SetBorderWidth(0).                 // No border width needed if colors touch
							SetPosition(common.Position{       // Explicit position
				X: float32(x) * squareSize,
				Y: float32(y) * squareSize,
			}).
			SetSize(common.Vec2{ // Explicit size
				X: squareSize,
				Y: squareSize,
			})
		children = append(children, square)
	}

	// The main container holding the board squares
	return common.NewContainer().
		SetID("chessboard_container"). // Optional ID for the board itself
		// BackgroundColor defaults to transparent, no need to set if squares cover it
		// BorderColor defaults to transparent
		SetSize(common.Vec2{X: boardSize, Y: boardSize}). // Explicit size for the main board container
		AddChildren(children...)
}
