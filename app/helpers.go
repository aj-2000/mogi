package app

/*
#cgo LDFLAGS: -L../renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "../renderer/include/renderer.h"
*/
import "C"
import "mogi/color"

func goColorToCColorRGBA(c color.RGBA) C.ColorRGBA {
	return C.ColorRGBA{
		r: C.float(c.R),
		g: C.float(c.G),
		b: C.float(c.B),
		a: C.float(c.A),
	}
}
