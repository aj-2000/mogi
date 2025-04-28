package common

type ColorRGBA struct {
	R, G, B, A float32
}

// --- Position Types ---

type PositionType int

const (
	PositionTypeAbsolute PositionType = iota
	PositionTypeRelative
)

type Position struct {
	X, Y float32
	Type PositionType
}

// --- Core Interface ---

type IComponent interface {
	Kind() ComponentKind
	Pos() Position
	Size() Vec2
	ID() string
	Children() []IComponent
}

// --- Base Structs ---

type Vec2 struct {
	X, Y float32
}

type ComponentKind int

const (
	ContainerKind ComponentKind = iota
	TextKind
	ButtonKind
)

type Component struct {
	kind     ComponentKind
	pos      Position
	size     Vec2
	id       string
	children []IComponent
}

// --- Implement IComponent automatically ---

func (c *Component) Kind() ComponentKind {
	return c.kind
}

func (c *Component) Pos() Position {
	return c.pos
}

func (c *Component) Size() Vec2 {
	return c.size
}

func (c *Component) ID() string {
	return c.id
}

func (c *Component) Children() []IComponent {
	return c.children
}

// --- Components ---

type Container struct {
	Component
	BackgroundColor ColorRGBA
	BorderColor     ColorRGBA
	BorderWidth     float32
	BorderRadius    float32
}

type Text struct {
	Component
	Content  string
	Color    ColorRGBA
	FontSize float32
}

type Button struct {
	Component
	Label     string
	Callback  func()
	Pressed   bool
	Released  bool
	MouseOver bool
}

// --- Constructors ---

type ContainerOptions struct {
	BackgroundColor ColorRGBA
	BorderColor     ColorRGBA
	BorderWidth     float32
	BorderRadius    float32
	Position        Position
	ID              string
	Children        []IComponent
	Size            Vec2 // Optional, can be zero for default
}

func NewContainer(opts ContainerOptions) *Container {
	return &Container{
		Component: Component{
			kind:     ContainerKind,
			id:       opts.ID,
			pos:      opts.Position,
			size:     opts.Size,
			children: opts.Children,
		},
		BackgroundColor: opts.BackgroundColor,
		BorderColor:     opts.BorderColor, // Default black
		BorderWidth:     opts.BorderWidth,
		BorderRadius:    opts.BorderRadius,
	}
}

type TextOptions struct {
	ID       string
	Content  string
	FontSize float32
	Color    ColorRGBA
	Position Position
	Children []IComponent
	Size     Vec2 // Optional, can be zero for default
}

func NewText(opts TextOptions) *Text {
	size := opts.Size
	if size == (Vec2{}) {
		size = Vec2{X: 100, Y: 30}
	}
	return &Text{
		Component: Component{
			kind:     TextKind,
			id:       opts.ID,
			pos:      opts.Position,
			size:     size,
			children: opts.Children,
		},
		Content:  opts.Content,
		Color:    opts.Color,
		FontSize: opts.FontSize,
	}
}

type ButtonOptions struct {
	ID       string
	Label    string
	Callback func()
	Position Position
	Children []IComponent
	Size     Vec2 // Optional, can be zero for default
}

func NewButton(opts ButtonOptions) *Button {
	size := opts.Size
	if size == (Vec2{}) {
		size = Vec2{X: 120, Y: 40}
	}
	return &Button{
		Component: Component{
			kind:     ButtonKind,
			id:       opts.ID,
			pos:      opts.Position,
			size:     size,
			children: opts.Children,
		},
		Label:    opts.Label,
		Callback: opts.Callback,
	}
}
