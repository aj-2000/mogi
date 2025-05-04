package color

import (
	"strconv"
	"strings"
)

type HSLA struct {
	H, S, L, A float32 // H in [0,360), S/L/A in [0,1]
}

func NewHSLA(h, s, l, a float32) HSLA { return HSLA{h, s, l, a} }
func (h HSLA) ToRGBA() RGBA {
	hh := h.H / 360
	s := h.S
	l := h.L
	a := h.A
	var r, g, b float32
	if s == 0 {
		// achromatic
		r, g, b = l, l, l
	} else {
		var q float32
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, hh+1.0/3.0)
		g = hue2rgb(p, q, hh)
		b = hue2rgb(p, q, hh-1.0/3.0)
	}
	return RGBA{r, g, b, a}
}
func (h HSLA) String() string {
	var sb strings.Builder
	sb.Grow(40) // rough upper bound for "hsla(360.00,1.00,1.00,1.00)"
	sb.WriteString("hsla(")

	buf := make([]byte, 0, 12)
	// Hue
	buf = strconv.AppendFloat(buf[:0], float64(h.H), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	// Saturation
	buf = strconv.AppendFloat(buf[:0], float64(h.S), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	// Lightness
	buf = strconv.AppendFloat(buf[:0], float64(h.L), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	// Alpha
	buf = strconv.AppendFloat(buf[:0], float64(h.A), 'f', 2, 32)
	sb.Write(buf)

	sb.WriteByte(')')
	return sb.String()
}
