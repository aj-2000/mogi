package main

/*
#cgo LDFLAGS: -L./renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
}

// different colors
var ColorRed = C.ColorRGBA{r: 1.0, g: 0.0, b: 0.0, a: 1.0}
var ColorGreen = C.ColorRGBA{r: 0.0, g: 1.0, b: 0.0, a: 1.0}
var ColorBlue = C.ColorRGBA{r: 0.0, g: 0.0, b: 1.0, a: 1.0}
var ColorWhite = C.ColorRGBA{r: 1.0, g: 1.0, b: 1.0, a: 1.0}
var ColorBlack = C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0}
var ColorYellow = C.ColorRGBA{r: 1.0, g: 1.0, b: 0.0, a: 1.0}
var ColorCyan = C.ColorRGBA{r: 0.0, g: 1.0, b: 1.0, a: 1.0}
var ColorMagenta = C.ColorRGBA{r: 1.0, g: 0.0, b: 1.0, a: 1.0}
var ColorGray = C.ColorRGBA{r: 0.5, g: 0.5, b: 0.5, a: 1.0}
var ColorDarkGray = C.ColorRGBA{r: 0.25, g: 0.25, b: 0.25, a: 1.0}
var ColorLightGray = C.ColorRGBA{r: 0.75, g: 0.75, b: 0.75, a: 1.0}
var ColorOrange = C.ColorRGBA{r: 1.0, g: 0.5, b: 0.0, a: 1.0}
var ColorPurple = C.ColorRGBA{r: 0.5, g: 0.0, b: 1.0, a: 1.0}
var ColorPink = C.ColorRGBA{r: 1.0, g: 0.0, b: 0.5, a: 1.0}
var ColorBrown = C.ColorRGBA{r: 0.5, g: 0.25, b: 0.0, a: 1.0}
var ColorGold = C.ColorRGBA{r: 1.0, g: 0.84, b: 0.0, a: 1.0}
var ColorSilver = C.ColorRGBA{r: 0.75, g: 0.75, b: 0.75, a: 1.0}
var ColorTeal = C.ColorRGBA{r: 0.0, g: 0.5, b: 0.5, a: 1.0}
var ColorNavy = C.ColorRGBA{r: 0.0, g: 0.0, b: 0.5, a: 1.0}
var ColorOlive = C.ColorRGBA{r: 0.5, g: 0.5, b: 0.0, a: 1.0}
var ColorMaroon = C.ColorRGBA{r: 0.5, g: 0.0, b: 0.0, a: 1.0}
var ColorLime = C.ColorRGBA{r: 0.0, g: 1.0, b: 0.0, a: 1.0}
var ColorAqua = C.ColorRGBA{r: 0.0, g: 1.0, b: 1.0, a: 1.0}
var ColorFuchsia = C.ColorRGBA{r: 1.0, g: 0.0, b: 1.0, a: 1.0}
var ColorCoral = C.ColorRGBA{r: 1.0, g: 0.5, b: 0.31, a: 1.0}
var ColorKhaki = C.ColorRGBA{r: 0.94, g: 0.9, b: 0.55, a: 1.0}
var ColorSalmon = C.ColorRGBA{r: 0.98, g: 0.5, b: 0.45, a: 1.0}
var ColorPeach = C.ColorRGBA{r: 1.0, g: 0.8, b: 0.65, a: 1.0}
var ColorMint = C.ColorRGBA{r: 0.68, g: 1.0, b: 0.68, a: 1.0}
var ColorPlum = C.ColorRGBA{r: 0.87, g: 0.63, b: 0.87, a: 1.0}
var ColorSlate = C.ColorRGBA{r: 0.44, g: 0.5, b: 0.56, a: 1.0}
var ColorSteel = C.ColorRGBA{r: 0.27, g: 0.51, b: 0.71, a: 1.0}
var ColorIndigo = C.ColorRGBA{r: 0.29, g: 0.0, b: 0.51, a: 1.0}
var ColorViolet = C.ColorRGBA{r: 0.93, g: 0.51, b: 0.93, a: 1.0}
var ColorThistle = C.ColorRGBA{r: 0.85, g: 0.75, b: 0.85, a: 1.0}
var ColorWheat = C.ColorRGBA{r: 0.96, g: 0.87, b: 0.7, a: 1.0}
var ColorTan = C.ColorRGBA{r: 0.82, g: 0.71, b: 0.55, a: 1.0}
var ColorChocolate = C.ColorRGBA{r: 0.82, g: 0.41, b: 0.12, a: 1.0}
var ColorSienna = C.ColorRGBA{r: 0.65, g: 0.16, b: 0.16, a: 1.0}
var ColorPeru = C.ColorRGBA{r: 0.8, g: 0.52, b: 0.25, a: 1.0}
var ColorBurlywood = C.ColorRGBA{r: 0.87, g: 0.72, b: 0.53, a: 1.0}

// UI container

type PositionType int

const (
	PositionTypeAbsolute PositionType = iota
	PositionTypeRelative
)

type Position struct {
	X, Y float32
	Type PositionType
}

type IComponent interface {
	Type() ComponentType
	Pos() Position
	Size() Vec2
	ID() string
	Parent() IComponent
	Children() []IComponent
}

type Component struct {
	ComponentType ComponentType
	Pos           Position
	Size          Vec2
	ID            string
	Parent        IComponent
	Children      []IComponent
}

type ComponentType int

const (
	TContainer ComponentType = iota
	TText      ComponentType = iota
	TButton    ComponentType = iota
)

type Vec2 struct {
	X, Y float32
}
type Container struct {
	Component
}

func (c *Container) Type() ComponentType {
	return c.Component.ComponentType
}

func (c *Container) Pos() Position {
	return c.Component.Pos
}

func (c *Container) Size() Vec2 {
	return c.Component.Size
}

func (c *Container) ID() string {
	return c.Component.ID
}

func (c *Container) Parent() IComponent {
	return c.Component.Parent
}

func (c *Container) Children() []IComponent {
	return c.Component.Children
}

type Text struct {
	Component
	Text string
}

func (t *Text) Type() ComponentType {
	return t.Component.ComponentType
}

func (t *Text) Pos() Position {
	return t.Component.Pos
}

func (t *Text) Size() Vec2 {
	return t.Component.Size
}

func (t *Text) ID() string {
	return t.Component.ID
}

func (t *Text) Parent() IComponent {
	return t.Component.Parent
}

func (t *Text) Children() []IComponent {
	return t.Component.Children
}

type Button struct {
	Component
	Text      string
	Callback  func()
	Pressed   bool
	Released  bool
	MouseOver bool
}

func (b *Button) Type() ComponentType {
	return b.Component.ComponentType
}

func (b *Button) Pos() Position {
	return b.Component.Pos
}

func (b *Button) Size() Vec2 {
	return b.Component.Size
}

func (b *Button) ID() string {
	return b.Component.ID
}

func (b *Button) Parent() IComponent {
	return b.Component.Parent
}

func (b *Button) Children() []IComponent {
	return b.Component.Children
}

type ComponentRenderer struct {
	Component IComponent
}

func (cr *ComponentRenderer) Render(renderer unsafe.Pointer) {
	font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 24.0)

	if cr.Component == nil {
		return
	}
	pos := cr.Component.Pos()
	posVec2 := C.Vec2{x: 0, y: 0}

	switch pos.Type {
	case PositionTypeAbsolute:
		posVec2.x = C.float(pos.X)
		posVec2.y = C.float(pos.Y)
	default:
		if cr.Component.Parent() != nil {
			parentPos := cr.Component.Parent().Pos()
			posVec2.x = C.float(parentPos.X + pos.X)
			posVec2.y = C.float(parentPos.Y + pos.Y)
		} else {
			posVec2.x = C.float(pos.X)
			posVec2.y = C.float(pos.Y)
		}
	}

	switch cr.Component.Type() {
	case TContainer:
		container, ok := cr.Component.(*Container)
		if !ok {
			return
		}
		C.draw_rectangle_filled(
			renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(container.Size().X),
				height:   C.float(container.Size().Y),
			},
			ColorGray,
		)

	case TText:
		text, ok := cr.Component.(*Text)
		if !ok {
			return
		}
		cstr := C.CString(text.Text)
		defer C.free(unsafe.Pointer(cstr))

		defer C.destroy_font(font)
		C.draw_text(
			renderer,
			font,
			cstr,
			posVec2,
			ColorWhite,
		)

	case TButton:
		button, ok := cr.Component.(*Button)
		if !ok {
			return
		}
		C.draw_rectangle_filled(
			renderer,
			C.Rect{
				position: posVec2,
				width:    C.float(button.Size().X),
				height:   C.float(button.Size().Y),
			},
			ColorBlue,
		)
		cstr := C.CString(button.Text)
		defer C.free(unsafe.Pointer(cstr))

		defer C.destroy_font(font)
		pos := C.Vec2{x: posVec2.x, y: posVec2.y}
		C.draw_text(
			renderer,
			font,
			cstr,
			pos,
			ColorWhite,
		)
	}

	if cr.Component.Children() != nil {
		for _, child := range cr.Component.Children() {
			childRenderer := &ComponentRenderer{Component: child}
			childRenderer.Render(renderer)
		}
	}
}

func main() {
	renderer := C.create_renderer(C.int(800), C.int(600), C.CString("Go with C Renderer"))
	defer C.destroy_renderer(renderer)
	font := C.load_font(C.CString("JetBrainsMonoNL-Regular.ttf"), 24.0)

	if renderer == nil {
		fmt.Println("Failed to create renderer")
		return
	}

	// Create color structure
	bgColor := C.ColorRGBA{r: 0.0, g: 0.0, b: 0.0, a: 1.0}
	// rectColor := C.ColorRGBA{r: 1.0, g: 0.0, b: 0.0, a: 1.0}

	// // Create rectangle structure
	// rect := C.Rect{
	// 	position: C.Vec2{x: 100.0, y: 100.0},
	// 	width:    200.0,
	// 	height:   100.0,
	// }

	defer C.destroy_font(font)

	for C.window_should_close(renderer) == 0 {
		C.clear_screen(renderer, bgColor)
		// C.draw_rectangle_filled(renderer, rect, rectColor)
		// C.draw_rectangle_outline(renderer, rect, ColorWhite)
		// C.draw_rectangle_filled_outline(renderer, rect, ColorGray, ColorDarkGray)

		// // Example of formatted string
		// smallFont := font                                              // Replace with actual small font if available
		// buffer := fmt.Sprintf("Small Font Example (Size: %.0f)", 24.0) // Replace 60.0 with actual font height if available

		// pos3 := C.Vec2{x: 50.0, y: 250.0}
		// C.draw_text(renderer, smallFont, C.CString(buffer), pos3, ColorGold)

		// pos4 := C.Vec2{x: 50.0, y: 280.0}
		// C.draw_text(renderer, smallFont, C.CString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), pos4, ColorCyan)

		// pos5 := C.Vec2{x: 50.0, y: 300.0}
		// C.draw_text(renderer, smallFont, C.CString("abcdefghijklmnopqrstuvwxyz"), pos5, ColorNavy)

		// pos6 := C.Vec2{x: 50.0, y: 320.0}
		// C.draw_text(renderer, smallFont, C.CString("0123456789 .,:;!?()[]{}"), pos6, ColorOrange)

		// // Create Circle struct and draw
		// circle := C.Circle{
		// 	position: C.Vec2{x: 400.0, y: 300.0},
		// 	radius:   50.0,
		// }
		// C.draw_circle_filled(renderer, circle, ColorViolet)
		// C.draw_circle_outline(renderer, circle, ColorBrown)
		// C.draw_rectangle_filled_outline(renderer, rect, ColorGray, ColorDarkGray)

		// // Create Line struct and draw
		// line := C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 670, y: 200.0},
		// }
		// C.draw_line(renderer, line, ColorRed)
		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 750, y: 200.0},
		// }
		// C.draw_line_thick(renderer, line, ColorGreen, 5.0)
		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 850, y: 200.0},
		// }
		// C.draw_line_dashed(renderer, line, ColorBlue, 5.0, 10.0)

		// line = C.Line{
		// 	start: C.Vec2{x: 100.0, y: 100.0},
		// 	end:   C.Vec2{x: 950, y: 200.0},
		// }
		// C.draw_line_dotted(renderer, line, ColorWhite, 2.0)

		text := &Text{
			Component: Component{
				ComponentType: TText,
				Pos:           Position{X: 10, Y: 10, Type: PositionTypeRelative},
				Size:          Vec2{X: 180, Y: 30},
				ID:            "text1",
				Parent:        nil,
				Children:      nil,
			},
			Text: "Hello, World!",
		}

		button := &Button{
			Component: Component{
				ComponentType: TButton,
				Pos:           Position{X: 10, Y: 50, Type: PositionTypeRelative},
				Size:          Vec2{X: 180, Y: 30},
				ID:            "button1",
				Parent:        nil,
				Children:      nil,
			},
			Text:     "Click Me",
			Callback: func() { fmt.Println("Button clicked!") },
			Pressed:  false,
			Released: false,
		}

		container := &Container{
			Component: Component{
				ComponentType: TContainer,
				Pos:           Position{X: 50, Y: 50},
				Size:          Vec2{X: 200, Y: 200},
				ID:            "container1",
				Parent:        nil,
				Children: []IComponent{
					text,
					button,
				},
			},
		}

		text.Component.Parent = container
		button.Component.Parent = container

		// render
		componentRenderer := &ComponentRenderer{Component: container}
		componentRenderer.Render(renderer)

		C.present_screen(renderer)

		// Don't forget to process events
		C.handle_events(renderer)
	}

	//component demo

	// Set parent-child relationships

	fmt.Println("Exiting")
}
