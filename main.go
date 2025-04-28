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
	"log"
	"runtime"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
}

type App struct {
	fonts       map[string]*C.FontData
	renderer    unsafe.Pointer
	totalTime   float64
	totalFrames int64
}

func (app *App) LoadFont(path string, size float32) (*C.FontData, error) {
	fontKey := fmt.Sprintf("%s_%f", path, size)
	if font, ok := app.fonts[fontKey]; ok {
		return font, nil
	}

	font := C.load_font(C.CString(path), C.float(size))
	if font == nil {
		return nil, fmt.Errorf("failed to load font: %s", path)
	}

	app.fonts[fontKey] = font
	return font, nil
}

func (app *App) UnloadFont(path string, size float32) {
	fontKey := fmt.Sprintf("%s_%f", path, size)
	if font, ok := app.fonts[fontKey]; ok {
		C.destroy_font(font)
		delete(app.fonts, fontKey)
	}
}

func (app *App) Run(f func(app *App) common.IComponent) {
	if app.renderer == nil {
		log.Fatalln("Renderer is not initialized")
	}

	for C.window_should_close(app.renderer) == 0 {
		deltaTime := app.GetDeltaTime()
		app.totalTime += float64(deltaTime)
		app.totalFrames++
		root := f(app)
		componentRenderer := &ComponentRenderer{Component: root}

		C.clear_screen(app.renderer, C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0})
		componentRenderer.Render(app)
		C.present_screen(app.renderer)
		C.handle_events(app.renderer)
	}
}

func (app *App) GetDeltaTime() float32 {
	deltaTime := C.get_delta_time(app.renderer)
	return float32(deltaTime)
}

// TODO: it's not correct for some reason
func (app *App) GetFPS() float32 {
	deltaTime := app.GetDeltaTime()
	if deltaTime == 0 {
		return 0
	}
	return 1.0 / deltaTime
}

func (app *App) GetAvgFPS() float32 {
	if app.totalFrames == 0 {
		return 0
	}
	avgFPS := float64(app.totalFrames) / app.totalTime
	return float32(avgFPS)
}

func (app *App) SetVSync(vsync bool) {
	if vsync {
		C.set_vsync(app.renderer, 1)
	} else {
		C.set_vsync(app.renderer, 0)
	}
}

func (app *App) Destroy() {
	if app.renderer != nil {
		C.destroy_renderer(app.renderer)
		app.renderer = nil
	}
	for _, font := range app.fonts {
		C.destroy_font(font)
	}
	app.fonts = nil
}

func NewApp(title string, width int, height int) *App {
	renderer := C.create_renderer(C.int(width), C.int(height), C.CString(title))
	if renderer == nil {
		log.Fatalln("Failed to create renderer")
	}

	app := &App{
		fonts:    make(map[string]*C.FontData),
		renderer: renderer,
	}
	app.SetVSync(true)
	return app
}

type ComponentRenderer struct {
	Component common.IComponent
	Parent    common.IComponent
}

func goColorToCColorRGBA(color common.ColorRGBA) C.ColorRGBA {
	return C.ColorRGBA{
		r: C.float(color.R),
		g: C.float(color.G),
		b: C.float(color.B),
		a: C.float(color.A),
	}
}

func (cr *ComponentRenderer) Render(app *App) {
	if cr.Component == nil {
		return
	}

	pos := cr.Component.Pos()
	posVec2 := C.Vec2{x: 0, y: 0}

	switch pos.Type {
	case common.PositionTypeAbsolute:
		posVec2.x = C.float(pos.X)
		posVec2.y = C.float(pos.Y)
	case common.PositionTypeRelative:
		if cr.Parent != nil {
			parentPos := cr.Parent.Pos()
			posVec2.x = C.float(parentPos.X + pos.X)
			posVec2.y = C.float(parentPos.Y + pos.Y)
		} else {
			posVec2.x = C.float(pos.X)
			posVec2.y = C.float(pos.Y)
		}
	}

	switch cr.Component.Kind() {
	case common.ContainerKind:
		container, ok := cr.Component.(*common.Container)
		if !ok {
			return
		}
		C.draw_rectangle_filled_outline(
			app.renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(container.Size().X),
				height:   C.float(container.Size().Y),
			},
			goColorToCColorRGBA(container.BackgroundColor),
			goColorToCColorRGBA(container.BorderColor),
		)

	case common.TextKind:
		text, ok := cr.Component.(*common.Text)
		if !ok {
			return
		}
		cstr := C.CString(text.Content)
		defer C.free(unsafe.Pointer(cstr))

		fontColor := goColorToCColorRGBA(text.Color)
		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", text.FontSize)
		if err != nil {
			log.Fatalf("Failed to load font: %v", err)
		}

		C.draw_text(
			app.renderer,
			font,
			cstr,
			posVec2,
			fontColor,
		)

	case common.ButtonKind:
		button, ok := cr.Component.(*common.Button)
		if !ok {
			return
		}
		C.draw_rectangle_filled(
			app.renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(button.Size().X),
				height:   C.float(button.Size().Y),
			},
			goColorToCColorRGBA(consts.ColorBlue),
		)

		cstr := C.CString(button.Label)
		defer C.free(unsafe.Pointer(cstr))

		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0)
		if err != nil {
			log.Fatalf("Failed to load font: %v", err)
		}

		C.draw_text(
			app.renderer,
			font,
			cstr,
			posVec2,
			goColorToCColorRGBA(consts.ColorWhite),
		)
	}

	// Render Children
	for _, child := range cr.Component.Children() {
		childRenderer := &ComponentRenderer{
			Component: child,
			Parent:    cr.Component,
		}
		childRenderer.Render(app)
	}
}

func main() {
	app := NewApp("GoUI", 800, 800)
	if app == nil {
		log.Fatalln("Failed to create app")
	}
	// TODO: is it needed?
	defer app.Destroy()

	app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0)

	// TODO: should we need to expose app?
	// Define these outside the closure to persist state across frames
	var fpsCounterComponentPos = common.Vec2{X: 10, Y: 10}
	var velocity = common.Vec2{X: 5, Y: 5}

	app.Run(func(app *App) common.IComponent {
		if fpsCounterComponentPos.X > 800-200 {
			velocity.X = -5
		} else if fpsCounterComponentPos.X < 0 {
			velocity.X = 5
		}
		if fpsCounterComponentPos.Y > 800-28 {
			velocity.Y = -5
		} else if fpsCounterComponentPos.Y < 0 {
			velocity.Y = 5
		}

		fpsCounterComponentPos.X += velocity.X
		fpsCounterComponentPos.Y += velocity.Y

		return common.NewContainer(common.ContainerOptions{
			BackgroundColor: consts.ColorCyan,
			BorderWidth:     2,
			BorderRadius:    10,
			Position:        common.Position{X: 0, Y: 0},
			ID:              "main_container",
			Children: []common.IComponent{
				examples.ChessboardComponent(),
				examples.BuyNowCardComponent(),
				examples.FPSCounterComponent(fpsCounterComponentPos.X, fpsCounterComponentPos.Y, app.GetAvgFPS()),
			},
			Size: common.Vec2{X: 800, Y: 800},
		})
	})
}
