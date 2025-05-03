package common

import (
	"fmt"
)

// Assume Vec2, Position, PositionTypeAbsolute, PositionTypeRelative,
// DisplayBlock, IComponent, Container, Text, Button, and their methods
// (ID, Size, Pos, Children, Parent, setPos, SetSize, setSize, Display,
// Content, Label) are defined elsewhere.

// Assume a max helper function exists:
// func max(a, b float32) float32 { ... }

// LayoutEngine calculates component positions and sizes based on layout rules.
type LayoutEngine struct {
	// CalculateTextWidth estimates the width of a string given a font size.
	CalculateTextWidth func(text string, fontSize float32) float32
}

// NewLayoutEngine creates a layout engine.
func NewLayoutEngine(f func(string, float32) float32) *LayoutEngine {
	return &LayoutEngine{
		CalculateTextWidth: f,
	}
}

// Layout computes the layout for the component tree starting at root.
// availableSize provides the initial constraints for the root component.
func (le *LayoutEngine) Layout(root IComponent, origin Vec2, availableSize Vec2) {
	// Pass 1: Calculate intrinsic sizes (bottom-up)
	le.calculateSizeRecursive(root, availableSize)

	// Pass 2: Calculate positions (top-down)
	// Root component starts at (0, 0) relative to its own frame.
	le.calculatePositionRecursive(root, origin)

	// Optional: Print the final tree for debugging
	// le.printComponentTree(root, "	")
}

// printComponentTree is a helper for debugging the layout structure.
func (le *LayoutEngine) printComponentTree(comp IComponent, indent string) {
	fmt.Printf("%s[%s: Size:%.1f,%.1f Pos:%.1f,%.1f (%v)]\n",
		indent, comp.ID(), comp.Size().X, comp.Size().Y,
		comp.Pos().X, comp.Pos().Y, comp.Pos().Type)
	children := comp.Children()
	if len(children) > 0 {
		// fmt.Printf("%s  Children:\n", indent)
		for _, child := range children {
			le.printComponentTree(child, indent+"	")
		}
	}
}

// calculateSizeRecursive determines the size of each component, starting from
// the leaves and moving up. It respects fixed sizes and calculates content-based
// sizes otherwise. availableSize provides the width constraint for wrapping.
func (le *LayoutEngine) calculateSizeRecursive(comp IComponent, availableSize Vec2) Vec2 {

	// Use component's fixed size if provided
	fixedSize := comp.Size() // Assuming Size() returns the pre-set size
	hasFixedWidth := fixedSize.X > 0
	hasFixedHeight := fixedSize.Y > 0

	var calculatedSize Vec2

	switch c := comp.(type) {
	case *Container:
		// If container has fixed dimensions, use them as constraints for children
		childAvailableSize := availableSize
		if hasFixedWidth {
			childAvailableSize.X = fixedSize.X - comp.Padding().X - comp.Border().X
		}
		if hasFixedHeight {
			// Height constraint isn't typically used for width calculation
			// but could be relevant for scrollable content later.
			childAvailableSize.Y = fixedSize.Y - comp.Padding().Y - comp.Border().Y
		}

		// Calculate children sizes first
		var childrenSizes []Vec2
		for _, child := range c.Children() {
			childrenSizes = append(childrenSizes, le.calculateSizeRecursive(child, childAvailableSize))
		}

		// If size is fixed, we don't need to calculate based on children
		if hasFixedWidth && hasFixedHeight {
			calculatedSize = fixedSize
		} else {
			// Calculate container size based on children layout (wrapping)
			contentHeight := float32(0.0)
			currentLineWidth := float32(0.0)
			currentLineMaxHeight := float32(0.0)
			maxLineWidth := float32(0.0)

			for i, child := range c.Children() {
				childSize := childrenSizes[i]
				if child.Pos().Type == PositionTypeAbsolute {
					// Skip absolutely positioned children for wrapping calculations for now (will handle at the end)
					continue
				}

				// Determine if wrapping is needed. Use childAvailableSize.X as the limit.
				// Wrap if not the first element on the line and adding it exceeds available width.
				// Also wrap if the child itself forces a block display.
				needsWrap := (currentLineWidth > 0 && currentLineWidth+childSize.X > childAvailableSize.X) ||
					child.Display() == DisplayBlock

				if needsWrap {
					// Finish previous line
					contentHeight += currentLineMaxHeight
					maxLineWidth = max(maxLineWidth, currentLineWidth)

					// Start new line
					currentLineWidth = childSize.X + 2*c.Margin().X
					currentLineMaxHeight = childSize.Y + 2*child.Margin().Y
				} else {
					// Add to current line
					currentLineWidth += childSize.X + 2*c.Margin().X
					currentLineMaxHeight = max(currentLineMaxHeight, childSize.Y+2*child.Margin().Y)
				}

				// If a single child is wider than available, it dictates the max width
				maxLineWidth = max(maxLineWidth, childSize.X)
			}

			// After processing all relative children, finalize the last line's height
			// check for absolutely positioned children that might affect the height and width
			for i, child := range c.Children() {
				if child.Pos().Type == PositionTypeAbsolute {
					childSize := childrenSizes[i]
					contentHeight = max(contentHeight, childSize.Y)
					maxLineWidth = max(maxLineWidth, childSize.X)
				}
			}

			// Add the height of the last line to the total content height
			// This is necessary to ensure the container's height accounts for all children.
			// If the last line was not empty, add its height to the total content height.

			// Add the height of the last line
			contentHeight += currentLineMaxHeight
			// Ensure the max width accounts for the last line's width
			maxLineWidth = max(maxLineWidth, currentLineWidth)

			// Use calculated dimensions unless overridden by fixed dimensions
			calculatedSize.X = maxLineWidth
			calculatedSize.Y = contentHeight
			if hasFixedWidth {
				calculatedSize.X = fixedSize.X
			}
			if hasFixedHeight {
				calculatedSize.Y = fixedSize.Y
			}
		}

	case *Text:
		// TODO: Font size should ideally come from component style/props
		fontSize := float32(24.0)
		// TODO: Padding should ideally come from component style/props
		padding := float32(30.0)
		width := le.CalculateTextWidth(c.Content, fontSize) + padding
		height := fontSize // Basic height based on font size
		calculatedSize = Vec2{X: width, Y: height}
		// Respect fixed size if set
		if hasFixedWidth {
			calculatedSize.X = fixedSize.X
		}
		if hasFixedHeight {
			calculatedSize.Y = fixedSize.Y
		}

	case *Button:
		textWidth := le.CalculateTextWidth(c.Label, c.FontSize()) + 30.0
		width := textWidth
		height := c.FontSize()
		calculatedSize = Vec2{X: width, Y: height}
		// Respect fixed size if set
		if hasFixedWidth {
			calculatedSize.X = fixedSize.X
		}
		if hasFixedHeight {
			calculatedSize.Y = fixedSize.Y
		}

	case *Image:
		calculatedSize = c.Size()
	default:
		// Consider logging a warning instead of panicking for unknown types?
		panic(fmt.Sprintf("Unsupported component type for size calculation: %T", comp))
		// Or return a zero size:
		// calculatedSize = Vec2{X: 0, Y: 0}
	}

	if !hasFixedWidth {
		calculatedSize.X += 2 * (comp.Padding().X + comp.Border().X) // Include padding and border in width
	}
	if !hasFixedHeight {
		calculatedSize.Y += 2 * (comp.Padding().Y + comp.Border().Y) // Include padding and border in height
	}

	comp.setSize(calculatedSize) // Set the calculated size on the component
	// Use the internal setter
	return calculatedSize
}

// calculatePositionRecursive determines the position of each component relative
// to its parent, starting from the root and moving down.
// parentTopLeft is the absolute coordinate where the parent *starts* placing this component.
func (le *LayoutEngine) calculatePositionRecursive(comp IComponent, parentTopLeft Vec2) {
	compPosInfo := comp.Pos()
	if compPosInfo.Type == PositionTypeAbsolute {
		parentTopLeft.X = compPosInfo.X
		parentTopLeft.Y = compPosInfo.Y
	}

	containerSize := comp.Size()
	containerSize.X -= comp.Padding().X + comp.Border().X
	containerSize.Y -= comp.Padding().Y + comp.Border().Y
	var contentOrigin Vec2 = parentTopLeft

	// Now, handle the layout *within* this component (positioning its children)
	switch c := comp.(type) {
	case *Container:
		// Use the determined contentOrigin for placing children.
		containerContentOrigin := contentOrigin
		// Use the size calculated in the first pass.

		// Track position within the current line for relative layout.
		currentLineXOffset := float32(0.0) + comp.Padding().X + comp.Border().X
		currentLineYOffset := float32(0.0) + comp.Padding().Y + comp.Border().Y
		currentLineMaxHeight := float32(0.0)

		// Iterate through children to position them.
		for _, child := range c.Children() {
			childPosInfo := child.Pos()
			childSize := child.Size()
			childSize.X += 2 * child.Margin().X // Include margin in size for wrapping calculations
			childSize.Y += 2 * child.Margin().Y // Include margin in size for wrapping calculations

			// Handle absolutely positioned children first.
			if childPosInfo.Type == PositionTypeAbsolute {
				// Position the absolute child relative to *this* container's content origin.
				// The recursive call will handle calculating its final screen position.
				le.calculatePositionRecursive(child, containerContentOrigin)
				continue // Skip relative flow calculations for this child.
			}

			// --- Relative Child Layout Logic ---
			needsWrap := false
			// Check for wrapping only if container has a positive width.
			if containerSize.X > 0 {
				// Wrap if:
				// 1. Not the first element on the line (X offset > 0) AND
				// 2. Adding the child exceeds the container's width OR
				// 3. The child forces a block display (e.g., like a <p> or <div> in HTML).
				needsWrap = (currentLineXOffset > 0 && currentLineXOffset+childSize.X > containerSize.X) ||
					(child.Display() == DisplayBlock) // Assuming DisplayBlock constant exists
			} else if currentLineXOffset > 0 {
				// If container has no width, wrap after every element to prevent infinite horizontal layout.
				needsWrap = true
			}

			// If wrapping is needed, move to the start of the next line.
			if needsWrap {
				currentLineYOffset += currentLineMaxHeight              // Add height of the completed line.
				currentLineXOffset = comp.Padding().X + comp.Border().X // Reset X offset for the new line.
				currentLineMaxHeight = 0                                // Reset max height for the new line.
			}

			// Calculate the child's position *relative* to this container's content origin.
			childRelativePos := Position{
				Type: PositionTypeRelative,                  // Ensure type is set correctly.
				X:    currentLineXOffset + child.Margin().X, // Add margin to the relative position
				Y:    currentLineYOffset + child.Margin().Y, // Add margin to the relative position
			}

			// Attempt to set the calculated relative position on the child component.
			child.setPos(childRelativePos)

			// Calculate the child's absolute top-left screen coordinate.
			// This is needed as the reference point (`parentTopLeft`) for positioning the child's *own* children.
			childAbsoluteTopLeft := Vec2{
				X: containerContentOrigin.X + currentLineXOffset, // Parent origin + child relative offset
				Y: containerContentOrigin.Y + currentLineYOffset, // Parent origin + child relative offset
			}

			// Recursively call positioning for the child's children.
			le.calculatePositionRecursive(child, childAbsoluteTopLeft)

			// Advance the X offset on the current line for the *next* sibling.
			currentLineXOffset += childSize.X
			// Update the maximum height encountered on the current line.
			currentLineMaxHeight = max(currentLineMaxHeight, childSize.Y)
		}

	case *Text, *Button, *Image:
		// Leaf node. Position was set by its parent container if relative.
		// Absolute positioning was handled when calculating contentOrigin.
		// No children to position.
		break

	default:
		// Should not happen if calculateSizeRecursive covers all component types.
		panic(fmt.Sprintf("Unsupported component type for position calculation: %T", comp))
	}
}
