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

	font := C.init_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 60)
	defer C.destroy_font(font)

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, bgColor)
		C.draw_rectangle(renderer, rect, rectColor)

		// Prepare line structs as expected by C API
		line1 := C.Line{
			start: C.Vec2{x: 0.0, y: 0.0},
			end:   C.Vec2{x: 800.0, y: 600.0},
		}
		line2 := C.Line{
			start: C.Vec2{x: 800.0, y: 0.0},
			end:   C.Vec2{x: 0.0, y: 600.0},
		}
		lineColor := C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0}
		C.draw_line(renderer, line1, lineColor)
		C.draw_line(renderer, line2, lineColor)

		// Prepare circle struct as expected by C API
		circle := C.Circle{
			position: C.Vec2{x: 400.0, y: 300.0},
			radius:   50.0,
		}
		circleColor := C.ColorRGBA{r: 0.0, g: 1.0, b: 0.0, a: 1.0}
		C.draw_circle(renderer, circle, circleColor)

		C.draw_text(renderer, C.CString("Hello from Go!"), C.Vec2{x: 10.0, y: 10.0}, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0}, font)
		C.draw_text(renderer, C.CString("Press ESC to exit"), C.Vec2{x: 10.0, y: 30.0}, C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0}, font)

		C.present_screen(renderer)

		// Don't forget to process events
		C.handle_events(renderer)
	}

	fmt.Println("Exiting")
}
