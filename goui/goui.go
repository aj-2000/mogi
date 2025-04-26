package goui

import "fmt"

func Init() {
	println("GOUI Init!")
}

// Button defines a simple button widget.
type Button struct {
	X, Y, Width, Height float32
	Label               string
	OnClick             func()
}

// Draw renders the button.
func (b *Button) Draw() {
	// Here you would render a rectangle (button) and display text
	// with OpenGL. For simplicity, this is just a placeholder.

	fmt.Printf("Drawing button: %s at (%.2f, %.2f)\n", b.Label, b.X, b.Y)

	// Check for click (use some event handler)
	if isMouseOver(b) && isMouseClicked() {
		b.OnClick()
	}
}

// Helper function to check if mouse is over the button
func isMouseOver(b *Button) bool {
	// Placeholder: Assume mouse is over button
	return true
}

// Placeholder function for mouse click event
func isMouseClicked() bool {
	// Placeholder: Assume mouse is clicked
	return true
}


func Begin() {

}

func End() {

}
