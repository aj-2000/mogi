package color

const (
	inv6    = 1.0 / 6.0 // reuse instead of writing 1.0/6.0 each time
	twoInv3 = 2.0 / 3.0 // reuse instead of 2.0/3.0
)

func hue2rgb(p, q, t float32) float32 {
	if t < 0 {
		t += 1
	} else if t > 1 {
		t -= 1
	}

	d := q - p

	if t < inv6 {
		// rising edge
		return p + d*6*t
	} else if t < 0.5 {
		// top plateau
		return q
	} else if t < twoInv3 {
		// falling edge
		return p + d*6*(twoInv3-t)
	}
	// bottom plateau
	return p
}
