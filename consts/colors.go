package consts

import "GoUI/common"

// These cannot be true Go constants because common.ColorRGBA is a struct.
// Instead, use unexported variables and provide exported getter functions for immutability.

var (
	colorRed         = common.ColorRGBA{R: 1, G: 0, B: 0, A: 1}
	colorGreen       = common.ColorRGBA{R: 0, G: 1, B: 0, A: 1}
	colorBlue        = common.ColorRGBA{R: 0, G: 0, B: 1, A: 1}
	colorWhite       = common.ColorRGBA{R: 1, G: 1, B: 1, A: 1}
	colorBlack       = common.ColorRGBA{R: 0, G: 0, B: 0, A: 1}
	colorGray        = common.ColorRGBA{R: 0.5, G: 0.5, B: 0.5, A: 1}
	colorYellow      = common.ColorRGBA{R: 1, G: 1, B: 0, A: 1}
	colorCyan        = common.ColorRGBA{R: 0, G: 1, B: 1, A: 1}
	colorMagenta     = common.ColorRGBA{R: 1, G: 0, B: 1, A: 1}
	colorOrange      = common.ColorRGBA{R: 1, G: 0.5, B: 0, A: 1}
	colorPink        = common.ColorRGBA{R: 1, G: 0.75, B: 0.8, A: 1}
	colorPurple      = common.ColorRGBA{R: 0.5, G: 0, B: 0.5, A: 1}
	colorBrown       = common.ColorRGBA{R: 0.6, G: 0.3, B: 0.2, A: 1}
	colorTransparent = common.ColorRGBA{R: 0, G: 0, B: 0, A: 0}
)

func ColorRed() common.ColorRGBA         { return colorRed }
func ColorGreen() common.ColorRGBA       { return colorGreen }
func ColorBlue() common.ColorRGBA        { return colorBlue }
func ColorWhite() common.ColorRGBA       { return colorWhite }
func ColorBlack() common.ColorRGBA       { return colorBlack }
func ColorGray() common.ColorRGBA        { return colorGray }
func ColorYellow() common.ColorRGBA      { return colorYellow }
func ColorCyan() common.ColorRGBA        { return colorCyan }
func ColorMagenta() common.ColorRGBA     { return colorMagenta }
func ColorOrange() common.ColorRGBA      { return colorOrange }
func ColorPink() common.ColorRGBA        { return colorPink }
func ColorPurple() common.ColorRGBA      { return colorPurple }
func ColorBrown() common.ColorRGBA       { return colorBrown }
func ColorTransparent() common.ColorRGBA { return colorTransparent }
