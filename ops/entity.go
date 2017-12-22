package ops

import "math"

// Locator describes types which are able to show their coordinates.
type Locator interface {
	Coords() (float64, float64)
}

// Sizer describes types which are able to show their size.
type Sizer interface {
	Sweep() float64
	Width() float64
}

// Marker describes types which are able to be located and sized.
type Marker interface {
	Locator
	Sizer
}

// Entity represents common attributes shared by items in a game map.
type Entity struct {
	ID     int
	Owner  int
	X      float64
	Y      float64
	Radius float64
	Health float64
}

// Coords returns the current x and y coordinates.
func (e Entity) Coords() (float64, float64) {
	return e.X, e.Y
}

// Sweep returns the current radius.
func (e Entity) Sweep() float64 {
	return e.Radius
}

// Width returns the current diameter.
func (e Entity) Width() float64 {
	return e.Radius * 2
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
func Nearest(min float64, b, a Marker) Entity {
	dist := Distance(a, b) - b.Sweep() - min
	angle := Radians(b, a)

	bx, by := b.Coords()
	x := bx + dist*math.Cos(angle)
	y := by + dist*math.Sin(angle)

	return Entity{
		X:      x,
		Y:      y,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		ID:     -1,
	}
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
