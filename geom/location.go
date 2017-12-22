package geom

// Location ...
type Location struct {
	x, y   float64
	radius float64
}

// MakeLocation ...
func MakeLocation(x, y, radius float64) Location {
	return Location{
		x:      x,
		y:      y,
		radius: radius,
	}
}

// Coords returns the current x and y coordinates.
func (l Location) Coords() (float64, float64) {
	return l.x, l.y
}

// Radius returns the current radius.
func (l Location) Radius() float64 {
	return l.radius
}
