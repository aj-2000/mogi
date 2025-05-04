package main

/*
#cgo LDFLAGS: -L./renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"mogi/common"
	"mogi/consts"
	"mogi/examples"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"unsafe"
)

func (app *App) Container() *common.Container {
	return common.NewContainer()
}

func (app *App) Text(content string) *common.Text {
	return common.NewText(content)
}

func (app *App) Button(label string) *common.Button {
	return common.NewButton(label)
}

type ComponentState struct {
}

type App struct {
	fonts         map[string]*C.FontData
	textures      map[string]C.GLuint
	renderer      unsafe.Pointer
	totalTime     float64
	totalFrames   int64
	deltaTime     float32
	lastFrameTime float32
	fps           float32
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

func (app *App) LoadTexture(imagePath string) C.GLuint {
	cImagePath := C.CString(imagePath)
	defer C.free(unsafe.Pointer(cImagePath))

	texture := C.load_texture(cImagePath)
	if texture == 0 {
		log.Printf("Failed to load texture: %s", imagePath)
	}
	return texture
}

func (app *App) LoadTextureFromMemory(imageData []byte, width, height, channels int) C.GLuint {
	if len(imageData) == 0 {
		return 0
	}
	ptr := (*C.uchar)(unsafe.Pointer(&imageData[0]))
	texture := C.load_texture_from_memory(ptr, C.int(width), C.int(height), C.int(channels))
	if texture == 0 {
		log.Printf("Failed to load texture from memory")
	}
	return texture
}

func (app *App) UnloadTexture(texture C.GLuint) {
	if texture != 0 {
		C.free_texture(texture)
	}
}

func (app *App) DrawTexture(texture C.GLuint, pos common.Vec2, size common.Vec2) {
	cRect := C.Rect{
		position: C.Vec2{x: C.float(pos.X), y: C.float(pos.Y)},
		width:    C.float(size.X),
		height:   C.float(size.Y),
	}
	C.draw_texture(app.renderer, texture, cRect, goColorToCColorRGBA(consts.ColorWhite()))
}

func (app *App) Run(f func(app *App) common.IComponent) {
	if app.renderer == nil {
		log.Fatalln("Renderer is not initialized")
	}
	// Define these outside the closure to persist state across frames
	le := common.NewLayoutEngine(func(s string, fontSize float32) float32 {
		font, _ := app.LoadFont("JetBrainsMonoNL-Regular.ttf", fontSize)
		return app.CalculateTextWidth(font, s)
	})

	for C.window_should_close(app.renderer) == 0 {
		le.BeginLayout()
		windowSize := app.GetWindowSize()

		current_time := float32(C.get_current_time(app.renderer))
		deltaTime := current_time - app.lastFrameTime
		app.lastFrameTime = current_time
		app.deltaTime = deltaTime
		app.fps = 1.0 / deltaTime

		app.totalTime += float64(app.deltaTime)
		app.totalFrames++
		root := f(app)
		componentRenderer := &ComponentRenderer{Component: root}

		le.Layout(root, common.Vec2{}, windowSize)
		// HandleOnClicks(app, root)

		C.clear_screen(app.renderer, C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0})
		componentRenderer.Render(app)
		C.present_screen(app.renderer)
		C.handle_events(app.renderer)
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

	// Unload all fonts
	for _, font := range app.fonts {
		C.destroy_font(font)
	}
	app.fonts = nil

	// Unload all textures
	for _, texture := range app.textures {
		app.UnloadTexture(texture)
	}
	app.textures = nil

	fmt.Printf("Avg FPS: %f\n", app.GetAvgFPS())
}

func (app *App) GetWindowSize() common.Vec2 {
	if app.renderer == nil {
		log.Fatalln("Renderer is not initialized")
	}
	size := C.get_window_size(app.renderer)
	return common.Vec2{X: float32(size.x), Y: float32(size.y)}
}

func NewApp(title string, width int, height int) *App {
	runtime.LockOSThread()
	renderer := C.create_renderer(C.int(width), C.int(height), C.CString(title))
	if renderer == nil {
		log.Fatalln("Failed to create renderer")
	}

	app := &App{
		fonts:    make(map[string]*C.FontData),
		textures: make(map[string]C.GLuint),
		renderer: renderer,
	}
	app.SetVSync(false)
	return app
}

type ComponentRenderer struct {
	Component common.IComponent
}

func goColorToCColorRGBA(color common.ColorRGBA) C.ColorRGBA {
	return C.ColorRGBA{
		r: C.float(color.R),
		g: C.float(color.G),
		b: C.float(color.B),
		a: C.float(color.A),
	}
}

func HandleOnClicks(app *App, component common.IComponent) {
	if component == nil || component.Display() == common.DisplayNone {
		return
	}

	// Get cursor state once per component
	cursor := C.get_cursor_pos(app.renderer)
	cPos := common.Vec2{X: float32(cursor.x), Y: float32(cursor.y)}
	mouseDown := C.is_mouse_button_pressed(app.renderer, 0) != 0      // held this frame
	mouseReleased := C.is_mouse_button_released(app.renderer, 0) != 0 // just went up this frame

	// If this is a Button, handle its pressed/released logic
	if btn, ok := component.(*common.Button); ok {
		over := btn.IsPointInComponent(cPos)
		btn.MouseOver = over

		// 1) if the cursor goes down inside the button, mark it pressed
		if over && mouseDown {
			btn.Pressed = true
		}

		if btn.Pressed && !over {
			btn.Pressed = false // reset your state
		}

		// 2) if it was pressed and now you see the release, fire **once**:
		if btn.Pressed && mouseReleased {
			btn.Pressed = false // reset your state
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

func (cr *ComponentRenderer) Render(app *App) { // Pass your App struct or RendererPtr directly
	if cr.Component == nil || cr.Component.Display() == common.DisplayNone {
		return
	}

	// Get calculated position and size
	pos := cr.Component.AbsolutePos()
	size := cr.Component.Size()
	cPosVec2 := C.Vec2{x: C.float(pos.X), y: C.float(pos.Y)}
	// print parent with componen  type

	cRect := C.Rect{
		position: cPosVec2,
		width:    C.float(size.X),
		height:   C.float(size.Y),
	}

	cBorderWidth := C.Vec2{x: C.float(cr.Component.Border().X), y: C.float(cr.Component.Border().Y)}
	borderRadius := C.float(cr.Component.BorderRadius())
	borderColor := goColorToCColorRGBA(cr.Component.BorderColor())

	// --- Render based on Kind ---
	switch comp := cr.Component.(type) { // Type switch is often cleaner
	case *common.Container:
		// Draw container background/border
		C.draw_rectangle_filled_border_rounded(
			app.renderer, // Pass C renderer pointer
			cRect,
			goColorToCColorRGBA(comp.BackgroundColor),
			cBorderWidth,
			borderColor,
			borderRadius,
		)

	case *common.Text:
		// Draw text
		cstr := C.CString(comp.Content)
		defer C.free(unsafe.Pointer(cstr))

		fontColor := goColorToCColorRGBA(comp.Color)
		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", comp.FontSize)
		if err != nil {
			log.Printf("Failed to load font during render: %v", err) // Log non-fatally
			// Potentially draw fallback text or nothing
			return
		}
		if font == nil {
			log.Printf("Font pointer is nil for %s", comp.Content)
			return
		}

		// Adjust text position slightly? Often drawing starts at baseline bottom-left.
		// This depends heavily on your C draw_text implementation.
		// For simplicity, using component's top-left for now.
		textPos := cPosVec2
		// textPos.y += C.float(size.Y) // Example if C func expects baseline

		C.draw_text(
			app.renderer,
			font,
			cstr,
			textPos,
			fontColor,
		)

	case *common.Button:
		// Determine background color based on state
		bgColor := comp.BackgroundColor
		if comp.Pressed {
			bgColor = comp.PressedColor
		} else if comp.MouseOver {
			bgColor = comp.HoverColor
		}

		// Draw button background
		C.draw_rectangle_filled_border_rounded(
			app.renderer,
			cRect,
			goColorToCColorRGBA(bgColor),
			cBorderWidth,
			borderColor,
			borderRadius,
		)

		// Draw button label (centered?)
		cstr := C.CString(comp.Label)
		defer C.free(unsafe.Pointer(cstr))

		// TODO: Font loading/caching
		font, err := app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0) // Use a reasonable default or button specific size
		if err != nil {
			log.Printf("Failed to load font during render: %v", err)
			return
		}
		if font == nil {
			log.Printf("Font pointer is nil for button %s", comp.Label)
			return
		}
		// TODO: Calculate centered text position
		// Needs text measurement capabilities from C or Go font library
		textWidth := float32(len(comp.Label)) * 12.0 // Very rough estimate!
		textHeight := float32(24.0)                  // Rough estimate!
		textPos := C.Vec2{
			x: C.float(float32(cPosVec2.x) + (size.X-textWidth)/2.0),
			y: C.float(float32(cPosVec2.y) + (size.Y-textHeight)/2.0), // Adjust based on C func baseline
		}

		C.draw_text(
			app.renderer,
			font,
			cstr,
			textPos,
			goColorToCColorRGBA(comp.TextColor),
		)

	case *common.Image:
		// Draw image
		textureID, ok := app.textures[comp.Path]
		if !ok {
			textureID = app.LoadTexture(comp.Path)
			app.textures[comp.Path] = textureID
		}
		if textureID == 0 {
			log.Printf("Failed to load texture for image: %s", comp.Path)
			return
		}
		pos := common.Vec2{X: float32(cPosVec2.x), Y: float32(cPosVec2.y)}
		app.DrawTexture(textureID, pos, size)
	}

	// --- Render Children ---
	for _, child := range cr.Component.Children() {
		childRenderer := &ComponentRenderer{
			Component: child,
		}
		childRenderer.Render(app) // Pass app/renderer down
	}
}

func main() {

	// --- pprof setup ---
	fCPU, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	defer fCPU.Close()
	if err := pprof.StartCPUProfile(fCPU); err != nil {
		log.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	fMem, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatalf("could not create memory profile: %v", err)
	}
	defer func() {
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(fMem); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
		fMem.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		app := NewApp("mogi", 800, 800)
		if app == nil {
			log.Fatalln("Failed to create app")
		}
		// TODO: is it needed?
		defer app.Destroy()

		app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0)

		children := []common.IComponent{
			examples.ChessboardComponent(),
			examples.BuyNowCardComponent(),
			examples.BoxesOneComponent(),
			examples.BoxesNLevelComponent(3, 3, 100),
			examples.NestedContainersComponent(),
			examples.ClayDemoComponent(app.GetWindowSize()),
			examples.ExampleMarginPaddingBorder(),
		}

		setTabIndex := func(index int) {
			for i, comp := range children {
				display := common.DisplayNone
				if i == index {
					display = common.DisplayBlock
				}
				switch c := comp.(type) {
				case *common.Button:
					c.SetDisplay(display)
				case *common.Text:
					c.SetDisplay(display)
				case *common.Image:
					c.SetDisplay(display)
				case *common.Container:
					c.SetDisplay(display)
				default:
					log.Printf("Unknown component type: %T", c)

				}

			}
		}

		_ = common.NewContainer().
			SetID("tabs").
			SetBackgroundColor(consts.ColorMagenta()).
			AddChildren(
				common.NewContainer().
					SetID("tabs_container").
					SetBackgroundColor(consts.ColorPink()).
					AddChildren(children...),
				common.NewContainer().
					SetID("tab_bar").
					SetBackgroundColor(consts.ColorGray()).
					SetBorderRadius(5).
					SetBorder(common.Vec2{X: 1, Y: 1}).
					SetBorderColor(consts.ColorBlack()).
					SetPadding(common.Vec2{X: 4, Y: 4}).
					SetGap(common.Vec2{X: 3, Y: 3}).
					SetMargin(common.Vec2{X: 3, Y: 3}).
					SetPosition(common.Position{X: 10, Y: 10, Type: common.PositionTypeAbsolute}).
					AddChildren(
						common.NewButton("Chessboard").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(0)
							}),
						common.NewButton("Buy Now").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(1)
							}),
						common.NewButton("Boxes One").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(2)
							}),
						common.NewButton("Boxes N Level").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(3)
							}),
						common.NewButton("Nested Containers").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(4)
							}),
						common.NewButton("Clay Demo").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(5)
							}),
						common.NewButton("Margin Padding Border").
							SetDisplay(common.DisplayInline).
							SetOnClick(func(self *common.Button) {
								setTabIndex(6)
							}),
					),
			).SetSize(app.GetWindowSize())

		app.Run(func(app *App) common.IComponent {

			bgColor := consts.ColorYellow()
			r := app.Container().
				SetID("app_container").
				SetBackgroundColor(bgColor).
				AddChildren( // Add all children at once
					// tabs,
					examples.FPSCounterComponent(common.Vec2{X: app.GetWindowSize().X - 180, Y: 20}, app.GetFPS),
				).
				// SetMargin(common.Vec2{X: 3, Y: 3}).
				SetPadding(common.Vec2{X: 4, Y: 4})

			return r
		})
	}()
	wg.Wait()
}
