package common

type ColorRGBA struct {
	R, G, B, A float32
}

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
	Children() []IComponent
}

type Component struct {
	ComponentType ComponentType
	Pos           Position
	Size          Vec2
	ID            string
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
	BackgroundColor ColorRGBA
	BorderColor     ColorRGBA
	BorderWidth     float32
	BorderRadius    float32
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

func (b *Button) Children() []IComponent {
	return b.Component.Children
}
