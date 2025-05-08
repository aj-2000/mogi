package ui

import (
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

// ——————————————————————————————————————————————————————————————————————————————
// Button Component
// ——————————————————————————————————————————————————————————————————————————————

type Button struct {
	Component
	Label        string
	Callback     func(self *Button)
	HoverColor   color.RGBA
	PressedColor color.RGBA
	TextColor    color.RGBA
	IsPressed    bool
	IsMouseOver  bool
}

func NewButton(label string) *Button {
	b := &Button{
		Component:    newComponentBase(ButtonKind),
		Label:        label,
		Callback:     nil,
		HoverColor:   color.RGBA{R: 0.3, G: 0.5, B: 0.9, A: 1},
		PressedColor: color.RGBA{R: 0.1, G: 0.3, B: 0.7, A: 1},
		TextColor:    color.RGBA{R: 1, G: 1, B: 1, A: 1},
	}
	b.Component.setDisplay(DisplayBlock)
	return b
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters
// ——————————————————————————————————————————————————————————————————————————————

func (b *Button) SetID(id string) *Button {
	b.Component.setID(id)
	return b
}

func (b *Button) SetDisplay(d Display) *Button {
	b.Component.setDisplay(d)
	return b
}

func (b *Button) SetLabel(label string) *Button {
	b.Label = label
	return b
}

func (b *Button) SetOnClick(callback func(self *Button)) *Button {
	b.Callback = callback
	return b
}

func (b *Button) SetBackgroundColor(color color.RGBA) *Button {
	b.Component.setBackgroundColor(color)
	return b
}

func (b *Button) SetHoverColor(color color.RGBA) *Button {
	b.HoverColor = color
	return b
}

func (b *Button) SetPressedColor(color color.RGBA) *Button {
	b.PressedColor = color
	return b
}

func (b *Button) SetTextColor(color color.RGBA) *Button {
	b.TextColor = color
	return b
}

func (b *Button) SetPosition(pos Position) *Button {
	b.Component.setPos(pos)
	return b
}

func (b *Button) SetSize(size math.Vec2f32) *Button {
	b.Component.setSize(size)
	return b
}

func (b *Button) SetZIndex(zIndex int) *Button {
	b.Component.setZIndex(zIndex)
	return b
}

func (b *Button) FontSize() float32 {
	return 24.0
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters for Flex Item Properties )
// ——————————————————————————————————————————————————————————————————————————————

func (b *Button) SetFlexGrow(grow float32) *Button {
	b.Component.SetFlexGrow(grow)
	return b
}

func (b *Button) SetFlexShrink(shrink float32) *Button {
	b.Component.SetFlexShrink(shrink)
	return b
}

func (b *Button) SetFlexBasis(basis float32) *Button {
	b.Component.SetFlexBasis(basis)
	return b
}

func (b *Button) SetFlexBasisAuto() *Button {
	b.Component.SetFlexBasis(FlexBasisAuto)
	return b
}

func (b *Button) SetAlignSelf(align AlignItems) *Button {
	b.Component.SetAlignSelf(align)
	return b
}

func (b *Button) SetOrder(order int) *Button {
	b.Component.SetOrder(order)
	return b
}
