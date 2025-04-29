package common

import "fmt"

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
	le.printComponentTree(root, "  ")
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
		for _, child := range comp.Children() {
			if child.Pos().Type == PositionTypeRelative {
				// Check if child fits in current row
				if currPos.X+child.Size().X > comp.Size().X+startingPos.X {
					// Move to next line
					currPos.X = startingPos.X
					currPos.Y += child.Size().Y
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
			if comp.Size().X != 0 && comp.Size().Y != 0 {
				return comp.Size()
			}
			// Calculate the size of the container based on its children
			var totalWidth float32 = 0.0
			var totalHeight float32 = 0.0
			var maxWidth float32 = 0.0
			for i, child := range comp.Children() {
				childSize := getSize(child, availableSize)

				if i == 0 {
					totalHeight = childSize.Y
					totalWidth = childSize.X
					maxWidth = max(totalWidth, childSize.X)
				} else {
					// Check if the child fits in the available width
					if totalWidth+childSize.X > availableSize.X {
						// Move to the next line
						totalWidth = 0
						totalHeight += childSize.Y
					}
					totalWidth += childSize.X
					maxWidth = max(maxWidth, childSize.X)
				}

			}
			// Set the size of the container based on its children
			containerSize := Vec2{
				X: maxWidth,
				Y: totalHeight}
			comp.SetSize(containerSize)

			return containerSize
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
