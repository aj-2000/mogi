package ui

import (
	"mogi/color"
	"mogi/math"
)

//
// ——————————————————————————————————————————————————————————————————————————————
// Component
// ——————————————————————————————————————————————————————————————————————————————
//

type Component struct {
	kind            ComponentKind
	pos             Position
	size            math.Vec2f32
	id              string
	fullID          string
	children        []IComponent
	parent          IComponent
	display         Display
	margin          math.Vec2f32
	padding         math.Vec2f32
	border          math.Vec2f32
	borderRadius    float32
	gap             math.Vec2f32
	borderColor     color.RGBA
	backgroundColor color.RGBA
	flexItemProps   FlexItemProps
}

func newComponentBase(kind ComponentKind) Component {
	return Component{
		kind:            kind,
		id:              "",
		children:        make([]IComponent, 0),
		flexItemProps:   NewFlexItemProps(),
		pos:             Position{X: 0, Y: 0, Type: PositionTypeRelative},
		display:         DisplayInline,
		border:          math.Vec2f32{X: 0, Y: 0},
		margin:          math.Vec2f32{X: 0, Y: 0},
		padding:         math.Vec2f32{X: 0, Y: 0},
		borderRadius:    0,
		borderColor:     color.Black,
		backgroundColor: color.Transparent,
	}
}

//
// ——————————————————————————————————————————————————————————————————————————————
// Implementation of IComponent for Component
// ——————————————————————————————————————————————————————————————————————————————
//

// ——————————————————————————————————————————————————————————————————————————————
// Getters
// ——————————————————————————————————————————————————————————————————————————————

func (c *Component) Kind() ComponentKind    { return c.kind }
func (c *Component) Pos() Position          { return c.pos }
func (c *Component) Size() math.Vec2f32     { return c.size }
func (c *Component) FullID() string         { return c.fullID }
func (c *Component) ID() string             { return c.id }
func (c *Component) Children() []IComponent { return c.children }

// should we return a copy?
func (c *Component) FlexItem() *FlexItemProps    { return &c.flexItemProps }
func (c *Component) Margin() math.Vec2f32        { return c.margin }
func (c *Component) Padding() math.Vec2f32       { return c.padding }
func (c *Component) Border() math.Vec2f32        { return c.border }
func (c *Component) Gap() math.Vec2f32           { return c.gap }
func (c *Component) BorderColor() color.RGBA     { return c.borderColor }
func (c *Component) BackgroundColor() color.RGBA { return c.backgroundColor }
func (c *Component) AbsolutePos() math.Vec2f32 {
	if c.parent != nil {
		return math.Vec2f32{
			X: c.pos.X + c.parent.AbsolutePos().X,
			Y: c.pos.Y + c.parent.AbsolutePos().Y,
		}
	}
	return c.pos.Vec2f32()
}
func (c *Component) IsPointInComponent(point math.Vec2f32) bool {
	absPos := c.AbsolutePos()
	return point.X >= absPos.X && point.X <= absPos.X+c.size.X &&
		point.Y >= absPos.Y && point.Y <= absPos.Y+c.size.Y
}
func (c *Component) SetParent(p IComponent) {
	c.parent = p
}
func (c *Component) BorderRadius() float32 { return c.borderRadius }
func (c *Component) Parent() IComponent {
	return c.parent
}
func (c *Component) Display() Display { return c.display }

// ——————————————————————————————————————————————————————————————————————————————
// Internal Setters (used by layout engine)
// ——————————————————————————————————————————————————————————————————————————————
func (c *Component) setPos(p Position) { c.pos = p }
func (c *Component) setDisplay(d Display) {
	c.display = d
}
func (c *Component) setSize(s math.Vec2f32)        { c.size = s }
func (c *Component) setMargin(margin math.Vec2f32) { c.margin = margin }
func (c *Component) setPadding(padding math.Vec2f32) {
	c.padding = padding
}
func (c *Component) setBorder(border math.Vec2f32) { c.border = border }
func (c *Component) setBorderRadius(radius float32) {
	if radius < 0 {
		radius = 0
	}
	c.borderRadius = radius
}
func (c *Component) setGap(gap math.Vec2f32) { c.gap = gap }
func (c *Component) setBorderColor(color color.Color) {
	c.borderColor = color.ToRGBA()
}
func (c *Component) setID(id string) {
	c.id = id
}
func (c *Component) setFullID(fullID string) {
	c.fullID = fullID
}
func (c *Component) setBackgroundColor(color color.Color) {
	c.backgroundColor = color.ToRGBA()
}

// ——————————————————————————————————————————————————————————————————————————————
// Fluent Setters for Flex Item Properties )
// ——————————————————————————————————————————————————————————————————————————————
func (c *Component) SetFlexGrow(grow float32) *Component {
	if grow < 0 {
		grow = 0
	}
	c.flexItemProps.Grow = grow
	return c
}
func (c *Component) SetFlexShrink(shrink float32) *Component {
	if shrink < 0 {
		shrink = 0
	}
	c.flexItemProps.Shrink = shrink
	return c
}
func (c *Component) SetFlexBasis(basis float32) *Component {
	c.flexItemProps.Basis = basis
	return c
}
func (c *Component) SetAlignSelf(align AlignItems) *Component {
	c.flexItemProps.AlignSelf = align
	return c
}
func (c *Component) SetOrder(order int) *Component {
	c.flexItemProps.Order = order
	return c
}
