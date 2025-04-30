package examples

import (
	"GoUI/common"
	"GoUI/consts"
	"math/rand"
	"strconv"
)

func randomColor() common.ColorRGBA {
	return common.ColorRGBA{
		R: float32(rand.Intn(256)) / 255.0,
		G: float32(rand.Intn(256)) / 255.0,
		B: float32(rand.Intn(256)) / 255.0,
		A: 1.0,
	}
}

func BoxesOneComponent() common.IComponent {
	const numBoxes = 5000

	children := make([]common.IComponent, numBoxes)
	for i := range children {
		width := float32(rand.Intn(11)) + 4.05
		height := float32(rand.Intn(11)) + 4.05
		id := "box_" + strconv.Itoa(i+1)
		color := randomColor()

		child := common.NewContainer().
			SetID(id).
			SetBackgroundColor(color).
			SetSize(common.Vec2{X: width, Y: height})

		children[i] = child
	}

	return common.NewContainer().
		SetID("boxes_container").
		SetBackgroundColor(consts.ColorBlack()).
		AddChildren(children...)
}

func recursiveHelper(currentLevel, maxLevel int, baseID string, maxChildrenPerNode int) common.IComponent {
	id := baseID + "_lvl" + strconv.Itoa(currentLevel) + "_r" + strconv.Itoa(rand.Intn(10000))
	color := randomColor()

	container := common.NewContainer().
		SetID(id).
		SetBackgroundColor(color)

	if currentLevel >= maxLevel {
		width := float32(rand.Intn(11)) + 4.05
		height := float32(rand.Intn(11)) + 4.05

		container.SetSize(common.Vec2{X: width, Y: height})
		return container
	}

	numChildren := rand.Intn(maxChildrenPerNode + 1)
	if numChildren > 0 {
		children := make([]common.IComponent, numChildren)
		for i := 0; i < numChildren; i++ {
			childID := id + "_c" + strconv.Itoa(i)
			children[i] = recursiveHelper(currentLevel+1, maxLevel, childID, maxChildrenPerNode)
		}
		container.AddChildren(children...)
	}

	return container
}

func BoxesNLevelComponent(maxLevel int, numRootChildren int, maxChildrenPerNode int) common.IComponent {
	if maxLevel < 1 {
		maxLevel = 1
	}
	if numRootChildren < 0 {
		numRootChildren = 0
	}
	if maxChildrenPerNode < 0 {
		maxChildrenPerNode = 0
	}

	rootChildren := make([]common.IComponent, numRootChildren)
	for i := 0; i < numRootChildren; i++ {
		rootChildID := "root_" + strconv.Itoa(i)
		rootChildren[i] = recursiveHelper(1, maxLevel, rootChildID, maxChildrenPerNode)
	}

	rootContainer := common.NewContainer().
		SetID("boxes_n_level_root").
		SetBackgroundColor(consts.ColorBlack()).
		AddChildren(rootChildren...)

	return rootContainer
}
