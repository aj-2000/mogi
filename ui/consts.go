package ui

import "github.com/aj-2000/mogi/math"

//
// ——————————————————————————————————————————————————————————————————————————————
// Display
// ——————————————————————————————————————————————————————————————————————————————
//

type Display int

const (
	DisplayBlock Display = iota
	DisplayInline
	DisplayFlex
	DisplayGrid
	DisplayNone
)

//
// ——————————————————————————————————————————————————————————————————————————————
// Position
// ——————————————————————————————————————————————————————————————————————————————
//

type PositionType int

const (
	PositionTypeAbsolute PositionType = iota
	PositionTypeRelative
)

type Position struct {
	X, Y float32
	Type PositionType
}

func (p Position) Vec2f32() math.Vec2f32 {
	return math.Vec2f32{X: p.X, Y: p.Y}
}

//
// ——————————————————————————————————————————————————————————————————————————————
// Component
// ——————————————————————————————————————————————————————————————————————————————
//

type ComponentKind int

const (
	ContainerKind ComponentKind = iota
	TextKind
	ButtonKind
	ImageKind
)

func (k ComponentKind) String() string {
	switch k {
	case ContainerKind:
		return "Container"
	case TextKind:
		return "Text"
	case ButtonKind:
		return "Button"
	case ImageKind:
		return "Image"
	default:
		return "Unknown"
	}
}

//
// ——————————————————————————————————————————————————————————————————————————————
// Flex
// ——————————————————————————————————————————————————————————————————————————————
//

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
		Shrink:    1,
		Basis:     FlexBasisAuto,
		AlignSelf: AlignSelfAuto,
		Order:     0,
	}
}

type FlexContainerProps struct {
	Enabled      bool
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
		AlignItems:   AlignItemsStretch,
		AlignContent: AlignItemsStretch,
		Gap:          0,
	}
}
