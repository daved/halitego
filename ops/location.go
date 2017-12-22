package ops

// Location ...
type Location struct {
	Entity
}

// MakeLocation ...
func MakeLocation(x, y float64) Location {
	return Location{
		Entity{
			x:      x,
			y:      y,
			radius: 0,
			health: 0,
			owner:  0,
			id:     -1,
		},
	}
}
