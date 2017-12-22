package ops

import (
	"fmt"
	"math"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops/internal/msg"
)

// ShipDockingStatus represents possible ship docking states.
type ShipDockingStatus int

// ShipDockingStatus states.
const (
	Undocked ShipDockingStatus = iota
	Docking
	Docked
	Undocking
)

// makeShipStatus converts an int to a ShipStatus.
func makeShipStatus(i int) ShipDockingStatus {
	s := [4]ShipDockingStatus{Undocked, Docking, Docked, Undocking}

	return s[i]
}

// Ship represents ship state.
type Ship struct {
	Entity
	VelX     float64
	VelY     float64
	PlanetID int
	SDStatus ShipDockingStatus
	Docking  float64
	Cooldown float64
}

// makeShip from a slice of game state tokens
func makeShip(playerID int, tokens []string) (Ship, []string) {
	s := Ship{
		Entity: Entity{
			Location: geom.MakeLocation(
				readTokenFloat(tokens, 1),
				readTokenFloat(tokens, 2),
				0.5,
			),
			id:     readTokenInt(tokens, 0),
			health: readTokenFloat(tokens, 3),
			owner:  playerID,
		},
		VelX:     readTokenFloat(tokens, 4),
		VelY:     readTokenFloat(tokens, 5),
		PlanetID: readTokenInt(tokens, 7),
		SDStatus: makeShipStatus(readTokenInt(tokens, 6)),
		Docking:  readTokenFloat(tokens, 8),
		Cooldown: readTokenFloat(tokens, 9),
	}

	return s, tokens[10:]
}

// DockingStatus ...
func (ship Ship) DockingStatus() ShipDockingStatus {
	return ship.SDStatus
}

// NoOp ...
func (ship Ship) NoOp() msg.NoOp {
	return msg.MakeNoOp()
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (ship Ship) Dock(planet Planet) (msg.Dock, error) {
	isClose := geom.Distance(planet, ship) <= (ship.Radius() + planet.Radius() + 4)

	var err error
	if !isClose {
		err = fmt.Errorf("cannot dock")
	}

	return msg.MakeDock(ship.id, planet.id), err
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship Ship) Undock() msg.Undock {
	return msg.MakeUndock(ship.id)
}

// NavigateTo demonstrates how the player might move ships through space
func (ship Ship) NavigateTo(target geom.Marker, gameMap Board) msg.Thrust {
	dist := geom.Distance(target, ship)
	safeDistance := dist - ship.Radius() - target.Radius() - .1

	angle := geom.Degrees(target, ship)
	speed := 7.0
	if dist < 10 {
		speed = 3.0
	}

	speed = math.Min(speed, safeDistance)
	return ship.thrust(speed, angle)
}

// thrust generates a string describing the ship's intension to move during the current turn
func (ship Ship) thrust(magnitude float64, angle float64) msg.Thrust {
	bounded := int(math.Ceil(angle - .5))

	if angle > 0.0 {
		bounded = int(math.Floor(angle + .5))
	}

	bounded = ((bounded % 360) + 360) % 360

	return msg.MakeThrust(ship.id, int(magnitude), bounded)
}
