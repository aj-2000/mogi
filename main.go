package main

/*
#cgo LDFLAGS: -L./renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "renderer/include/renderer.h"
*/
import "C"
import (
	"fmt"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	renderer := C.create_renderer(C.int(800), C.int(600), C.CString("Go with C Renderer"))
	defer C.destroy_renderer(renderer)

	if renderer == nil {
		fmt.Println("Failed to create renderer")
		return
	}

	// Create color structure
	bgColor := C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0}
	rectColor := C.ColorRGBA{r: 1.0, g: 0.0, b: 0.0, a: 1.0}

	// Create rectangle structure
	rect := C.Rect{
		position: C.Vec2{x: 100.0, y: 100.0},
		width:    200.0,
		height:   100.0,
	}

	font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 16.0)
	defer C.destroy_font(font)

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, bgColor)
		C.draw_rectangle(renderer, rect, rectColor)

		// Example of formatted string
		smallFont := font                                              // Replace with actual small font if available
		buffer := fmt.Sprintf("Small Font Example (Size: %.0f)", 16.0) // Replace 60.0 with actual font height if available

		pos3 := C.Vec2{x: 50.0, y: 250.0}
		C.draw_text(renderer, smallFont, C.CString(buffer), pos3, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0})

		pos4 := C.Vec2{x: 50.0, y: 280.0}
		C.draw_text(renderer, smallFont, C.CString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), pos4, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0})

		pos5 := C.Vec2{x: 50.0, y: 300.0}
		C.draw_text(renderer, smallFont, C.CString("abcdefghijklmnopqrstuvwxyz"), pos5, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0})

		pos6 := C.Vec2{x: 50.0, y: 320.0}
		C.draw_text(renderer, smallFont, C.CString("0123456789 .,:;!?()[]{}"), pos6, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0})

		// Create Circle struct and draw
		circle := C.Circle{
			position: C.Vec2{x: 400.0, y: 300.0},
			radius:   50.0,
		}
		C.draw_circle(renderer, circle, C.ColorRGBA{r: 0.0, g: 0.0, b: 1.0, a: 1.0})

		// Create Line struct and draw
		line := C.Line{
			start: C.Vec2{x: 100.0, y: 100.0},
			end:   C.Vec2{x: 200.0, y: 200.0},
		}
		C.draw_line(renderer, line, C.ColorRGBA{r: 1.0, g: 1.0, b: 0.0, a: 1.0})

		C.present_screen(renderer)

		// Don't forget to process events
		C.handle_events(renderer)
	}

	fmt.Println("Exiting")
}
