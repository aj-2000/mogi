package color

import (
	"strconv"
	"strings"
)

type RGBA struct {
	R, G, B, A float32
}

func NewRGBA(r, g, b, a float32) RGBA { return RGBA{r, g, b, a} }

// BlendOver alpha-composites c over dst.
func (c RGBA) BlendOver(dst RGBA) RGBA {
	invA := 1 - c.A // reuse
	outA := c.A + dst.A*invA
	if outA == 0 {
		return RGBA{} // fully transparent
	}
	// factor dst.A*invA once
	dstMul := dst.A * invA
	return RGBA{
		R: (c.R*c.A + dst.R*dstMul) / outA,
		G: (c.G*c.A + dst.G*dstMul) / outA,
		B: (c.B*c.A + dst.B*dstMul) / outA,
		A: outA,
	}
}

// String returns "rgba(r, g, b, a)", with two decimal places.
// NOTE: it loses precision for large values (e.g. 1.23456789 becomes 1.23).
// It is not intended for high-precision use, but for human-readable output.
func (c RGBA) String() string {
	var sb strings.Builder
	sb.Grow(32) // rough upper bound

	sb.WriteString("rgba(")
	// reuse a small byte-slice for each float
	buf := make([]byte, 0, 12)

	buf = strconv.AppendFloat(buf[:0], float64(c.R), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	buf = strconv.AppendFloat(buf[:0], float64(c.G), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	buf = strconv.AppendFloat(buf[:0], float64(c.B), 'f', 2, 32)
	sb.Write(buf)
	sb.WriteByte(',')

	buf = strconv.AppendFloat(buf[:0], float64(c.A), 'f', 2, 32)
	sb.Write(buf)

	sb.WriteByte(')')
	return sb.String()
}
func (c RGBA) ToRGBA() RGBA {
	return c
}
