package color

import (
	"fmt"
)

type Hex struct {
	hex  string
	rgba RGBA
}

func NewHex(input string) Hex {
	// 1) collect up to 8 hex digits, skipping whitespace and optional '#'
	var raw [8]byte
	n := 0
	for i := 0; i < len(input) && n < 8; i++ {
		c := input[i]
		switch {
		case c == '#', c == ' ', c == '\t', c == '\n', c == '\r':
			// skip
		case '0' <= c && c <= '9', 'A' <= c && c <= 'F', 'a' <= c && c <= 'f':
			// normalize to uppercase
			if c >= 'a' && c <= 'f' {
				c -= 'a' - 'A'
			}
			raw[n] = c
			n++
		default:
			panic(fmt.Sprintf("invalid hex character %q in %q", c, input))
		}
	}

	// 2) expand to exactly 8 digits in buf
	var buf [8]byte
	switch n {
	case 3:
		// #RGB → R R  G G  B B  FF
		buf[0], buf[1] = raw[0], raw[0]
		buf[2], buf[3] = raw[1], raw[1]
		buf[4], buf[5] = raw[2], raw[2]
		buf[6], buf[7] = 'F', 'F'
	case 4:
		// #RGBA → R R  G G  B B  A A
		buf[0], buf[1] = raw[0], raw[0]
		buf[2], buf[3] = raw[1], raw[1]
		buf[4], buf[5] = raw[2], raw[2]
		buf[6], buf[7] = raw[3], raw[3]
	case 6:
		// #RRGGBB → RRGGBBFF
		copy(buf[0:6], raw[0:6])
		buf[6], buf[7] = 'F', 'F'
	case 8:
		// #RRGGBBAA
		copy(buf[:], raw[:])
	default:
		panic(fmt.Sprintf("invalid hex length %d in %q", n, input))
	}

	// 3) manual hex→byte conversion
	nibble := func(c byte) uint8 {
		switch {
		case '0' <= c && c <= '9':
			return c - '0'
		case 'A' <= c && c <= 'F':
			return c - 'A' + 10
		default:
			panic(fmt.Sprintf("invalid hex digit %q", c))
		}
	}

	// combine two nibbles → one byte, then normalize to [0,1]
	makeFloat := func(hi, lo byte) float32 {
		v := nibble(hi)<<4 | nibble(lo)
		return float32(v) / 255.0
	}

	r := makeFloat(buf[0], buf[1])
	g := makeFloat(buf[2], buf[3])
	b := makeFloat(buf[4], buf[5])
	a := makeFloat(buf[6], buf[7])

	return Hex{
		hex:  string(buf[:]),
		rgba: RGBA{r, g, b, a},
	}
}
func (h Hex) ToRGBA() RGBA {
	return h.rgba
}
func (h Hex) String() string {
	return h.hex
}
