package hlt

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

// Distance returns the distance between two instances of Locator.
func Distance(a, b Locator) float64 {
	ax, ay := a.Coords()
	bx, by := b.Coords()

	return distanceBetween(ax, ay, bx, by)
}

func distanceBetween(ax, ay, bx, by float64) float64 {
	dx := bx - ax
	dy := by - ay

	return math.Sqrt(dx*dx + dy*dy)
}

// Direction returns an angle in degrees to the target.
func Direction(a, b Locator) float64 {
	return RadToDeg(DirectionRadians(a, b))
}

// DirectionRadians returns an angle in radians to the target
func DirectionRadians(a, b Locator) float64 {
	ax, ay := a.Coords()
	bx, by := b.Coords()

	return radiansToFacing(ax, ay, bx, by)
}

func radiansToFacing(ax, ay, bx, by float64) float64 {
	dx := bx - ax
	dy := by - ay

	return math.Atan2(dy, dx)
}

// NearestPoint returns the closest point that is at least minDistance from the target
func NearestPoint(a, b Marker, minDistance float64) Entity {
	dist := Distance(a, b) - b.Sweep() - minDistance
	angle := DirectionRadians(b, a)

	bx, by := b.Coords()
	x := bx + dist*math.Cos(angle)
	y := by + dist*math.Sin(angle)

	return Entity{
		X:      x,
		Y:      y,
		radius: 0,
		Health: 0,
		Owner:  -1,
		ID:     -1,
	}
}
