package ops

import "github.com/daved/halitego/geom"

// Entity represents common attributes shared by items in a game map.
type Entity struct {
	id     int
	owner  int
	health float64

	geom.Location
}

// MakeEntity ...
func MakeEntity(x, y, radius, health float64, id, owner int) Entity {
	return Entity{
		id:       id,
		owner:    owner,
		health:   health,
		Location: geom.MakeLocation(x, y, radius),
	}
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
