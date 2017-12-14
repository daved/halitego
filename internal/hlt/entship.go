package hlt

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

// IntToShipStatus converts an int to a ShipStatus.
func IntToShipStatus(i int) ShipStatus {
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
	id, _ := strconv.Atoi(tokens[0])
	x, _ := strconv.ParseFloat(tokens[1], 64)
	y, _ := strconv.ParseFloat(tokens[2], 64)
	health, _ := strconv.ParseFloat(tokens[3], 64)
	velX, _ := strconv.ParseFloat(tokens[4], 64)
	velY, _ := strconv.ParseFloat(tokens[5], 64)
	status, _ := strconv.Atoi(tokens[6])
	planetID, _ := strconv.Atoi(tokens[7])
	docking, _ := strconv.ParseFloat(tokens[8], 64)
	cooldown, _ := strconv.ParseFloat(tokens[9], 64)

	s := Ship{
		Entity: Entity{
			x:      x,
			y:      y,
			radius: .5,
			health: health,
			owner:  playerID,
			id:     id,
		},
		VelX:     velX,
		VelY:     velY,
		PlanetID: planetID,
		Status:   IntToShipStatus(status),
		Docking:  docking,
		Cooldown: cooldown,
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
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.id), strconv.Itoa(int(magnitude)), strconv.Itoa(boundedAngle))
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (ship Ship) Dock(planet Planet) (string, error) {
	isClose := Distance(ship, planet) <= (ship.radius + planet.radius + 4)

	if !isClose {
		return "", fmt.Errorf("cannot dock")
	}

	return fmt.Sprintf("d %d %d", ship.id, planet.id), nil
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship Ship) Undock() string {
	return fmt.Sprintf("u %s", strconv.Itoa(ship.id))
}

// NavigateBasic demonstrates how the player might move ships through space
func (ship Ship) NavigateBasic(target Entity, gameMap Map) string {
	dist := Distance(ship, target)
	safeDistance := dist - ship.Entity.radius - target.radius - .1

	angle := DegreesTo(ship, target)
	speed := 7.0
	if dist < 10 {
		speed = 3.0
	}

	speed = math.Min(speed, safeDistance)
	return ship.Thrust(speed, angle)
}

// Navigate demonstrates how the player might negotiate obsticles between
// a ship and its target
func (ship Ship) Navigate(target Entity, gameMap Map) string {
	ob := gameMap.ObstaclesBetween(ship.Entity, target)

	if !ob {
		return ship.NavigateBasic(target, gameMap)
	}

	x0 := math.Min(ship.x, target.x)
	x2 := math.Max(ship.x, target.x)
	y0 := math.Min(ship.y, target.y)
	y2 := math.Max(ship.y, target.y)

	dx := (x2 - x0) / 5
	dy := (y2 - y0) / 5
	bestdist := 1000.0
	bestTarget := target

	for x1 := x0; x1 <= x2; x1 += dx {
		for y1 := y0; y1 <= y2; y1 += dy {
			intermediateTarget := Entity{
				x:      x1,
				y:      y1,
				radius: 0,
				health: 0,
				owner:  0,
				id:     -1,
			}
			ob1 := gameMap.ObstaclesBetween(ship.Entity, intermediateTarget)
			if !ob1 {
				ob2 := gameMap.ObstaclesBetween(intermediateTarget, target)
				if !ob2 {
					totdist := math.Sqrt(math.Pow(x1-x0, 2)+math.Pow(y1-y0, 2)) + math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2))
					if totdist < bestdist {
						bestdist = totdist
						bestTarget = intermediateTarget

					}
				}
			}
		}
	}

	return ship.NavigateBasic(bestTarget, gameMap)
}
