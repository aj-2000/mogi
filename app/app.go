package app

/*
#cgo LDFLAGS: -L../renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "../renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"mogi/math"
	"mogi/ui"
	"runtime"
	"unsafe"
)

type App struct {
	renderer      *renderer
	totalTime     float64
	totalFrames   int64
	deltaTime     float32
	lastFrameTime float32
	fps           float32
}

func (app *App) Container() *ui.Container {
	return ui.NewContainer()
}

func (app *App) Text(content string) *ui.Text {
	return ui.NewText(content)
}

func (app *App) Button(label string) *ui.Button {
	return ui.NewButton(label)
}

func (app *App) Image(path string) *ui.Image {
	return ui.NewImage(path)
}

func (app *App) Run(f func(app *App) ui.IComponent) {
	if app.renderer == nil {
		log.Fatalln("Renderer is not initialized")
	}
	// Define these outside the closure to persist state across frames
	le := ui.NewLayoutEngine(func(s string, fontSize float32) float32 {
		font, _ := app.LoadFont("JetBrainsMonoNL-Regular.ttf", fontSize)
		return app.CalculateTextWidth(font, s)
	})

	for !app.renderer.windowShouldClose() {
		le.BeginLayout()
		windowSize := app.GetWindowSize()

		current_time := app.renderer.getTime()
		deltaTime := current_time - app.lastFrameTime
		app.lastFrameTime = current_time
		app.deltaTime = deltaTime
		app.fps = 1.0 / deltaTime

		app.totalTime += float64(app.deltaTime)
		app.totalFrames++
		root := f(app)
		componentRenderer := &ComponentRenderer{Component: root}

		le.Layout(root, math.Vec2f32{}, windowSize)
		// HandleOnClicks(app, root)

		app.renderer.clear()
		componentRenderer.Render(app)
		app.renderer.present()
		app.renderer.handleEvents()
		le.EndLayout()
	}
}

func (app *App) CalculateTextWidth(font *C.FontData, text string) float32 {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	width := C.calculate_text_width(font, cstr)
	return float32(width)
}

// TODO: it's not correct for some reason
func (app *App) GetFPS() float32 {
	deltaTime := app.deltaTime
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
	app.renderer.setVSync(vsync)
}

func (app *App) Destroy() {
	// TODO: we should not expose font manager, text manager here
	app.renderer.destroy()
	fmt.Printf("Avg FPS: %f\n", app.GetAvgFPS())
}

func (app *App) GetWindowSize() math.Vec2f32 {
	return app.renderer.getWindowSize()
}

// TODO: should we expose c font data to public?
func (app *App) LoadFont(path string, size float32) (*FontData, error) {
	return app.renderer.fontManager.load(path, size)
}

func NewApp(width, height int, title string) *App {
	runtime.LockOSThread()
	app := &App{
		renderer: newRenderer(width, height, title),
	}
	app.SetVSync(true)
	return app
}

type ComponentRenderer struct {
	Component ui.IComponent
}

// func HandleOnClicks(app *App, component IComponent) {
// 	if component == nil || component.Display() == DisplayNone {
// 		return
// 	}

// 	// Get cursor state once per component
// 	cursor := C.get_cursor_pos(app.renderer)
// 	cPos := Vec2{X: float32(cursor.x), Y: float32(cursor.y)}
// 	mouseDown := C.is_mouse_button_pressed(app.renderer, 0) != 0      // held this frame
// 	mouseReleased := C.is_mouse_button_released(app.renderer, 0) != 0 // just went up this frame

// 	// If this is a Button, handle its pressed/released logic
// 	if btn, ok := component.(*Button); ok {
// 		over := btn.IsPointInComponent(cPos)
// 		btn.MouseOver = over

// 		// 1) if the cursor goes down inside the button, mark it pressed
// 		if over && mouseDown {
// 			btn.Pressed = true
// 		}

// 		if btn.Pressed && !over {
// 			btn.Pressed = false // reset your state
// 		}

// 		// 2) if it was pressed and now you see the release, fire **once**:
// 		if btn.Pressed && mouseReleased {
// 			btn.Pressed = false // reset your state
// 			if btn.Callback != nil {
// 				btn.Callback(btn)
// 			}
// 		}
// 	}

// 	// recurse into children
// 	for _, child := range component.Children() {
// 		HandleOnClicks(app, child)
// 	}
// }

func (cr *ComponentRenderer) Render(app *App) { // Pass your App struct or RendererPtr directly
	if cr.Component == nil || cr.Component.Display() == ui.DisplayNone {
		return
	}

	pos := cr.Component.AbsolutePos()
	size := cr.Component.Size()

	borderWidth := cr.Component.Border()
	borderRadius := cr.Component.BorderRadius()
	borderColor := cr.Component.BorderColor()
	backgroundColor := cr.Component.BackgroundColor()

	switch comp := cr.Component.(type) { // Type switch is often cleaner
	case *ui.Container:
		app.renderer.drawRectangle(pos, size, backgroundColor, borderWidth, borderColor, borderRadius)

	case *ui.Text:
		app.renderer.drawText("JetBrainsMonoNL-Regular.ttf", comp.FontSize, comp.Content, pos, comp.Color)

	case *ui.Button:
		textHeight := float32(24.0) // Rough estimate!

		// if comp.Pressed {
		// 	bgColor = comp.PressedColor
		// } else if comp.MouseOver {
		// 	bgColor = comp.HoverColor
		// }
		app.renderer.drawRectangle(pos, size, backgroundColor, borderWidth, comp.BorderColor(), borderRadius)
		// TODO: fontsize
		// TODO: text centering
		// TODO: fallback font?
		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0) // Use a reasonable default or button specific size
		if err != nil {
			log.Printf("Failed to load font during render: %v", err)
			return
		}
		textWidth := app.renderer.calculateTextWidth(font, comp.Label)
		offset := size.Sub(*math.NewVec2f32(textWidth, textHeight)).Scale(0.5)
		textPos := pos.Add(*offset)

		app.renderer.drawText("JetBrainsMonoNL-Regular.ttf", textHeight, comp.Label, *textPos, comp.TextColor)

	case *ui.Image:
		// TODO: how to expose Gluint to Go? or use C directly?
		textureID, err := app.renderer.textureManager.load(comp.Path)
		if err != nil {
			log.Printf("Failed to load texture: %v", err)
			return
		}
		app.renderer.drawTexture(textureID, pos, size)
	}

	for _, child := range cr.Component.Children() {
		childRenderer := &ComponentRenderer{
			Component: child,
		}
		childRenderer.Render(app)
	}
}
