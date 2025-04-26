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

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, 0.0, 0.0, 0.0, 1.0)
		C.draw_rectangle(renderer, 100.0, 100.0, 200.0, 100.0, 1.0, 0.0, 0.0, 1.0)
		C.present_screen(renderer)
		time.Sleep(16 * time.Millisecond)
	}

	fmt.Println("Exiting")
}
