package geom

import "math"

// Locator describes types which are able to show their coordinates.
type Locator interface {
	Coords() (float64, float64)
}

// Sizer describes types which are able to show their size.
type Sizer interface {
	Radius() float64
}

// Marker describes types which are able to be located and sized.
type Marker interface {
	Locator
	Sizer
}

// Distance returns the Distance between two instances of Locator types.
func Distance(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	return distanceBetween(bx, by, ax, ay)
}

// Degrees returns the angle in degrees between two instances of Locator types.
func Degrees(b, a Locator) float64 {
	return radiansToDegrees(Radians(b, a))
}

// Radians returns the angle in radians between two instances of Locator types.
func Radians(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	return radiansBetween(bx, by, ax, ay)
}

// Nearest returns the closest point from Marker "a" to Marker "b" that is at
// least a distance of "min" from Marker "b".
func Nearest(min float64, b, a Marker) Location {
	dist := Distance(b, a) - b.Radius() - min
	angle := Radians(b, a)

	bx, by := b.Coords()
	x := bx + dist*math.Cos(angle)
	y := by + dist*math.Sin(angle)

	return MakeLocation(x, y, 0)
}

func distanceBetween(bx, by, ax, ay float64) float64 {
	dx := ax - bx
	dy := ay - by

	return math.Sqrt(dx*dx + dy*dy)
}

func radiansBetween(bx, by, ax, ay float64) float64 {
	dx := bx - ax
	dy := by - ay

	return math.Atan2(dy, dx)
}

func radiansToDegrees(r float64) float64 {
	return r / math.Pi * 180
}
