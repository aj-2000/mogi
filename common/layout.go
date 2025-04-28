package common

import (
	"sort"
)

// LayoutEngine calculates component positions and sizes based on Flexbox rules.
type LayoutEngine struct {
	// Potential future state: Font measurement interface, caching, etc.
	// FontMeasurer FontMeasurer // Interface to measure text size
}

// NewLayoutEngine creates a layout engine.
func NewLayoutEngine( /* fontMeasurer FontMeasurer */ ) *LayoutEngine {
	return &LayoutEngine{
		// FontMeasurer: fontMeasurer,
	}
}

// CalculateLayout computes the layout for the component tree starting at root.
// availableSize is the constraint for the root component.
func (le *LayoutEngine) CalculateLayout(root IComponent, availableSize Vec2) {
	// Start the recursive layout calculation
	le.calculateNodeLayout(root, availableSize, Position{X: 0, Y: 0})
}

// calculateNodeLayout is the recursive core of the layout engine.
// node: The current component to calculate layout for.
// availableSize: The space allocated to this node by its parent.
// offset: The position of this node's top-left corner relative to the root.
func (le *LayoutEngine) calculateNodeLayout(node IComponent, availableSize Vec2, offset Position) {
	// 1. Set initial size based on available space (this might be adjusted by flex)
	node.setSize(availableSize)
	node.setPos(offset) // Set initial position

	// 2. Check if this node is a Flex Container
	container, isContainer := node.(*Container)
	isFlexContainer := isContainer && container.FlexContainerProps().Enabled // Need FlexContainerProps() accessor

	if isFlexContainer && len(node.Children()) > 0 {
		// --- Apply Flexbox Layout to Children ---
		println("Laying out flex container:", node.ID())
		le.layoutFlexContainer(container, availableSize)
	} else {
		// --- Non-Flex Container or Leaf Node ---
		// Children (if any) are laid out recursively, but without flex rules
		// from *this* container. They might still be flex items if the parent
		// was a flex container, or they might become flex containers themselves.
		// For now, we just recurse. A non-flex container doesn't impose layout
		// rules in this basic model (it just takes up its calculated space).
		children := node.Children()
		for _, child := range children {
			// TODO: How to determine available size for children in non-flex?
			// Simplistic: give child its preferred size? Or zero?
			// Or maybe non-flex containers *do* need basic stacking logic?
			// For now, let's assume children of non-flex containers need
			// explicit positioning or are handled by *their* parent's flex rules.
			// We still need to recurse so *their* children get laid out.

			// Calculate child offset relative to the root
			childOffset := Position{
				X: offset.X + child.Pos().X, // Use child's current pos as relative offset
				Y: offset.Y + child.Pos().Y,
			}
			// Calculate child available size (hard problem without layout rules)
			// Let's pass the child's current size for now, assuming it was set somehow
			// or will be determined by its own flex container properties.
			childAvailableSize := child.Size()

			le.calculateNodeLayout(child, childAvailableSize, childOffset)
		}
	}
}

// layoutFlexContainer applies Flexbox rules to the children of a container.
func (le *LayoutEngine) layoutFlexContainer(container *Container, containerSize Vec2) {
	props := container.FlexContainerProps() // Need accessor
	children := container.Children()
	if len(children) == 0 {
		return
	}

	// --- Determine Axes ---
	isRow := props.Direction == FlexDirectionRow || props.Direction == FlexDirectionRowReverse
	isColumn := !isRow
	isReverse := props.Direction == FlexDirectionRowReverse || props.Direction == FlexDirectionColumnReverse

	mainAxisIsX := isRow
	crossAxisIsX := !mainAxisIsX

	// --- Get Container Inner Size (excluding border/padding if you add them) ---
	// For now, assume containerSize is the inner size.
	mainAvailableSize := containerSize.X
	crossAvailableSize := containerSize.Y
	if isColumn {
		mainAvailableSize, crossAvailableSize = crossAvailableSize, mainAvailableSize
	}

	// --- Sort children by Order property ---
	// Create a slice of indices to sort without modifying original children slice directly
	orderIndices := make([]int, len(children))
	for i := range children {
		orderIndices[i] = i
	}
	sort.SliceStable(orderIndices, func(i, j int) bool {
		idxI := orderIndices[i]
		idxJ := orderIndices[j]
		return children[idxI].FlexItem().Order < children[idxJ].FlexItem().Order
	})
	// Create a sorted view or copy if needed for processing
	sortedChildren := make([]IComponent, len(children))
	for i, idx := range orderIndices {
		sortedChildren[i] = children[idx]
	}

	// --- Line Breaking (Simplified: No Wrapping) ---
	// We assume all items fit on a single line for now.
	// TODO: Implement line breaking logic for FlexWrapWrap/WrapReverse

	// --- Calculate Initial Main Sizes (Flex Basis) ---
	totalBasis := float32(0)
	totalGrow := float32(0)
	totalShrink := float32(0)
	// Store calculated basis per item
	calculatedBasis := make([]float32, len(sortedChildren))

	for i, child := range sortedChildren {
		flexItem := child.FlexItem()
		basis := flexItem.Basis

		if basis == FlexBasisAuto || basis < 0 {
			// Use intrinsic size or fallback
			// TODO: Need intrinsic size calculation (e.g., text measurement)
			// Fallback: Use current size if available, otherwise 0 or a default?
			intrinsicSize := child.Size() // Placeholder!
			if mainAxisIsX {
				basis = intrinsicSize.X
			} else {
				basis = intrinsicSize.Y
			}
			// If still zero, maybe default to 0? Or a small value?
			if basis <= 0 {
				basis = 0 // Or some minimum content size?
			}
		}

		// Clamp basis to container size? (Flexbox doesn't strictly do this here)
		calculatedBasis[i] = basis
		totalBasis += basis
		totalGrow += flexItem.Grow
		totalShrink += flexItem.Shrink
	}

	// Add gaps to total basis
	numGaps := float32(len(sortedChildren) - 1)
	if numGaps < 0 {
		numGaps = 0
	}
	totalBasis += numGaps * props.Gap

	// --- Resolve Flexible Lengths ---
	remainingSpace := mainAvailableSize - totalBasis
	finalMainSizes := make([]float32, len(sortedChildren))

	if remainingSpace > 0 && totalGrow > 0 {
		// Distribute positive space (Grow)
		unitGrow := remainingSpace / totalGrow
		for i, child := range sortedChildren {
			finalMainSizes[i] = calculatedBasis[i] + child.FlexItem().Grow*unitGrow
		}
	} else if remainingSpace < 0 && totalShrink > 0 {
		// Distribute negative space (Shrink) - more complex
		// Need to consider scaled shrink factors (basis * shrink)
		totalScaledShrink := float32(0)
		for i, child := range sortedChildren {
			totalScaledShrink += calculatedBasis[i] * child.FlexItem().Shrink
		}

		if totalScaledShrink > 0 {
			unitShrink := -remainingSpace / totalScaledShrink // remainingSpace is negative
			for i, child := range sortedChildren {
				finalMainSizes[i] = calculatedBasis[i] - (calculatedBasis[i]*child.FlexItem().Shrink)*unitShrink
				// Clamp minimum size (cannot shrink below min-content, often 0)
				if finalMainSizes[i] < 0 {
					finalMainSizes[i] = 0
				}
			}
		} else {
			// Cannot shrink, just use basis (items will overflow)
			copy(finalMainSizes, calculatedBasis)
		}

	} else {
		// No flex needed or possible, use calculated basis
		copy(finalMainSizes, calculatedBasis)
	}

	// --- Calculate Cross Sizes and Align Items ---
	// Determine max cross size on the line (needed for align-items baseline/stretch)
	maxCrossSize := float32(0)
	finalCrossSizes := make([]float32, len(sortedChildren))

	for i, child := range sortedChildren {
		// Determine cross size based on child's preference or stretch
		alignSelf := child.FlexItem().AlignSelf
		if alignSelf == AlignSelfAuto {
			alignSelf = props.AlignItems // Inherit from container
		}

		crossSize := float32(0)
		if alignSelf == AlignItemsStretch && crossAvailableSize > 0 {
			// Stretch if container has a definite cross size
			crossSize = crossAvailableSize
		} else {
			// Use intrinsic cross size
			// TODO: Need intrinsic size calculation
			intrinsicSize := child.Size() // Placeholder!
			if crossAxisIsX {
				crossSize = intrinsicSize.X
			} else {
				crossSize = intrinsicSize.Y
			}
			// Fallback if intrinsic is zero
			if crossSize <= 0 {
				// Maybe use the available cross size? Or 0?
				// Let's use 0 for now if intrinsic is unknown/zero
				crossSize = 0
			}
		}
		// Clamp cross size if it exceeds available?
		if crossAvailableSize > 0 && crossSize > crossAvailableSize {
			// crossSize = crossAvailableSize // Optional clamping
		}

		finalCrossSizes[i] = crossSize
		if crossSize > maxCrossSize {
			maxCrossSize = crossSize
		}
	}
	// If container has no definite cross size, maxCrossSize might be 0.
	// If AlignItems is Stretch, and container cross size *is* definite,
	// items should stretch to that size. Our logic above handles this.

	// --- Position Items along Main Axis (Justify Content) ---
	actualTotalMainSize := float32(0)
	for _, size := range finalMainSizes {
		actualTotalMainSize += size
	}
	actualTotalMainSize += numGaps * props.Gap

	mainRemainingSpace := mainAvailableSize - actualTotalMainSize
	mainOffset := float32(0) // Offset for the first item

	numItemsF := float32(len(sortedChildren))

	switch props.Justify {
	case JustifyContentFlexEnd:
		mainOffset = mainRemainingSpace
	case JustifyContentCenter:
		mainOffset = mainRemainingSpace / 2.0
	case JustifyContentSpaceBetween:
		if numItemsF > 1 {
			mainOffset = 0 // First item at start
			// Gap is distributed; props.Gap is ignored here
		} else {
			mainOffset = mainRemainingSpace / 2.0 // Center single item
		}
	case JustifyContentSpaceAround:
		if numItemsF > 0 {
			space := mainRemainingSpace / numItemsF
			mainOffset = space / 2.0
			// Gap is distributed; props.Gap is ignored here
		}
	case JustifyContentSpaceEvenly:
		if numItemsF > 0 {
			space := mainRemainingSpace / (numItemsF + 1)
			mainOffset = space
			// Gap is distributed; props.Gap is ignored here
		}
	case JustifyContentFlexStart:
		fallthrough // Explicit fallthrough
	default:
		mainOffset = 0
	}

	// --- Position Items along Cross Axis (Align Items / Align Self) ---
	// And Set Final Size/Position for each child

	currentMainPos := mainOffset
	if isReverse { // Start from the end if reversed
		currentMainPos = mainAvailableSize - mainOffset
	}

	containerPos := container.Pos() // Get container's absolute position

	for i, child := range sortedChildren {
		mainSize := finalMainSizes[i]
		crossSize := finalCrossSizes[i]

		// Determine cross axis offset
		alignSelf := child.FlexItem().AlignSelf
		if alignSelf == AlignSelfAuto {
			alignSelf = props.AlignItems
		}

		crossOffset := float32(0)
		availableCross := maxCrossSize // Use max line size for alignment baseline
		if props.AlignItems == AlignItemsStretch && crossAvailableSize > 0 {
			// If container stretches, alignment context is container's cross size
			availableCross = crossAvailableSize
		}

		switch alignSelf {
		case AlignItemsFlexEnd:
			crossOffset = availableCross - crossSize
		case AlignItemsCenter:
			crossOffset = (availableCross - crossSize) / 2.0
		case AlignItemsStretch:
			// Size already set, offset is 0 unless container cross size > maxCrossSize
			if crossAvailableSize > maxCrossSize {
				// Center the stretched line within the container cross axis
				crossOffset = (crossAvailableSize - maxCrossSize) / 2.0
			} else {
				crossOffset = 0
			}
			// Ensure the item's cross size is stretched if needed
			if crossAvailableSize > 0 {
				crossSize = crossAvailableSize // Item stretches to fill container cross axis
			}

		case AlignItemsBaseline:
			// TODO: Baseline alignment requires font metrics - complex!
			// Fallback to FlexStart for now.
			crossOffset = 0
		case AlignItemsFlexStart:
			fallthrough
		default:
			crossOffset = 0
		}

		// --- Calculate Final Absolute Position and Size ---
		var childPos Position
		var childSize Vec2

		if isReverse {
			currentMainPos -= mainSize // Move cursor backward
		}

		if mainAxisIsX { // Row direction
			childPos.X = containerPos.X + currentMainPos
			childPos.Y = containerPos.Y + crossOffset
			childSize.X = mainSize
			childSize.Y = crossSize
		} else { // Column direction
			childPos.X = containerPos.X + crossOffset
			childPos.Y = containerPos.Y + currentMainPos
			childSize.X = crossSize
			childSize.Y = mainSize
		}

		child.setPos(childPos)
		child.setSize(childSize)

		// --- Recurse: Layout the child's own children ---
		// Pass the newly calculated size as the available size for the child.
		le.calculateNodeLayout(child, childSize, childPos)

		// --- Advance Main Position ---
		advance := mainSize + props.Gap
		// Handle spacing for justify content modes that distribute space
		if numItemsF > 1 {
			switch props.Justify {
			case JustifyContentSpaceBetween:
				advance = mainSize + mainRemainingSpace/(numItemsF-1)
			case JustifyContentSpaceAround:
				advance = mainSize + mainRemainingSpace/numItemsF
			case JustifyContentSpaceEvenly:
				advance = mainSize + mainRemainingSpace/(numItemsF+1)
			}
		}

		if isReverse {
			currentMainPos -= props.Gap // Move cursor backward for gap
			if props.Justify == JustifyContentSpaceBetween && i < len(sortedChildren)-1 {
				currentMainPos -= mainRemainingSpace / (numItemsF - 1)
			} else if props.Justify == JustifyContentSpaceAround {
				currentMainPos -= mainRemainingSpace / numItemsF
			} else if props.Justify == JustifyContentSpaceEvenly {
				currentMainPos -= mainRemainingSpace / (numItemsF + 1)
			}
		} else {
			currentMainPos += advance // Move cursor forward
		}
	}
}

// Helper needed in Container to access props
func (c *Container) FlexContainerProps() FlexContainerProps {
	return c.flexContainerProps
}
