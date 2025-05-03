package examples

import (
	"mogi/common"
	"mogi/consts"
	"strconv"
)

func ExampleMarginPaddingBorder() common.IComponent {
	numBoxes := 64
	boxSize := float32(50.0)
	boxes := make([]common.IComponent, numBoxes)
	for i := range boxes {
		boxes[i] = common.NewContainer().
			SetID("box_" + strconv.Itoa(i+1)).
			SetSize(common.Vec2{X: boxSize, Y: boxSize}).
			SetBackgroundColor(common.ColorRGBA{
				R: float32(i) / float32(numBoxes),
				G: float32(i) / float32(numBoxes),
				B: float32(i) / float32(numBoxes),
				A: 1.0,
			}).
			SetMargin(common.Vec2{X: 5, Y: 5}).
			AddChild(
				common.NewText("Box " + strconv.Itoa(i+1)).
					SetID("text_" + strconv.Itoa(i+1)).
					SetFontSize(12).
					SetColor(consts.ColorRed()))
	}

	return common.NewContainer().
		SetID("example_mpb_main_container").
		SetBackgroundColor(consts.ColorOrange()).
		AddChildren(
			boxes...,
		)
}
