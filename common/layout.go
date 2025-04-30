package common

import (
	"fmt"
)

// LayoutEngine calculates component positions and sizes based on Flexbox rules.
type LayoutEngine struct {
	CalculateTextWidth func(string, float32) float32
}

// NewLayoutEngine creates a layout engine.
func NewLayoutEngine(f func(string, float32) float32) *LayoutEngine {
	return &LayoutEngine{
		CalculateTextWidth: f,
	}
}

// CalculateLayout computes the layout for the component tree starting at root.
// availableSize is the constraint for the root component.
func (le *LayoutEngine) Layout(root IComponent, startingPos Vec2, availableSize Vec2) {
	le.populateSize(root, availableSize)
	le.populatePosition(root, startingPos)
	// le.printComponentTree(root, "  ")
	// panic("LayoutEngine.Layout: Layout calculation complete")
}

func (le *LayoutEngine) printComponentTree(root IComponent, indent string) {
	fmt.Printf("%s[%s: Size:%v, Position:%v]\n", indent, root.ID(), root.Size(), root.Pos())
	children := root.Children()
	if len(children) > 0 {
		fmt.Printf("%s  Children:\n", indent)
		for _, child := range children {
			le.printComponentTree(child, indent+"    ")
		}
	}
}

func (le *LayoutEngine) populatePosition(root IComponent, startingPos Vec2) {
	// log.Printf(
	// 	"populatePosition: ID=%s, startingPos=%v, currentPos=%v, size=%v, parentPos=%v, parentSize=%v\n",
	// 	root.ID(),
	// 	startingPos,
	// 	root.Pos(),
	// 	root.Size(),
	// 	func() *Vec2 {
	// 		if root.Parent() != nil {
	// 			return &Vec2{root.Parent().Pos().X, root.Parent().Pos().Y}
	// 		}
	// 		return nil
	// 	}(),
	// 	func() *Vec2 {
	// 		if root.Parent() != nil {
	// 			return &Vec2{root.Parent().Size().X, root.Parent().Size().Y}
	// 		}
	// 		return nil
	// 	}(),
	// )
	if root.Pos().Type != PositionTypeAbsolute {
		var parentPos Vec2
		if root.Parent() != nil {
			parentPos = Vec2{root.Parent().Pos().X, root.Parent().Pos().Y}
		}
		pos := Position{
			Type: PositionTypeRelative,
			X:    startingPos.X - parentPos.X,
			Y:    startingPos.Y - parentPos.Y,
		}
		root.setPos(pos)
	}

	switch comp := root.(type) {
	case *Container:
		// Calculate the position of the container's children
		currPos := startingPos
		longestChildHeight := float32(0.0)
		for _, child := range comp.Children() {
			if child.Pos().Type == PositionTypeRelative {
				longestChildHeight = max(longestChildHeight, child.Size().Y)

				// Check if child fits in current row
				if (currPos.X+child.Size().X > comp.Size().X+startingPos.X) || child.Display() == DisplayBlock {
					// Move to next line
					currPos.X = startingPos.X
					currPos.Y += longestChildHeight
					longestChildHeight = 0.0
				}

				le.populatePosition(child, currPos)

				// Move X for next child
				currPos.X += child.Size().X
			} else {
				le.populatePosition(child, Vec2{child.Pos().X, child.Pos().Y})
			}
		}

	case *Text:
	case *Button:
	default:
		panic("Unsupported component type")
	}
}

func (le *LayoutEngine) populateSize(root IComponent, availableSize Vec2) {

	var getSize func(IComponent, Vec2) Vec2
	getSize = func(root IComponent, availableSize Vec2) Vec2 {
		switch comp := root.(type) {
		case *Container:
			// Calculate the size of the container based on its children
			var currWidth float32 = 0.0
			var totalHeightTillPreviousRow float32 = 0.0
			var currRowMaxHeight float32 = 0.0
			var maxWidth float32 = 0.0
			if comp.Size().X != 0 && comp.Size().Y != 0 {
				availableSize = comp.Size()
			}
			for i, child := range comp.Children() {
				childSize := getSize(child, availableSize)
				currRowMaxHeight = max(childSize.Y, currRowMaxHeight)

				if i == 0 {
					currWidth = childSize.X
					maxWidth = max(currWidth, childSize.X)
				} else {
					if currWidth+childSize.X > availableSize.X {
						// Move to the next line
						currWidth = 0
						totalHeightTillPreviousRow += currRowMaxHeight
						currRowMaxHeight = childSize.Y
					}
					currWidth += childSize.X
					maxWidth = max(maxWidth, currWidth)
				}

			}
			// Set the size of the container based on its children
			if comp.Size().X != 0 && comp.Size().Y != 0 {
				return comp.Size()
			} else {
				containerSize := Vec2{
					X: maxWidth,
					Y: totalHeightTillPreviousRow + currRowMaxHeight}
				comp.SetSize(containerSize)
				return containerSize
			}
		case *Text:
			width := le.CalculateTextWidth(comp.Content, 24.0) + 30
			height := 24.0
			comp.SetSize(Vec2{X: width, Y: float32(height)})
			return Vec2{X: width, Y: float32(height)}
		case *Button:
			width := le.CalculateTextWidth(comp.Label, 24.0)
			height := 24.0
			comp.SetSize(Vec2{X: width, Y: float32(height)})
			return Vec2{X: width, Y: float32(height)}
		default:
			panic("Unsupported component type")
		}
	}

	size := getSize(root, availableSize)
	root.setSize(size)
}

// Helper needed in Container to access props
func (c *Container) FlexContainerProps() FlexContainerProps {
	return c.flexContainerProps
}
