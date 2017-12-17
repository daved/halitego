package ops

import (
	"fmt"
	"math"
	"strconv"
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
		Status:   IntToShipStatus(readTokenInt(tokens, 6)),
		Docking:  readTokenFloat(tokens, 8),
		Cooldown: readTokenFloat(tokens, 9),
	}

	return s, tokens[10:]
}

// Thrust generates a string describing the ship's intension to move during the current turn
func (ship Ship) Thrust(magnitude float64, angle float64) string {
	var boundedAngle int
	if angle > 0.0 {
		boundedAngle = int(math.Floor(angle + .5))
	} else {
		boundedAngle = int(math.Ceil(angle - .5))
	}
	boundedAngle = ((boundedAngle % 360) + 360) % 360
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.ID), strconv.Itoa(int(magnitude)), strconv.Itoa(boundedAngle))
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (ship Ship) Dock(planet Planet) (string, error) {
	isClose := Distance(ship, planet) <= (ship.Radius + planet.Radius + 4)

	if !isClose {
		return "", fmt.Errorf("cannot dock")
	}

	return fmt.Sprintf("d %d %d", ship.ID, planet.ID), nil
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship Ship) Undock() string {
	return fmt.Sprintf("u %s", strconv.Itoa(ship.ID))
}

// NavigateTo demonstrates how the player might move ships through space
func (ship Ship) NavigateTo(target Entity, gameMap Board) string {
	dist := Distance(ship, target)
	safeDistance := dist - ship.Entity.Radius - target.Radius - .1

	angle := DegreesTo(ship, target)
	speed := 7.0
	if dist < 10 {
		speed = 3.0
	}

	speed = math.Min(speed, safeDistance)
	return ship.Thrust(speed, angle)
}

// IntToShipStatus converts an int to a ShipStatus.
func IntToShipStatus(i int) ShipStatus {
	s := [4]ShipStatus{Undocked, Docking, Docked, Undocking}

	return s[i]
}
