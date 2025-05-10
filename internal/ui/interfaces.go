package ui

import (
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

type IComponent interface {
	Kind() ComponentKind
	Pos() Position
	Size() math.Vec2f32
	ID() string
	FullID() string

	BorderColor() color.RGBA
	BackgroundColor() color.RGBA
	Margin() math.Vec2f32
	Padding() math.Vec2f32
	Border() math.Vec2f32
	Gap() math.Vec2f32
	BorderRadius() float32
	SetParent(p IComponent)

	AbsolutePos() math.Vec2f32
	AbsoluteZIndex() int
	Parent() IComponent
	Children() []IComponent
	FlexItem() *FlexItemProps
	Display() Display
	IsPointInsideComponent(point math.Vec2f32) bool
	ZIndex() int
	WidthPercent() float32
	HeightPercent() float32

	// --- Fluent Setters ---

	// --- Internal Setters (used by layout engine) ---
	// These need to be part of the interface if the layout engine
	// works purely on IComponent. Alternatively, the layout engine
	// could use type assertions, but this is cleaner.
	setPos(Position)
	setSize(math.Vec2f32)
	setGap(gap math.Vec2f32)
	setMargin(margin math.Vec2f32)
	setBorderRadius(radius float32)
	setDisplay(Display)
	setBorderColor(color color.Color)
	setBackgroundColor(color color.Color)
	setBorder(border math.Vec2f32)
	setPadding(padding math.Vec2f32)
	setID(id string)
	setFullID(fullID string)
	setZIndex(zIndex int)
	setWidthPercent(widthPercent float32)
	setHeightPercent(heightPercent float32)
	// Optional: Method to get intrinsic size (needed for flex-basis: auto)
	// CalculateIntrinsicSize(available math.Vec2f32 ) math.Vec2f32
}
