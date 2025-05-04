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
	alive              map[string]bool
	count              map[string]int
	state              map[string]ComponentState
}

type ComponentState struct {
}

// Reset at the top of each frame:
func (le *LayoutEngine) BeginLayout() {
	le.alive = make(map[string]bool, len(le.state))
	le.count = make(map[string]int)
}

// Prune at the end of each frame:
func (le *LayoutEngine) EndLayout() {
	for id := range le.state {
		if !le.alive[id] {
			delete(le.state, id)
		}
	}
}

// NewLayoutEngine creates a layout engine.
func NewLayoutEngine(f func(string, float32) float32) *LayoutEngine {
	return &LayoutEngine{
		CalculateTextWidth: f,
		alive:              make(map[string]bool),
		count:              make(map[string]int),
		state:              make(map[string]ComponentState),
	}
}
func (le *LayoutEngine) Layout(root IComponent,
	origin Vec2,
	availableSize Vec2) {
	// ─── PASS 0: assign a stable full‐path ID to every node ───
	// Seed with “root” so your first widgets come out as “root/container#0(…)”
	le.assignIDsRecursive(root)

	// ─── PASS 1: intrinsic size (bottom‐up) ───
	le.calculateSizeRecursive(root, availableSize)

	// ─── PASS 2: positions (top‐down) ───
	le.calculatePositionRecursive(root, origin)

	// Optional debug dump
	// fmt.Println("--- Layout Complete ---")
	// le.printComponentTree(root, "")
	// fmt.Println("-----------------------")
}

func (le *LayoutEngine) registerID(fullID string) {
	if !le.alive[fullID] {
		le.alive[fullID] = true
		if _, ok := le.state[fullID]; !ok {
			le.state[fullID] = ComponentState{}
		}
	}
}

func (le *LayoutEngine) nextFullID(parent, widgetType, userID string) string {
	key := parent + "/" + widgetType
	idx := le.count[key]
	le.count[key] = idx + 1

	// fallback if caller didn't supply an ID
	id := userID

	segment := fmt.Sprintf("%s#%d(%s)", widgetType, idx, id)
	full := parent + "/" + segment

	le.registerID(full)
	return full
}

// assignIDsRecursive walks the tree and gives each component
// a full‐path ID before any layout work happens.
func (le *LayoutEngine) assignIDsRecursive(comp IComponent) {
	// Ask your App (or UIContext) for the new full ID:
	//   widgetType := comp.Type() // e.g. "container", "button", "text"
	//   userID     := comp.UserID() // whatever the caller set, or "" if none
	fullID := "root"
	if comp.Parent() != nil {
		fullID = comp.Parent().FullID()
	}
	newFullID := le.nextFullID(fullID, comp.Kind().String(), comp.ID())
	comp.setFullID(newFullID)

	// Recurse into children
	for _, child := range comp.Children() {
		le.assignIDsRecursive(child)
	}
}

// printComponentTree is a helper for debugging the layout structure.
func (le *LayoutEngine) printComponentTree(comp IComponent, indent string) {

	if comp.Display() == DisplayNone {
		return
	}
	displayStr := comp.FullID()

	fmt.Printf("%s[%s: Size:%.1f,%.1f Pos:%.1f,%.1f (%v) Padding:%.1f,%.1f Margin:%.1f,%.1f Gap:%.1f,%.1f Border:%.1f,%.1f]\n",
		indent, displayStr, comp.Size().X, comp.Size().Y,
		comp.Pos().X, comp.Pos().Y, comp.Pos().Type,
		comp.Padding().X, comp.Padding().Y,
		comp.Margin().X, comp.Margin().Y,
		comp.Gap().X, comp.Gap().Y,
		comp.Border().X, comp.Border().Y)
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
// sizes otherwise. availableSize provides the constraints from the parent.
func (le *LayoutEngine) calculateSizeRecursive(comp IComponent, availableSize Vec2) Vec2 {
	if comp == nil || comp.Display() == DisplayNone {
		// Skip if component is nil or marked as not displayed.
		return Vec2{X: 0, Y: 0}
	}
	fixedSize := comp.Size() // User-defined fixed size
	hasFixedWidth := fixedSize.X > 0
	hasFixedHeight := fixedSize.Y > 0

	padding := comp.Padding()
	border := comp.Border()
	paddingAndBorderX := padding.X + border.X
	paddingAndBorderY := padding.Y + border.Y

	var calculatedContentSize Vec2 // Size needed by content only

	switch c := comp.(type) {
	case *Container:
		// If container has fixed dimensions, children layout doesn't affect its size.
		// However, we still need to calculate children sizes recursively.
		// Calculate the space available *inside* this container for children.
		childAvailableSize := availableSize
		childAvailableSize.X -= 2 * paddingAndBorderX
		childAvailableSize.Y -= 2 * paddingAndBorderY

		// If container has fixed size, use that to constrain children instead.
		if hasFixedWidth {
			childAvailableSize.X = fixedSize.X - 2*paddingAndBorderX
		}
		if hasFixedHeight {
			childAvailableSize.Y = fixedSize.Y - 2*paddingAndBorderY
		}
		// Ensure available size for children isn't negative.
		childAvailableSize.X = max(0, childAvailableSize.X)
		childAvailableSize.Y = max(0, childAvailableSize.Y)

		// Calculate children sizes first, passing the constrained available size.
		var childrenSizes []Vec2
		for _, child := range c.Children() {
			childrenSizes = append(childrenSizes, le.calculateSizeRecursive(child, childAvailableSize))
		}

		// If size is fully fixed, we don't need to calculate based on children layout.
		// The recursive calls above were still needed for the children themselves.
		if hasFixedWidth && hasFixedHeight {
			// Use fixed size directly (padding/border are included implicitly)
			calculatedContentSize = Vec2{
				X: max(0, fixedSize.X-2*paddingAndBorderX),
				Y: max(0, fixedSize.Y-2*paddingAndBorderY),
			}
			// Skip layout calculation below
		} else {
			// Calculate container's content size based on children layout (wrapping).
			contentHeight := float32(0.0)
			currentLineWidth := float32(0.0)
			currentLineMaxHeight := float32(0.0)
			maxContentWidth := float32(0.0) // Tracks the widest line of content
			numberOfChildrenInLine := 0

			for i, child := range c.Children() {
				// Use the size calculated in the recursive call
				childSize := childrenSizes[i]
				childMargin := child.Margin()

				// Size used for layout includes the child's margins
				childLayoutWidth := childSize.X + 2*childMargin.X
				childLayoutHeight := childSize.Y + 2*childMargin.Y

				if child.Pos().Type == PositionTypeAbsolute {
					// Absolutely positioned children don't participate in flow layout.
					continue
				}

				// Determine if wrapping is needed. Use childAvailableSize.X as the limit.
				gapX := float32(0.0)
				if numberOfChildrenInLine > 0 {
					gapX = c.Gap().X // Gap only applies after the first item
				}

				needsWrap := (numberOfChildrenInLine > 0 && currentLineWidth+gapX+childLayoutWidth > childAvailableSize.X) ||
					child.Display() == DisplayBlock

				if needsWrap && numberOfChildrenInLine > 0 { // Ensure wrap only happens if line isn't empty
					// Finish previous line
					contentHeight += currentLineMaxHeight
					if numberOfChildrenInLine > 1 { // Only add line gap if there was more than one item
						contentHeight += c.Gap().Y
					}
					maxContentWidth = max(maxContentWidth, currentLineWidth)

					// Start new line with the current child
					currentLineWidth = childLayoutWidth
					currentLineMaxHeight = childLayoutHeight
					numberOfChildrenInLine = 1
				} else {
					// Add to current line
					if numberOfChildrenInLine > 0 {
						currentLineWidth += c.Gap().X // Add gap before adding the child
					}
					currentLineWidth += childLayoutWidth
					currentLineMaxHeight = max(currentLineMaxHeight, childLayoutHeight)
					numberOfChildrenInLine++
				}

				// If a single child is wider than available, it dictates the minimum content width needed.
				// This ensures maxContentWidth captures overflowing individual children.
				maxContentWidth = max(maxContentWidth, childLayoutWidth)
			}

			// Add the height of the last line
			if numberOfChildrenInLine > 0 {
				contentHeight += currentLineMaxHeight
				maxContentWidth = max(maxContentWidth, currentLineWidth)
			}

			// Store the calculated size needed *just for the content*.
			calculatedContentSize.X = maxContentWidth
			calculatedContentSize.Y = contentHeight
		}

	case *Text:
		// TODO: Font size should ideally come from component style/props
		fontSize := float32(24.0)
		// TODO: Padding should ideally come from component style/props
		// Note: Text usually uses its own padding, not the component's general padding.
		// Assuming CalculateTextWidth gives raw text width. Add specific text padding if needed.
		textPadding := float32(0.0) // Example: Add text-specific padding if necessary
		width := le.CalculateTextWidth(c.Content, fontSize) + 2*textPadding
		height := fontSize + 2*textPadding // Basic height based on font size
		calculatedContentSize = Vec2{X: width, Y: height}

	case *Button:
		// Assume button includes internal padding within its calculation logic
		// For example, CalculateTextWidth might already include some padding.
		// Or, add explicit button padding here.
		buttonPaddingX := float32(15.0) // Example internal padding
		buttonPaddingY := float32(5.0)  // Example internal padding
		textWidth := le.CalculateTextWidth(c.Label, c.FontSize())
		width := textWidth + 2*buttonPaddingX
		height := c.FontSize() + 2*buttonPaddingY
		calculatedContentSize = Vec2{X: width, Y: height}

	case *Image:
		// Assume c.Size() returns the intrinsic size of the image content.
		intrinsicSize := c.Size()
		calculatedContentSize = intrinsicSize

	default:
		// Return zero size for unknown types, maybe log a warning.
		fmt.Printf("Warning: Unsupported component type for size calculation: %T\n", comp)
		calculatedContentSize = Vec2{X: 0, Y: 0}
		// Or panic:
		// panic(fmt.Sprintf("Unsupported component type for size calculation: %T", comp))
	}

	// Calculate the total "natural" size including padding and border.
	naturalWidth := calculatedContentSize.X + 2*paddingAndBorderX
	naturalHeight := calculatedContentSize.Y + 2*paddingAndBorderY

	// Final size respects fixed dimensions if they are set.
	finalSize := Vec2{
		X: naturalWidth,
		Y: naturalHeight,
	}
	if hasFixedWidth {
		finalSize.X = fixedSize.X
	}
	if hasFixedHeight {
		finalSize.Y = fixedSize.Y
	}

	// Ensure final size is not negative (e.g., if fixed size is smaller than padding/border).
	finalSize.X = max(0, finalSize.X)
	finalSize.Y = max(0, finalSize.Y)

	// Set the calculated size on the component.
	comp.setSize(finalSize)

	// Return the final calculated size for the parent's layout logic.
	return finalSize
}

// calculatePositionRecursive determines the position of each component relative
// to its parent, starting from the root and moving down.
// parentTopLeft is the absolute coordinate where the parent *starts* placing this component.
func (le *LayoutEngine) calculatePositionRecursive(comp IComponent, parentTopLeft Vec2) {
	if comp == nil || comp.Display() == DisplayNone {
		// Skip if component is nil or marked as not displayed.
		return
	}
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
		numberOfChildrenInLine := 0

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
				numberOfChildrenInLine = 1                              // Reset for the new line
				currentLineYOffset += currentLineMaxHeight + c.Gap().Y  // Add height of the completed line.
				currentLineXOffset = comp.Padding().X + comp.Border().X // Reset X offset for the new line.
				currentLineMaxHeight = 0                                // Reset max height for the new line.
			} else {
				numberOfChildrenInLine++ // Increment the number of children in the current line.
			}
			if numberOfChildrenInLine > 1 {
				currentLineXOffset += c.Gap().X // Add gap between children
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
