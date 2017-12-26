package geom

import (
	"math"
)

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

// CenterDistance returns the distance between two instances of Locator types.
func CenterDistance(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	d := distanceBetween(bx, by, ax, ay)

	return d
}

// EdgeDistance returns the distance between two instances of Locator types.
func EdgeDistance(b, a Marker) float64 {
	return CenterDistance(b, a) - b.Radius() - a.Radius()
}

// Degrees returns the angle in degrees between two instances of Locator types.
func Degrees(b, a Locator) float64 {
	return radiansToDegrees(Radians(b, a))
}

// BoundDegrees ...
func BoundDegrees(b, a Locator) int {
	d := Degrees(b, a)

	bd := int(math.Ceil(d - .5))
	if d > 0.0 {
		bd = int(math.Floor(d + .5))
	}

	return ((bd % 360) + 360) % 360
}

// Radians returns the angle in radians between two instances of Locator types.
func Radians(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	return radiansBetween(bx, by, ax, ay)
}

// BufferedLocation returns the closest point from Marker "a" to Marker "b" that is at
// least a distance of "min" from Marker "b".
func BufferedLocation(buffer float64, b, a Marker) Location {
	d := EdgeDistance(b, a) - buffer

	r := Radians(b, a)

	ax, ay := a.Coords()
	x := ax + (d * math.Cos(r))
	y := ay + (d * math.Sin(r))

	return MakeLocation(x, y, 0)
}

// Obstacles demonstrates how the player might determine if the path
// between two enitities is clear
func Obstacles(ms []Marker, b, a Marker) bool {
	x1, y1 := b.Coords()
	x2, y2 := a.Coords()
	dx := x2 - x1
	dy := y2 - y1
	ptA := dx*dx + dy*dy + 1e-8
	crossterms := x1*x1 - x1*x2 + y1*y1 - y1*y2

	for _, e := range ms {
		x0, y0 := e.Coords()
		if x0 == x1 || x0 == x2 {
			continue
		}

		closestDistance := EdgeDistance(a, e)
		if closestDistance < e.Radius()+1 {
			return true
		}

		ptB := -2 * (crossterms + x0*dx + y0*dy)
		t := -ptB / (2 * ptA)

		if t <= 0 || t >= 1 {
			continue
		}

		closestX := x1 + dx*t
		closestY := x2 + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-x0, 2) * +math.Pow(closestY-y0, 2))

		if closestDistance <= e.Radius()+b.Radius()+1 {
			return true
		}
	}

	return false
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
