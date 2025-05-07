package ui

import (
	"log"

	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

type Container struct {
	Component
	// Flex properties for when THIS component IS a flex container
	flexContainerProps FlexContainerProps
}

// --- Container Constructor ---
func NewContainer() *Container {
	c := &Container{
		Component:          newComponentBase(ContainerKind),
		flexContainerProps: NewFlexContainerProps(),
	}
	// c.SetMargin(math.Vec2f32 {X: 3, Y: 3})                       // Default margin
	// c.SetPadding(math.Vec2f32 {X: 4, Y: 4})                      // Default padding
	// c.SetBorder(math.Vec2f32 {X: 2, Y: 2})                       // Default border
	// c.SetBorderColor(color.RGBA{R: 1, G: 1, B: 1, A: 1}) // white
	// c.SetBorderRadius(10)                               // Default border radius
	// c.SetGap(math.Vec2f32 {X: 5, Y: 5})                          // Default gap for flex items
	return c
}

// --- Fluent Setters for Container Visual Properties ---

func (c *Container) SetID(id string) *Container {
	c.Component.setID(id)
	return c
}

func (c *Container) SetBackgroundColor(color color.RGBA) *Container {
	c.Component.setBackgroundColor(color)
	return c
}

func (c *Container) SetPosition(pos Position) *Container {
	c.Component.setPos(pos)
	return c
}

func (c *Container) SetSize(size math.Vec2f32) *Container {
	c.Component.setSize(size)
	return c
}

func (c *Container) SetMargin(margin math.Vec2f32) *Container {
	c.Component.setMargin(margin)
	return c
}

func (c *Container) SetPadding(padding math.Vec2f32) *Container {
	c.Component.setPadding(padding)
	return c
}

func (c *Container) SetBorder(border math.Vec2f32) *Container {
	c.Component.setBorder(border)
	return c
}

func (c *Container) SetBorderRadius(radius float32) *Container {
	c.Component.setBorderRadius(radius)
	return c
}

func (c *Container) SetBorderColor(color color.Color) *Container {
	c.Component.setBorderColor(color)
	return c
}

func (c *Container) SetGap(gap math.Vec2f32) *Container {
	c.Component.setGap(gap)
	return c
}

// --- Fluent Setters for Container Flex Item Properties ---
// These allow a Container to act as a flex item within another container

func (c *Container) SetFlexGrow(grow float32) *Container {
	c.Component.SetFlexGrow(grow)
	return c
}

func (c *Container) SetFlexShrink(shrink float32) *Container {
	c.Component.SetFlexShrink(shrink)
	return c
}

func (c *Container) SetFlexBasis(basis float32) *Container {
	c.Component.SetFlexBasis(basis)
	return c
}

func (c *Container) SetFlexBasisAuto() *Container {
	c.Component.SetFlexBasis(FlexBasisAuto)
	return c
}

func (c *Container) SetFlexAlignSelf(align AlignItems) *Container {
	c.Component.SetAlignSelf(align)
	return c
}

func (c *Container) SetFlexOrder(order int) *Container {
	c.Component.SetOrder(order)
	return c
}

func (c *Container) SetDisplay(d Display) *Container {
	c.Component.setDisplay(d)
	return c
}

// ——————————————————————————————————————————————————————————————————————————————
// AddChild and AddChildren methods for adding child components
// ——————————————————————————————————————————————————————————————————————————————

func (c *Container) AddChild(child IComponent) *Container {
	if child != nil {
		c.children = append(c.children, child)
		child.SetParent(c)
	} else {
		log.Fatalf("Child component cannot be nil")
	}
	return c
}

func (c *Container) AddChildren(children ...IComponent) *Container {
	for _, child := range children {
		c.AddChild(child)
	}
	return c
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters for Flex Item Properties )
// ——————————————————————————————————————————————————————————————————————————————

func (c *Container) SetFlexEnabled(enabled bool) *Container {
	c.flexContainerProps.Enabled = enabled
	return c
}

func (c *Container) SetFlexDirection(direction FlexDirection) *Container {
	c.flexContainerProps.Direction = direction
	return c
}

func (c *Container) SetFlexWrap(wrap FlexWrap) *Container {
	c.flexContainerProps.Wrap = wrap
	return c
}

func (c *Container) SetJustifyContent(justify JustifyContent) *Container {
	c.flexContainerProps.Justify = justify
	return c
}

func (c *Container) SetAlignItems(align AlignItems) *Container {
	c.flexContainerProps.AlignItems = align
	return c
}

func (c *Container) SetAlignContent(align AlignItems) *Container {
	c.flexContainerProps.AlignContent = align
	return c
}
