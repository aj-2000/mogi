package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"math/rand"
	"strconv"
)

var colors = []func() common.ColorRGBA{
	consts.ColorRed,
	consts.ColorGreen,
	consts.ColorBlue,
	consts.ColorYellow,
	consts.ColorOrange,
	consts.ColorPurple,
	consts.ColorCyan,
	consts.ColorMagenta,
	consts.ColorWhite,
}

func BoxesOneComponent() common.IComponent {
	const numBoxes = 5000 // Number of child boxes to generate

	children := make([]common.IComponent, numBoxes)
	for i := range children {
		width := float32(rand.Intn(11)) + 4.05
		height := float32(rand.Intn(11)) + 4.05
		color := colors[i%len(colors)]()

		child := common.NewContainer().
			SetID("box_" + strconv.Itoa(i+1)).
			SetBackgroundColor(color).
			SetSize(common.Vec2{X: float32(width), Y: float32(height)}).
			SetPosition(common.Position{
				Type: common.PositionTypeRelative,
			})

		children[i] = child
	}

	return common.NewContainer().
		SetID("boxes_container").
		SetBackgroundColor(consts.ColorBlack()).
		AddChildren(children...).
		SetPosition(common.Position{
			Type: common.PositionTypeRelative,
		})
}

// recursiveHelper creates a container and potentially fills it with nested children.
// Only leaf nodes (at maxLevel) will have an explicit size set.
func recursiveHelper(currentLevel, maxLevel int, baseID string, maxChildrenPerNode int) common.IComponent {
	// 1. Create the container for the current level
	id := baseID + "_lvl" + strconv.Itoa(currentLevel) + "_r" + strconv.Itoa(rand.Intn(10000))
	color := colors[rand.Intn(len(colors))]() // Assign color regardless of level

	container := common.NewContainer().
		SetID(id).
		SetBackgroundColor(color).
		SetPosition(common.Position{
			Type: common.PositionTypeRelative, // Position relative to parent
		})
		// Optional: Add layout properties if your library supports them for auto-sizing
		// e.g., SetLayoutType(common.LayoutTypeFlex), SetFlexDirection(common.FlexDirectionRow) etc.

	// 2. Check if this is a leaf node (at the maximum desired level)
	if currentLevel >= maxLevel {
		// --- LEAF NODE ---
		// Calculate and set size ONLY for leaf nodes
		// Use simpler size calculation as scaling factor isn't needed here
		width := float32(rand.Intn(11)) + 4.05 // e.g., 4.05 to 14.05
		height := float32(rand.Intn(11)) + 4.05

		container.SetSize(common.Vec2{X: width, Y: height})
		// No children for leaf nodes, return the sized container
		return container
	}

	// 3. --- INTERMEDIATE NODE ---
	// If not at max level, generate children.
	// DO NOT set an explicit size on this container. Its size should be
	// determined by its children and the layout system.
	numChildren := rand.Intn(maxChildrenPerNode + 1) // 0 to maxChildrenPerNode
	if numChildren > 0 {
		children := make([]common.IComponent, numChildren)
		for i := 0; i < numChildren; i++ {
			// Recursively call for the next level
			childID := id + "_c" + strconv.Itoa(i)
			children[i] = recursiveHelper(currentLevel+1, maxLevel, childID, maxChildrenPerNode)
		}
		// Add the generated children to the current container
		container.AddChildren(children...)
	}
	// Return the intermediate container (without explicit size)
	return container
}

// Creates a root container with potentially nested children up to 'maxLevel' deep.
func BoxesNLevelComponent(maxLevel int, numRootChildren int, maxChildrenPerNode int) common.IComponent {
	// Seed random number generator (important!)
	// rand.Seed(time.Now().UnixNano()) // Do this once at app start

	// --- Input Validation ---
	if maxLevel < 1 {
		maxLevel = 1
	}
	if numRootChildren < 0 {
		numRootChildren = 0
	}
	if maxChildrenPerNode < 0 {
		maxChildrenPerNode = 0
	}

	// --- Create Root Children ---
	rootChildren := make([]common.IComponent, numRootChildren)
	for i := 0; i < numRootChildren; i++ {
		rootChildID := "root_" + strconv.Itoa(i)
		rootChildren[i] = recursiveHelper(1, maxLevel, rootChildID, maxChildrenPerNode)
	}

	// --- Create the Main Root Container ---
	rootContainer := common.NewContainer().
		SetID("boxes_n_level_root").
		SetBackgroundColor(consts.ColorBlack()).
		AddChildren(rootChildren...).
		SetPosition(common.Position{
			Type: common.PositionTypeRelative,
		})
		// The root container also doesn't get an explicit size here.
		// It will depend on the context it's placed in.
		// You might need to configure layout properties on this container
		// (e.g., flexbox, grid) depending on your GoUI library,
		// so it knows how to arrange the rootChildren.

	return rootContainer
}
