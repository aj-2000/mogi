package common

import "log"

// Added for default callback

type ColorRGBA struct {
	R, G, B, A float32
}

type PositionType int

const (
	PositionTypeAbsolute PositionType = iota // May be less used with Flexbox
	PositionTypeRelative                     // Flexbox handles relative positioning
)

type Display int

const (
	DisplayBlock Display = iota
	DisplayInline
	DisplayFlex
	DisplayGrid
	DisplayNone
)

type Position struct {
	X, Y float32
	Type PositionType
}

type Vec2 struct {
	X, Y float32
}

type ComponentKind int

const (
	ContainerKind ComponentKind = iota
	TextKind
	ButtonKind
)

// --- Flexbox Enums and Structs (from previous example) ---

type FlexDirection int

const (
	FlexDirectionRow FlexDirection = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

type FlexWrap int

const (
	FlexWrapNoWrap FlexWrap = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

type JustifyContent int

const (
	JustifyContentFlexStart JustifyContent = iota
	JustifyContentFlexEnd
	JustifyContentCenter
	JustifyContentSpaceBetween
	JustifyContentSpaceAround
	JustifyContentSpaceEvenly
)

type AlignItems int // Also used for AlignSelf, AlignContent

const (
	AlignItemsStretch AlignItems = iota
	AlignItemsFlexStart
	AlignItemsFlexEnd
	AlignItemsCenter
	AlignItemsBaseline
	AlignSelfAuto AlignItems = -1 // Sentinel value for AlignSelf default
)

const (
	FlexBasisAuto float32 = -1.0
)

// --- Flex Item Properties ---
type FlexItemProps struct {
	Grow      float32
	Shrink    float32
	Basis     float32
	AlignSelf AlignItems
	Order     int
	// --- Calculated values by layout engine ---
	// These store the intermediate results during calculation
	// You might not need to expose these publicly.
	// calculatedFlexBasis float32 // The basis used after considering content size
	// mainSize            float32 // Final size along the main axis
	// crossSize           float32 // Final size along the cross axis
}

func NewFlexItemProps() FlexItemProps {
	return FlexItemProps{
		Grow:      0,
		Shrink:    1, // Default shrink is 1
		Basis:     FlexBasisAuto,
		AlignSelf: AlignSelfAuto,
		Order:     0,
	}
}

// --- Flex Container Properties ---
type FlexContainerProps struct {
	Enabled      bool // Default: false
	Direction    FlexDirection
	Wrap         FlexWrap
	Justify      JustifyContent
	AlignItems   AlignItems
	AlignContent AlignItems
	Gap          float32
}

func NewFlexContainerProps() FlexContainerProps {
	return FlexContainerProps{
		Enabled:      false,
		Direction:    FlexDirectionRow,
		Wrap:         FlexWrapNoWrap,
		Justify:      JustifyContentFlexStart,
		AlignItems:   AlignItemsStretch, // Default is stretch
		AlignContent: AlignItemsStretch, // Default is stretch (or flex-start if not wrapping)
		Gap:          0,
	}
}

// --- Core Interface ---

type IComponent interface {
	Kind() ComponentKind
	Pos() Position // Position calculated by layout engine
	Size() Vec2    // Size calculated by layout engine
	ID() string
	SetParent(p IComponent)
	Parent() IComponent // Optional getter for parent component
	Children() []IComponent
	FlexItem() *FlexItemProps // Accessor for flex item properties
	Display() Display

	// --- Internal Setters (used by layout engine) ---
	// These need to be part of the interface if the layout engine
	// works purely on IComponent. Alternatively, the layout engine
	// could use type assertions, but this is cleaner.
	setPos(Position)
	setSize(Vec2)
	setDisplay(Display)
	// Optional: Method to get intrinsic size (needed for flex-basis: auto)
	// CalculateIntrinsicSize(available Vec2) Vec2
}

// --- Base Struct ---

type Component struct {
	kind     ComponentKind
	pos      Position // Calculated by layout engine
	size     Vec2     // Calculated by layout engine
	id       string
	children []IComponent
	parent   IComponent // Optional but useful for layout/events
	display  Display
	// Flex properties for when THIS component is an ITEM in a flex container
	flexItemProps FlexItemProps
}

// --- Base Constructor (Internal) ---
// Creates the base component with defaults. Specific types call this.
func newComponentBase(kind ComponentKind) Component {
	return Component{
		kind:          kind,
		id:            "",
		children:      make([]IComponent, 0), // Initialize slice
		flexItemProps: NewFlexItemProps(),
		display:       DisplayInline,
		// pos and size start zeroed, layout engine calculates them
	}
}

// --- Implement IComponent ---

func (c *Component) Kind() ComponentKind      { return c.kind }
func (c *Component) Pos() Position            { return c.pos }
func (c *Component) Size() Vec2               { return c.size }
func (c *Component) ID() string               { return c.id }
func (c *Component) Children() []IComponent   { return c.children }
func (c *Component) FlexItem() *FlexItemProps { return &c.flexItemProps }

func (c *Component) SetParent(p IComponent) {
	c.parent = p
}

// Optional Getter
func (c *Component) Parent() IComponent {
	return c.parent
}

// --- Internal Setters ---
func (c *Component) setPos(p Position) { c.pos = p }
func (c *Component) setDisplay(d Display) {
	c.display = d
}
func (c *Component) Display() Display { return c.display }
func (c *Component) setSize(s Vec2)   { c.size = s }

// --- Fluent Setters for Flex Item Properties (on Base Component) ---
// These return *Component so they can be chained from any component type

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

func (c *Component) SetID(id string) *Component {
	c.id = id
	return c
}

// --- Components ---

type Container struct {
	Component
	BackgroundColor ColorRGBA
	BorderColor     ColorRGBA
	BorderWidth     float32
	BorderRadius    float32

	// Flex properties for when THIS component IS a flex container
	flexContainerProps FlexContainerProps
}

// --- Container Constructor ---
func NewContainer() *Container {
	return &Container{
		Component:          newComponentBase(ContainerKind),
		flexContainerProps: NewFlexContainerProps(),
		BackgroundColor:    ColorRGBA{0, 0, 0, 0}, // Default transparent
		BorderColor:        ColorRGBA{0, 0, 0, 0}, // Default transparent
	}
}

// --- Fluent Setters for Container Visual Properties ---

func (c *Container) SetID(id string) *Container {
	c.Component.SetID(id)
	return c
}

func (c *Container) SetBackgroundColor(color ColorRGBA) *Container {
	c.BackgroundColor = color
	return c
}

func (c *Container) SetBorderColor(color ColorRGBA) *Container {
	c.BorderColor = color
	return c
}

func (c *Container) SetBorderWidth(width float32) *Container {
	if width < 0 {
		width = 0
	}
	c.BorderWidth = width
	return c
}

func (c *Container) SetBorderRadius(radius float32) *Container {
	if radius < 0 {
		radius = 0
	}
	c.BorderRadius = radius
	return c
}

// --- Fluent Setters for Container Position/Size (Manual - Layout Engine Overrides) ---

func (c *Container) SetPosition(pos Position) *Container {
	c.Component.setPos(pos)
	return c
}

func (c *Container) SetSize(size Vec2) *Container {
	c.Component.setSize(size)
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

// --- Fluent Setters for Flex Container Properties ---

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

func (c *Container) SetGap(gap float32) *Container {
	if gap < 0 {
		gap = 0
	}
	c.flexContainerProps.Gap = gap
	return c
}

// --- Add Child Method ---
// In common/container.go
func (c *Container) AddChild(child IComponent) *Container {
	if child != nil {
		c.children = append(c.children, child)
		child.SetParent(c) // Simply call the interface method!
	} else {
		log.Fatalf("Child component cannot be nil")
	}
	return c
}

// AddChildren adds multiple children components.
func (c *Container) AddChildren(children ...IComponent) *Container {
	for _, child := range children {
		c.AddChild(child) // Reuse existing AddChild logic (parent setting etc)
	}
	return c
}

// --- Text Component ---

type Text struct {
	Component
	Content  string
	Color    ColorRGBA
	FontSize float32
	// FontFace string // Consider adding FontFace property
}

// --- Text Constructor ---
func NewText(content string) *Text {
	return &Text{
		Component: newComponentBase(TextKind),
		Content:   content,
		Color:     ColorRGBA{0, 0, 0, 1}, // Default black
		FontSize:  16.0,                  // Default font size
	}
}

// --- Fluent Setters for Text ---

func (t *Text) SetID(id string) *Text {
	t.Component.SetID(id)
	return t
}

func (t *Text) SetContent(content string) *Text {
	t.Content = content
	return t
}

func (t *Text) SetColor(color ColorRGBA) *Text {
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

// --- Fluent Setters for Text Position/Size (Manual - Layout Engine Overrides) ---

func (t *Text) SetPosition(pos Position) *Text {
	t.Component.setPos(pos)
	return t
}

func (t *Text) SetSize(size Vec2) *Text {
	t.Component.setSize(size)
	return t
}

// --- Fluent Setters for Text Flex Item Properties ---
// These allow a Text component to act as a flex item

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

// --- Button Component ---

type Button struct {
	Component
	Label     string
	Callback  func()
	Pressed   bool // State: currently pressed down
	Released  bool // State: just released (useful for click detection)
	MouseOver bool // State: mouse is hovering
	// Add visual properties like background color, border, etc.
	BackgroundColor ColorRGBA
	HoverColor      ColorRGBA
	PressedColor    ColorRGBA
	TextColor       ColorRGBA
}

// --- Button Constructor ---
func NewButton(label string) *Button {
	b := &Button{
		Component: newComponentBase(ButtonKind),
		Label:     label,
		Callback:  func() {},
		// Default visual styles
		BackgroundColor: ColorRGBA{0.2, 0.4, 0.8, 1}, // Blueish
		HoverColor:      ColorRGBA{0.3, 0.5, 0.9, 1},
		PressedColor:    ColorRGBA{0.1, 0.3, 0.7, 1},
		TextColor:       ColorRGBA{1, 1, 1, 1}, // White
	}

	b.Component.setDisplay(DisplayBlock)
	return b
}

// --- Fluent Setters for Button ---

func (b *Button) SetID(id string) *Button {
	b.Component.SetID(id)
	return b
}

func (b *Button) SetLabel(label string) *Button {
	b.Label = label
	return b
}

func (b *Button) SetOnClick(callback func()) *Button {
	if callback != nil {
		b.Callback = callback
	} else {
		// Provide a no-op callback if nil is passed, or keep default
		b.Callback = func() {}
	}
	return b
}

func (b *Button) SetBackgroundColor(color ColorRGBA) *Button {
	b.BackgroundColor = color
	return b
}

func (b *Button) SetHoverColor(color ColorRGBA) *Button {
	b.HoverColor = color
	return b
}

func (b *Button) SetPressedColor(color ColorRGBA) *Button {
	b.PressedColor = color
	return b
}

func (b *Button) SetTextColor(color ColorRGBA) *Button {
	b.TextColor = color
	return b
}

// --- Fluent Setters for Button Position/Size (Manual - Layout Engine Overrides) ---

func (b *Button) SetPosition(pos Position) *Button {
	b.Component.setPos(pos)
	return b
}

func (b *Button) SetSize(size Vec2) *Button {
	b.Component.setSize(size)
	return b
}

// --- Fluent Setters for Button Flex Item Properties ---
// These allow a Button component to act as a flex item

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
