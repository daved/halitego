package ops

import "math"

// Locator describes types which are able to show their coordinates.
type Locator interface {
	Coords() (float64, float64)
	Distance(Locator) float64
	Degrees(Locator) float64
	Radians(Locator) float64
}

// Sizer describes types which are able to show their size.
type Sizer interface {
	Radius() float64
	Diameter() float64
}

// Marker describes types which are able to be located and sized.
type Marker interface {
	Locator
	Sizer
}

// Entity represents common attributes shared by items in a game map.
type Entity struct {
	id     int
	owner  int
	x, y   float64
	radius float64
	health float64
}

// ID ...
func (e Entity) ID() int {
	return e.id
}

// Owner ...
func (e Entity) Owner() int {
	return e.owner
}

// Health ...
func (e Entity) Health() float64 {
	return e.health
}

// Coords returns the current x and y coordinates.
func (e Entity) Coords() (float64, float64) {
	return e.x, e.y
}

// Radius returns the current radius.
func (e Entity) Radius() float64 {
	return e.radius
}

// Diameter returns the current diameter.
func (e Entity) Diameter() float64 {
	return e.radius * 2
}

// Distance returns the Distance between two instances of Locator types.
func (e Entity) Distance(b Locator) float64 {
	bx, by := b.Coords()
	ax, ay := e.Coords()

	return distanceBetween(bx, by, ax, ay)
}

// Degrees returns the angle in degrees between two instances of Locator types.
func (e Entity) Degrees(b Locator) float64 {
	return radiansToDegrees(e.Radians(b))
}

// Radians returns the angle in radians between two instances of Locator types.
func (e Entity) Radians(b Locator) float64 {
	bx, by := b.Coords()
	ax, ay := e.Coords()

	return radiansBetween(bx, by, ax, ay)
}

// Nearest returns the closest point from Marker "a" to Marker "b" that is at
// least a distance of "min" from Marker "b".
func (e Entity) Nearest(min float64, b Marker) Entity {
	dist := e.Distance(b) - b.Radius() - min
	angle := e.Radians(b)

	bx, by := b.Coords()
	x := bx + dist*math.Cos(angle)
	y := by + dist*math.Sin(angle)

	return Entity{
		x:      x,
		y:      y,
		radius: 0,
		health: 0,
		owner:  -1,
		id:     -1,
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
