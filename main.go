package main

/*
#cgo LDFLAGS: -L./renderer/lib/Release -lrender -lglfw3 -lgdi32 -static
#include "renderer/include/renderer.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"time"
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

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, bgColor)
		C.draw_rectangle(renderer, rect, rectColor)
		C.present_screen(renderer)

		// Don't forget to process events
		C.handle_events(renderer)

		time.Sleep(16 * time.Millisecond)
	}

	fmt.Println("Exiting")
}
