package ops

import (
	"fmt"
	"math"
)

// ShipStatus represents possible ship docking states.
type ShipStatus int

// ShipStatus states.
const (
	Undocked ShipStatus = iota
	Docking
	Docked
	Undocking
)

// MakeShipStatus converts an int to a ShipStatus.
func MakeShipStatus(i int) ShipStatus {
	s := [4]ShipStatus{Undocked, Docking, Docked, Undocking}

	return s[i]
}

// Ship represents ship state.
type Ship struct {
	Entity
	VelX     float64
	VelY     float64
	PlanetID int
	Status   ShipStatus
	Docking  float64
	Cooldown float64
}

// MakeShip from a slice of game state tokens
func MakeShip(playerID int, tokens []string) (Ship, []string) {
	s := Ship{
		Entity: Entity{
			ID:     readTokenInt(tokens, 0),
			X:      readTokenFloat(tokens, 1),
			Y:      readTokenFloat(tokens, 2),
			Radius: 0.5,
			Health: readTokenFloat(tokens, 3),
			Owner:  playerID,
		},
		VelX:     readTokenFloat(tokens, 4),
		VelY:     readTokenFloat(tokens, 5),
		PlanetID: readTokenInt(tokens, 7),
		Status:   MakeShipStatus(readTokenInt(tokens, 6)),
		Docking:  readTokenFloat(tokens, 8),
		Cooldown: readTokenFloat(tokens, 9),
	}

	return s, tokens[10:]
}

// NoOp ...
func (ship Ship) NoOp() NoOpMsg {
	return makeNoOpMsg()
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (ship Ship) Dock(planet Planet) (DockMsg, error) {
	isClose := Distance(ship, planet) <= (ship.Radius + planet.Radius + 4)

	var err error
	if !isClose {
		err = fmt.Errorf("cannot dock")
	}

	return makeDockMsg(ship.ID, planet.ID), err
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship Ship) Undock() UndockMsg {
	return makeUndockMsg(ship.ID)
}

// NavigateTo demonstrates how the player might move ships through space
func (ship Ship) NavigateTo(target Entity, gameMap Board) ThrustMsg {
	dist := Distance(ship, target)
	safeDistance := dist - ship.Entity.Radius - target.Radius - .1

	angle := Degrees(ship, target)
	speed := 7.0
	if dist < 10 {
		speed = 3.0
	}

	speed = math.Min(speed, safeDistance)
	return ship.thrust(speed, angle)
}

// thrust generates a string describing the ship's intension to move during the current turn
func (ship Ship) thrust(magnitude float64, angle float64) ThrustMsg {
	bounded := int(math.Ceil(angle - .5))

	if angle > 0.0 {
		bounded = int(math.Floor(angle + .5))
	}

	bounded = ((bounded % 360) + 360) % 360

	return makeThrustMsg(ship.ID, int(magnitude), bounded)
}
