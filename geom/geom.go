package geom

import (
	"math"
)

// Direction ...
type Direction int

// Direction constants.
const (
	Left Direction = iota
	Right
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

// PerpindicularLocation ...
func PerpindicularLocation(buffer float64, dir Direction, b, a Marker) Location {
	d := b.Radius() + buffer
	r := Radians(b, a)
	bx, by := b.Coords()

	dMult := 1.0
	if dir == Left {
		dMult = -1.0
	}

	x := bx + ((d * math.Cos((math.Pi/2.0)+r)) * dMult)
	y := by + ((d * math.Sin((math.Pi/2.0)-r)) * dMult)

	return MakeLocation(x, y, 0)
}

// Obstacles demonstrates how the player might determine if the path
// between two enitities is clear
func Obstacles(ms []Marker, b, a Marker) bool {
	bx, by := b.Coords()
	ax, ay := a.Coords()
	dx := ax - bx
	dy := ay - by
	ptA := dx*dx + dy*dy + 1e-8
	crossterms := bx*bx - bx*ax + by*by - by*ay

	for _, po := range ms {
		pox, poy := po.Coords()

		if (pox == bx && poy == by) || (pox == ax && poy == ay) {
			continue
		}

		closestDistance := EdgeDistance(a, po)
		if closestDistance < 0 {
			return true
		}

		ptB := -2 * (crossterms + pox*dx + poy*dy)
		t := -ptB / (2 * ptA)

		if t <= 0 || t >= 1 {
			continue
		}

		closestX := bx + dx*t
		closestY := ax + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-pox, 2) * +math.Pow(closestY-poy, 2))

		if closestDistance <= po.Radius()+b.Radius()+1 {
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
