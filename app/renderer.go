package app

/*
#cgo LDFLAGS: -L../renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "../renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

type renderer struct {
	ptr            unsafe.Pointer
	fontManager    *fontManager
	textureManager *textureManager
}

func newRenderer(width, height int, title string) *renderer {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	ptr := C.create_renderer(C.int(width), C.int(height), cTitle)
	if ptr == nil {
		log.Fatalln("failed to create renderer")
	}
	return &renderer{
		ptr:            ptr,
		fontManager:    NewFontManager(),
		textureManager: NewTextureManager(),
	}
}

func (r *renderer) setVSync(enabled bool) {
	var flag C.int
	if enabled {
		flag = 1
	}
	C.set_vsync(r.ptr, flag)
}

func (r *renderer) drawText(path string, size float32, text string, pos math.Vec2f32, color color.RGBA) {
	font, err := r.fontManager.load(path, size)
	if err != nil {
		log.Println("failed to load font:", err)
		return
	}
	cPos := C.Vec2{x: C.float(pos.X), y: C.float(pos.Y)}
	C.draw_text(r.ptr, font, C.CString(text), cPos, goColorToCColorRGBA(color))
}

func (r *renderer) drawRectangle(pos, size math.Vec2f32, backgroundColor color.RGBA, borderWidth math.Vec2f32, borderColor color.RGBA, radius float32) {
	cRect := C.Rect{
		position: C.Vec2{x: C.float(pos.X), y: C.float(pos.Y)},
		width:    C.float(size.X),
		height:   C.float(size.Y),
	}
	cFillColor := goColorToCColorRGBA(backgroundColor)
	cBorderWidth := C.Vec2{x: C.float(borderWidth.X), y: C.float(borderWidth.Y)}
	cBorderColor := goColorToCColorRGBA(borderColor)
	C.draw_rectangle_filled_border_rounded(
		r.ptr,
		cRect,
		cFillColor,
		cBorderWidth,
		cBorderColor,
		C.float(radius),
	)
}

// should we expose this to public?
func (r *renderer) drawTexture(textureID C.GLuint, pos, size math.Vec2f32) {
	cRect := C.Rect{
		position: C.Vec2{x: C.float(pos.X), y: C.float(pos.Y)},
		width:    C.float(size.X),
		height:   C.float(size.Y),
	}
	C.draw_texture(r.ptr, textureID, cRect, goColorToCColorRGBA(color.White))
}

func (r *renderer) loadTextureFromMemory(data []byte, w, h, ch int) C.GLuint {
	if len(data) == 0 {
		return 0
	}
	ptr := unsafe.Pointer(&data[0])
	tex := C.load_texture_from_memory((*C.uchar)(ptr), C.int(w), C.int(h), C.int(ch))
	if tex == 0 {
		log.Println("failed to load texture from memory")
	}
	return tex
}

func (r *renderer) clear() {
	C.clear_screen(r.ptr, goColorToCColorRGBA(color.Transparent))
}

func (r *renderer) present() {
	C.present_screen(r.ptr)
}

func (r *renderer) getTime() float32 {
	return float32(C.get_current_time(r.ptr))
}

func (r *renderer) getWindowSize() math.Vec2f32 {
	sz := C.get_window_size(r.ptr)
	return math.Vec2f32{X: float32(sz.x), Y: float32(sz.y)}
}

func (r *renderer) windowShouldClose() bool {
	return C.window_should_close(r.ptr) != 0
}

func (r *renderer) destroy() {
	C.destroy_renderer(r.ptr)
	r.fontManager.destroy()
	r.textureManager.destroy()
	r.ptr = nil
}

func (r *renderer) calculateTextWidth(fontData *FontData, text string) float32 {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return float32(C.calculate_text_width(fontData, cText))
}

func (r *renderer) handleEvents() {
	C.handle_events(r.ptr)
}

func (r *renderer) getMousePos() math.Vec2f32 {
	pos := C.get_cursor_pos(r.ptr)
	return math.Vec2f32{X: float32(pos.x), Y: float32(pos.y)}
}

func (r *renderer) IsMousePressed(button int) bool {
	return C.is_mouse_button_pressed(r.ptr, C.int(button)) != 0
}

func (r *renderer) IsMouseReleased(button int) bool {
	return C.is_mouse_button_released(r.ptr, C.int(button)) != 0
}
