package ui

import (
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

// ——————————————————————————————————————————————————————————————————————————————
// Text Component
// ——————————————————————————————————————————————————————————————————————————————

type Text struct {
	Component
	Content  string
	Color    color.RGBA
	FontSize float32
}

func NewText(content string) *Text {
	t := &Text{
		Component: newComponentBase(TextKind),
		Content:   content,
		Color:     color.Black, // Default black
		FontSize:  16.0,        // Default font size
	}
	return t
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters
// ——————————————————————————————————————————————————————————————————————————————

func (t *Text) SetID(id string) *Text {
	t.Component.setID(id)
	return t
}

func (t *Text) SetContent(content string) *Text {
	t.Content = content
	return t
}

func (t *Text) SetColor(color color.RGBA) *Text {
	t.Color = color
	return t
}

func (t *Text) SetFontSize(size float32) *Text {
	if size <= 0 {
		size = 16.0 // Reset to default if invalid
	}
	t.FontSize = size
	return t
}

func (t *Text) SetDisplay(d Display) *Text {
	t.Component.setDisplay(d)
	return t
}

func (t *Text) SetPosition(pos Position) *Text {
	t.Component.setPos(pos)
	return t
}

func (t *Text) SetSize(size math.Vec2f32) *Text {
	t.Component.setSize(size)
	return t
}

func (t *Text) SetFlexGrow(grow float32) *Text {
	t.Component.SetFlexGrow(grow)
	return t
}

func (t *Text) SetFlexShrink(shrink float32) *Text {
	t.Component.SetFlexShrink(shrink)
	return t
}

func (t *Text) SetFlexBasis(basis float32) *Text {
	t.Component.SetFlexBasis(basis)
	return t
}

func (t *Text) SetFlexBasisAuto() *Text {
	t.Component.SetFlexBasis(FlexBasisAuto)
	return t
}

func (t *Text) SetAlignSelf(align AlignItems) *Text {
	t.Component.SetAlignSelf(align)
	return t
}

func (t *Text) SetOrder(order int) *Text {
	t.Component.SetOrder(order)
	return t
}

func (t *Text) SetZIndex(zIndex int) *Text {
	t.Component.setZIndex(zIndex)
	return t
}
