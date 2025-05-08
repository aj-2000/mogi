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
	"runtime"
	"sort"
	"unsafe"

	"github.com/aj-2000/mogi/color"

	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

type App struct {
	renderer      *renderer
	totalTime     float64
	totalFrames   int64
	deltaTime     float32
	lastFrameTime float32
	fps           float32
	le            *ui.LayoutEngine
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
	// TODO: optimize fps calculation
	for !app.renderer.windowShouldClose() {
		app.le.BeginLayout()
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
		app.le.AssignIDsRecursive(root)
		if app.totalFrames != 1 {
			// should not run on the first frame
			app.le.CopyStateToComponentsRecursive(root)
		}
		app.le.Layout(root, math.Vec2f32{}, windowSize)
		// Logic that requires state from the previous frame
		HandleOnClicks(app, root)
		app.le.CopyStateFromComponentsRecursive(root)

		app.renderer.clear()
		componentRenderer.Render(app)
		app.renderer.present()
		app.renderer.handleEvents()
		app.le.EndLayout()
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
	// TODO: make it cleaner
	var app *App
	app = &App{
		le: ui.NewLayoutEngine(func(s string, fontSize float32) float32 {
			font, _ := app.LoadFont("JetBrainsMonoNL-Regular.ttf", fontSize)
			return app.CalculateTextWidth(font, s)
		}),
		renderer: newRenderer(width, height, title),
	}
	app.SetVSync(true)
	return app
}

type ComponentRenderer struct {
	Component ui.IComponent
}

func HandleOnClicks(app *App, component ui.IComponent) {
	if component == nil || component.Display() == ui.DisplayNone {
		return
	}

	// Get cursor state once per component
	// TODO: move it to wrapper?
	cursorPos := app.renderer.getMousePos()
	mouseDown := app.renderer.IsMousePressed(0)      // held this frame
	mouseReleased := app.renderer.IsMouseReleased(0) // just went up this frame

	// If this is a Button, handle its pressed/released logic
	if btn, ok := component.(*ui.Button); ok {
		over := btn.IsPointInsideComponent(cursorPos)
		btn.IsMouseOver = over

		// 1) if the cursor goes down inside the button, mark it pressed
		if over && mouseDown {
			btn.IsPressed = true
		}

		if btn.IsPressed && !over {
			btn.IsPressed = false // reset your state
		}

		// 2) if it was pressed and now you see the release, fire **once**:
		if btn.IsPressed && mouseReleased {
			btn.IsPressed = false // reset your state
			if btn.Callback != nil {
				btn.Callback(btn)
			}
		}
	}

	// recurse into children
	for _, child := range component.Children() {
		HandleOnClicks(app, child)
	}
}

func (app *App) GetMousePos() math.Vec2f32 {
	return app.renderer.getMousePos()
}

type RenderCommandKind int

const (
	RenderCommandNone RenderCommandKind = iota
	RenderCommandDrawRectangle
	RenderCommandDrawText
	RenderCommandDrawTexture
)

type RenderCommand struct {
	Kind            RenderCommandKind
	Pos             math.Vec2f32
	Size            math.Vec2f32
	Color           color.RGBA
	Font            *C.FontData
	Text            string
	BorderWidth     math.Vec2f32
	BorderColor     color.RGBA
	BorderRadius    float32
	BackgroundColor color.RGBA
	ZIndex          int
	HoverColor      color.RGBA
	PressedColor    color.RGBA
	FontSize        float32
	Path            string
	Display         ui.Display
}

type RenderCommandArray = []RenderCommand

func (cr *ComponentRenderer) GenerateRenderCommands(app *App) RenderCommandArray {
	if cr.Component == nil || cr.Component.Display() == ui.DisplayNone {
		return nil
	}

	pos := cr.Component.AbsolutePos()
	size := cr.Component.Size()

	borderWidth := cr.Component.Border()
	borderRadius := cr.Component.BorderRadius()
	borderColor := cr.Component.BorderColor()
	backgroundColor := cr.Component.BackgroundColor()
	zIndex := cr.Component.AbsoluteZIndex()
	var commands RenderCommandArray

	switch comp := cr.Component.(type) {
	case *ui.Container:
		commands = append(commands, RenderCommand{
			Kind:            RenderCommandDrawRectangle,
			Pos:             pos,
			Size:            size,
			Color:           backgroundColor,
			BorderWidth:     borderWidth,
			BorderColor:     borderColor,
			BorderRadius:    borderRadius,
			ZIndex:          zIndex,
			Display:         comp.Display(),
			BackgroundColor: backgroundColor,
		})

	case *ui.Text:
		commands = append(commands, RenderCommand{
			Kind:     RenderCommandDrawText,
			Text:     comp.Content,
			Color:    comp.Color,
			Pos:      pos,
			Display:  comp.Display(),
			FontSize: comp.FontSize,
			ZIndex:   zIndex})

	case *ui.Button:
		if comp.IsPressed {
			backgroundColor = comp.PressedColor
		} else if comp.IsMouseOver {
			backgroundColor = comp.HoverColor
		}
		buttonCommand := RenderCommand{
			Kind:            RenderCommandDrawRectangle,
			Color:           backgroundColor,
			ZIndex:          zIndex,
			Text:            "",
			Pos:             pos,
			Size:            size,
			HoverColor:      comp.HoverColor,
			PressedColor:    comp.PressedColor,
			BorderWidth:     borderWidth,
			BorderColor:     borderColor,
			Display:         comp.Display(),
			BorderRadius:    borderRadius,
			BackgroundColor: backgroundColor,
		}
		commands = append(commands, buttonCommand)
		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", comp.FontSize())
		if err != nil {
			log.Printf("Failed to load font during render: %v", err)
			return nil
		}
		textWidth := app.renderer.calculateTextWidth(font, comp.Label)
		offset := size.Sub(*math.NewVec2f32(textWidth, comp.FontSize())).Scale(0.5)
		textPos := *pos.Add(*offset)
		commands = append(commands, RenderCommand{
			Kind:     RenderCommandDrawText,
			Text:     comp.Label,
			Color:    comp.TextColor,
			Pos:      textPos,
			Display:  comp.Display(),
			FontSize: comp.FontSize(),
			ZIndex:   zIndex + 1})

	case *ui.Image:

		imageCommand := RenderCommand{
			Kind:    RenderCommandDrawTexture,
			Path:    comp.Path,
			Pos:     pos,
			Display: comp.Display(),
			Size:    size,
			ZIndex:  zIndex,
		}
		commands = append(commands, imageCommand)
	}

	for _, child := range cr.Component.Children() {
		childRenderer := &ComponentRenderer{Component: child}
		childCommands := childRenderer.GenerateRenderCommands(app)
		commands = append(commands, childCommands...)
	}

	return commands
}

func (cr *ComponentRenderer) Render(app *App) {
	commands := cr.GenerateRenderCommands(app)
	// Sort commands by ZIndex
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].ZIndex < commands[j].ZIndex
	})
	for _, command := range commands {
		if command.Display == ui.DisplayNone {
			continue
		}
		switch command.Kind {
		case RenderCommandDrawRectangle:
			app.renderer.drawRectangle(command.Pos, command.Size, command.BackgroundColor, command.BorderWidth, command.BorderColor, command.BorderRadius)

		case RenderCommandDrawText:
			app.renderer.drawText("JetBrainsMonoNL-Regular.ttf", command.FontSize, command.Text, command.Pos, command.Color)

		case RenderCommandDrawTexture:
			if command.Path == "" {
				log.Println("Texture path is empty, skipping texture render")
				continue
			}
			textureID, err := app.renderer.textureManager.load(command.Path)
			if err != nil {
				log.Printf("Failed to load texture: %v", err)
				return
			}
			app.renderer.drawTexture(textureID, command.Pos, command.Size)
		default:
			log.Printf("Unknown render command kind: %v", command.Kind)
		}
	}
}
