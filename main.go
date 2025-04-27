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
	fonts    map[string]*C.FontData
	renderer unsafe.Pointer
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

func (app *App) Run(root common.IComponent) {
	if root == nil {
		log.Fatalln("Root component is nil")
	}

	if app.renderer == nil {
		log.Fatalln("Renderer is not initialized")
	}

	componentRenderer := &ComponentRenderer{Component: root}
	for C.window_should_close(app.renderer) == 0 {
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

func (app *App) GetFPS() float32 {
	fps := 1.0 / app.GetDeltaTime()
	return float32(fps)
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

func goColortoCColorRGBA(color common.ColorRGBA) C.ColorRGBA {
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
			app.renderer,
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
		cstr := C.CString(text.Text())
		defer C.free(unsafe.Pointer(cstr))
		fontColor := goColortoCColorRGBA(text.Color)
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

	case common.TButton:
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
			goColortoCColorRGBA(consts.ColorBlue),
		)
		cstr := C.CString(button.Text)
		defer C.free(unsafe.Pointer(cstr))

		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0)
		if err != nil {
			log.Fatalf("Failed to load font: %v", err)
		}

		pos := C.Vec2{x: posVec2.x, y: posVec2.y}
		C.draw_text(
			app.renderer,
			font,
			cstr,
			pos,
			goColortoCColorRGBA(consts.ColorWhite),
		)
	}

	if cr.Component.Children() != nil {
		for _, child := range cr.Component.Children() {
			childRenderer := &ComponentRenderer{Component: child, Parent: cr.Component}
			childRenderer.Render(app)
		}
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

	var children []common.IComponent
	fpsCounterComponent := &common.Container{
		Component: common.Component{
			ComponentType: common.TContainer,
			Pos:           common.Position{X: 10, Y: 10},
			Size:          common.Vec2{X: 130, Y: 28},
			ID:            "fps_counter_background",
			Children: []common.IComponent{
				&common.Text{
					Component: common.Component{
						ComponentType: common.TText,
						Pos:           common.Position{X: 10, Y: 0, Type: common.PositionTypeRelative},
						Size:          common.Vec2{X: 100, Y: 30},
						ID:            "fps_counter_text",
					},
					Text:     func() string { return fmt.Sprintf("FPS: %.2f", app.GetFPS()) },
					Color:    consts.ColorRed,
					FontSize: 24,
				},
			},
		},
		BackgroundColor: consts.ColorGreen,
	}

	children = append(children, examples.ChessboardComponent(), examples.BuyNowCardComponent(), fpsCounterComponent)
	mainContainer := &common.Container{
		Component: common.Component{
			ComponentType: common.TContainer,
			Pos:           common.Position{X: 320, Y: 330},
			Size:          common.Vec2{X: 200, Y: 310},
			ID:            "buy_now_card",
			Children:      children,
		},
		BackgroundColor: consts.ColorBlack,
		BorderColor:     consts.ColorWhite,
		BorderWidth:     2,
		BorderRadius:    10,
	}

	app.Run(mainContainer)

	fmt.Println("Exiting")
}
