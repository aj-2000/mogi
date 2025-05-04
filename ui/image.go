package ui

import "mogi/math"

// ——————————————————————————————————————————————————————————————————————————————
// Image
// ——————————————————————————————————————————————————————————————————————————————

type Image struct {
	Component
	Path string
}

func NewImage(path string) *Image {
	i := &Image{
		Component: newComponentBase(ImageKind),
		Path:      path,
	}
	return i
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters
// ——————————————————————————————————————————————————————————————————————————————

func (i *Image) SetID(id string) *Image {
	i.Component.setID(id)
	return i
}

func (i *Image) SetDisplay(d Display) *Image {
	i.Component.setDisplay(d)
	return i
}

func (i *Image) SetPosition(pos Position) *Image {
	i.Component.setPos(pos)
	return i
}

func (i *Image) SetSize(size math.Vec2f32) *Image {
	i.Component.setSize(size)
	return i
}
