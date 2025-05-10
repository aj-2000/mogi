// math/vectors.go
package math

import "math"

//
// ——————————————————————————————————————————————————————————————————————————————
// mutable float64 vectors
// ——————————————————————————————————————————————————————————————————————————————
//

type Vec2 struct{ X, Y float64 }

func NewVec2(x, y float64) *Vec2 { return &Vec2{x, y} }

func (v *Vec2) Add(u Vec2) *Vec2 {
	v.X += u.X
	v.Y += u.Y
	return v
}

func (v *Vec2) Sub(u Vec2) *Vec2 {
	v.X -= u.X
	v.Y -= u.Y
	return v
}

func (v *Vec2) Scale(s float64) *Vec2 {
	v.X *= s
	v.Y *= s
	return v
}

func (v *Vec2) Dot(u Vec2) float64 {
	return v.X*u.X + v.Y*u.Y
}

func (v *Vec2) Norm() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v *Vec2) Normalize() *Vec2 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
	}
	return v
}

type Vec3 struct{ X, Y, Z float64 }

func NewVec3(x, y, z float64) *Vec3 { return &Vec3{x, y, z} }

func (v *Vec3) Add(u Vec3) *Vec3 {
	v.X += u.X
	v.Y += u.Y
	v.Z += u.Z
	return v
}

func (v *Vec3) Sub(u Vec3) *Vec3 {
	v.X -= u.X
	v.Y -= u.Y
	v.Z -= u.Z
	return v
}

func (v *Vec3) Scale(s float64) *Vec3 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

func (v *Vec3) Dot(u Vec3) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v *Vec3) Cross(u Vec3) *Vec3 {
	x := v.Y*u.Z - v.Z*u.Y
	y := v.Z*u.X - v.X*u.Z
	z := v.X*u.Y - v.Y*u.X
	v.X, v.Y, v.Z = x, y, z
	return v
}

func (v *Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vec3) Normalize() *Vec3 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
		v.Z /= n
	}
	return v
}

type Vec4 struct{ X, Y, Z, W float64 }

func NewVec4(x, y, z, w float64) *Vec4 { return &Vec4{x, y, z, w} }

func (v *Vec4) Add(u Vec4) *Vec4 {
	v.X += u.X
	v.Y += u.Y
	v.Z += u.Z
	v.W += u.W
	return v
}

func (v *Vec4) Sub(u Vec4) *Vec4 {
	v.X -= u.X
	v.Y -= u.Y
	v.Z -= u.Z
	v.W -= u.W
	return v
}

func (v *Vec4) Scale(s float64) *Vec4 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	v.W *= s
	return v
}

func (v *Vec4) Dot(u Vec4) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z + v.W*u.W
}

func (v *Vec4) Norm() float64 {
	sum := v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
	return math.Sqrt(sum)
}

func (v *Vec4) Normalize() *Vec4 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
		v.Z /= n
		v.W /= n
	}
	return v
}

//
// ——————————————————————————————————————————————————————————————————————————————
// mutable float32 vectors
// ——————————————————————————————————————————————————————————————————————————————
//

type Vec2f32 struct{ X, Y float32 }

func NewVec2f32(x, y float32) *Vec2f32 { return &Vec2f32{x, y} }

func (v *Vec2f32) Add(u Vec2f32) *Vec2f32 {
	v.X += u.X
	v.Y += u.Y
	return v
}

func (v *Vec2f32) Sub(u Vec2f32) *Vec2f32 {
	v.X -= u.X
	v.Y -= u.Y
	return v
}

func (v *Vec2f32) Scale(s float32) *Vec2f32 {
	v.X *= s
	v.Y *= s
	return v
}

func (v *Vec2f32) Dot(u Vec2f32) float32 {
	return v.X*u.X + v.Y*u.Y
}

func (v *Vec2f32) Norm() float32 {
	return float32(math.Hypot(float64(v.X), float64(v.Y)))
}

func (v *Vec2f32) Normalize() *Vec2f32 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
	}
	return v
}

func (v *Vec2f32) Clone() *Vec2f32 {
	return &Vec2f32{X: v.X, Y: v.Y}
}

func (v *Vec2f32) ToVec2() Vec2 {
	return Vec2{float64(v.X), float64(v.Y)}
}

type Vec3f32 struct{ X, Y, Z float32 }

func NewVec3f32(x, y, z float32) *Vec3f32 { return &Vec3f32{x, y, z} }

func (v *Vec3f32) Add(u Vec3f32) *Vec3f32 {
	v.X += u.X
	v.Y += u.Y
	v.Z += u.Z
	return v
}

func (v *Vec3f32) Sub(u Vec3f32) *Vec3f32 {
	v.X -= u.X
	v.Y -= u.Y
	v.Z -= u.Z
	return v
}

func (v *Vec3f32) Scale(s float32) *Vec3f32 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

func (v *Vec3f32) Dot(u Vec3f32) float32 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v *Vec3f32) Cross(u Vec3f32) *Vec3f32 {
	x := v.Y*u.Z - v.Z*u.Y
	y := v.Z*u.X - v.X*u.Z
	z := v.X*u.Y - v.Y*u.X
	v.X, v.Y, v.Z = x, y, z
	return v
}

func (v *Vec3f32) Norm() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v *Vec3f32) Normalize() *Vec3f32 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
		v.Z /= n
	}
	return v
}

func (v *Vec3f32) Clone() *Vec3f32 {
	return &Vec3f32{X: v.X, Y: v.Y, Z: v.Z}
}

func (v *Vec3f32) ToVec3() Vec3 {
	return Vec3{float64(v.X), float64(v.Y), float64(v.Z)}
}

type Vec4f32 struct{ X, Y, Z, W float32 }

func NewVec4f32(x, y, z, w float32) *Vec4f32 { return &Vec4f32{x, y, z, w} }

func (v *Vec4f32) Add(u Vec4f32) *Vec4f32 {
	v.X += u.X
	v.Y += u.Y
	v.Z += u.Z
	v.W += u.W
	return v
}

func (v *Vec4f32) Sub(u Vec4f32) *Vec4f32 {
	v.X -= u.X
	v.Y -= u.Y
	v.Z -= u.Z
	v.W -= u.W
	return v
}

func (v *Vec4f32) Scale(s float32) *Vec4f32 {
	v.X *= s
	v.Y *= s
	v.Z *= s
	v.W *= s
	return v
}

func (v *Vec4f32) Dot(u Vec4f32) float32 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z + v.W*u.W
}

func (v *Vec4f32) Norm() float32 {
	sum := v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
	return float32(math.Sqrt(float64(sum)))
}

func (v *Vec4f32) Normalize() *Vec4f32 {
	if n := v.Norm(); n != 0 {
		v.X /= n
		v.Y /= n
		v.Z /= n
		v.W /= n
	}
	return v
}

func (v *Vec4f32) Clone() *Vec4f32 {
	return &Vec4f32{X: v.X, Y: v.Y, Z: v.Z, W: v.W}
}

func (v *Vec4f32) ToVec4() Vec4 {
	return Vec4{float64(v.X), float64(v.Y), float64(v.Z), float64(v.W)}
}
