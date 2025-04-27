package main

/*
#cgo LDFLAGS: -L./renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"GoUI/common"
	"GoUI/consts"
	"GoUI/examples"
	"fmt"
	"runtime"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
}

type ComponentRenderer struct {
	Component common.IComponent
	Parent    common.IComponent
}

func goColortoCColorRGBA(color common.ColorRGBA) C.ColorRGBA {
	return C.ColorRGBA{
		r: C.float(color.R),
		g: C.float(color.G),
		b: C.float(color.B),
		a: C.float(color.A),
	}
}

func (cr *ComponentRenderer) Render(renderer unsafe.Pointer) {
	font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 24.0)

	if cr.Component == nil {
		return
	}
	pos := cr.Component.Pos()
	posVec2 := C.Vec2{x: 0, y: 0}

	switch pos.Type {
	case common.PositionTypeAbsolute:
		posVec2.x = C.float(pos.X)
		posVec2.y = C.float(pos.Y)
	default:
		if cr.Parent != nil {
			parentPos := cr.Parent.Pos()
			posVec2.x = C.float(parentPos.X + pos.X)
			posVec2.y = C.float(parentPos.Y + pos.Y)
		} else {
			posVec2.x = C.float(pos.X)
			posVec2.y = C.float(pos.Y)
		}
	}

	switch cr.Component.Type() {
	case common.TContainer:
		container, ok := cr.Component.(*common.Container)
		if !ok {
			return
		}
		// TODO: Add border radius and border width + only render if values are set
		C.draw_rectangle_filled_outline(
			renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(container.Size().X),
				height:   C.float(container.Size().Y),
			},
			goColortoCColorRGBA(container.BackgroundColor),
			goColortoCColorRGBA(container.BorderColor),
		)

	case common.TText:
		text, ok := cr.Component.(*common.Text)
		if !ok {
			return
		}
		cstr := C.CString(text.Text)
		defer C.free(unsafe.Pointer(cstr))
		fontSize := C.float(text.FontSize)
		fontColor := goColortoCColorRGBA(text.Color)
		font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), fontSize)
		defer C.destroy_font(font)
		C.draw_text(
			renderer,
			font,
			cstr,
			posVec2,
			fontColor,
		)

	case common.TButton:
		button, ok := cr.Component.(*common.Button)
		if !ok {
			return
		}
		C.draw_rectangle_filled(
			renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(button.Size().X),
				height:   C.float(button.Size().Y),
			},
			goColortoCColorRGBA(consts.ColorBlue),
		)
		cstr := C.CString(button.Text)
		defer C.free(unsafe.Pointer(cstr))

		defer C.destroy_font(font)
		pos := C.Vec2{x: posVec2.x, y: posVec2.y}
		C.draw_text(
			renderer,
			font,
			cstr,
			pos,
			goColortoCColorRGBA(consts.ColorWhite),
		)
	}

	if cr.Component.Children() != nil {
		for _, child := range cr.Component.Children() {
			childRenderer := &ComponentRenderer{Component: child, Parent: cr.Component}
			childRenderer.Render(renderer)
		}
	}
}

func main() {
	renderer := C.create_renderer(C.int(800), C.int(800), C.CString("Go with C Renderer"))
	defer C.destroy_renderer(renderer)
	font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 24.0)

	if renderer == nil {
		fmt.Println("Failed to create renderer")
		return
	}

	// Create color structure
	bgColor := C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0}
	// rectColor := C.ColorRGBA{r: 1.0, g: 0.0, b: 0.0, a: 1.0}

	// // Create rectangle structure
	// rect := C.Rect{
	// 	position: C.Vec2{x: 100.0, y: 100.0},
	// 	width:    200.0,
	// 	height:   100.0,
	// }

	defer C.destroy_font(font)

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, bgColor)
		// C.draw_rectangle_filled(renderer, rect, rectColor)
		// C.draw_rectangle_outline(renderer, rect, ColorWhite)
		// C.draw_rectangle_filled_outline(renderer, rect, ColorGray, ColorDarkGray)

		// // Example of formatted string
		// smallFont := font                                              // Replace with actual small font if available
		// buffer := fmt.Sprintf("Small Font Example (Size: %.0f)", 24.0) // Replace 60.0 with actual font height if available

		// pos3 := C.Vec2{x: 50.0, y: 250.0}
		// C.draw_text(renderer, smallFont, C.CString(buffer), pos3, ColorGold)

		// pos4 := C.Vec2{x: 50.0, y: 280.0}
		// C.draw_text(renderer, smallFont, C.CString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), pos4, ColorCyan)

		// pos5 := C.Vec2{x: 50.0, y: 300.0}
		// C.draw_text(renderer, smallFont, C.CString("abcdefghijklmnopqrstuvwxyz"), pos5, ColorNavy)

		// pos6 := C.Vec2{x: 50.0, y: 320.0}
		// C.draw_text(renderer, smallFont, C.CString("0123456789 .,:;!?()[]{}"), pos6, ColorOrange)

		// // Create Circle struct and draw
		// circle := C.Circle{
		// 	position: C.Vec2{x: 400.0, y: 300.0},
		// 	radius:   50.0,
		// }
		// C.draw_circle_filled(renderer, circle, ColorViolet)
		// C.draw_circle_outline(renderer, circle, ColorBrown)
		// C.draw_rectangle_filled_outline(renderer, rect, ColorGray, ColorDarkGray)

		// // Create Line struct and draw
		// line := C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 670, y: 200.0},
		// }
		// C.draw_line(renderer, line, ColorRed)
		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 750, y: 200.0},
		// }
		// C.draw_line_thick(renderer, line, ColorGreen, 5.0)
		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 850, y: 200.0},
		// }
		// C.draw_line_dashed(renderer, line, ColorBlue, 5.0, 10.0)

		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 950, y: 200.0},
		// }
		// C.draw_line_dotted(renderer, line, ColorWhite, 2.0)

		// render
		componentRenderer := &ComponentRenderer{Component: examples.BuyNowCardComponent()}
		componentRenderer.Render(renderer)

		C.present_screen(renderer)

		// Don't forget to process events
		C.handle_events(renderer)
	}

	//component demo

	// Set parent-child relationships
	fmt.Println("Exiting")
}
