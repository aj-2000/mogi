package examples

import (
	"math/rand"
	"strconv"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

func randomColor() color.RGBA {
	return color.RGBA{
		R: float32(rand.Intn(256)) / 255.0,
		G: float32(rand.Intn(256)) / 255.0,
		B: float32(rand.Intn(256)) / 255.0,
		A: 1.0,
	}
}

func BoxesOneComponent(app *mogiApp.App) ui.IComponent {
	const numBoxes = 200

	children := make([]ui.IComponent, numBoxes)
	for i := range children {
		width := float32(rand.Intn(11)) + 4.05
		height := float32(rand.Intn(11)) + 4.05
		id := "box_" + strconv.Itoa(i+1)
		color := randomColor()

		child := app.Container().
			SetID(id).
			SetBackgroundColor(color).
			SetSize(math.Vec2f32{X: width, Y: height})

		children[i] = child
	}

	return app.Container().
		SetID("boxes_container(red)").
		SetBackgroundColor(color.Red).
		AddChildren(children...)
}

func recursiveHelper(app *mogiApp.App, currentLevel, maxLevel int, baseID string, maxChildrenPerNode int) ui.IComponent {
	id := baseID + "_lvl" + strconv.Itoa(currentLevel) + "_r" + strconv.Itoa(rand.Intn(10000))
	color := randomColor()

	container := app.Container().
		SetID(id).
		SetBackgroundColor(color).
		// SetMargin(math.Vec2f32{X: 2, Y: 2}).
		SetPadding(math.Vec2f32{X: 3, Y: 3}).
		SetBorder(math.Vec2f32{X: 2, Y: 2}).
		SetBorderColor(randomColor()).
		SetBorderRadius(2).
		SetGap(math.Vec2f32{X: 2, Y: 2})

	if currentLevel >= maxLevel {
		width := float32(rand.Intn(11)) + 4.05
		height := float32(rand.Intn(11)) + 4.05

		container.SetSize(math.Vec2f32{X: width, Y: height})
		return container
	}

	numChildren := rand.Intn(maxChildrenPerNode + 1)
	if numChildren > 0 {
		children := make([]ui.IComponent, numChildren)
		for i := 0; i < numChildren; i++ {
			childID := id + "_c" + strconv.Itoa(i)
			children[i] = recursiveHelper(app, currentLevel+1, maxLevel, childID, maxChildrenPerNode)
		}
		container.AddChildren(children...)
	}

	return container
}

// TODO: how to prevent by value (app), is it even required?
func BoxesNLevelComponent(app *mogiApp.App, maxLevel int, numRootChildren int, maxChildrenPerNode int) ui.IComponent {
	if maxLevel < 1 {
		maxLevel = 1
	}
	if numRootChildren < 0 {
		numRootChildren = 0
	}
	if maxChildrenPerNode < 0 {
		maxChildrenPerNode = 0
	}

	rootChildren := make([]ui.IComponent, numRootChildren)
	for i := 0; i < numRootChildren; i++ {
		rootChildID := "root_" + strconv.Itoa(i)
		rootChildren[i] = recursiveHelper(app, 1, maxLevel, rootChildID, maxChildrenPerNode)
	}

	rootContainer := app.Container().
		SetID("boxes_n_level_root").
		SetBackgroundColor(color.Black).
		AddChildren(rootChildren...)

	return rootContainer
}
