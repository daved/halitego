package hlt

import (
	"fmt"
	"math"
	"strconv"
)

// ShipStatus represents possible ship states.
type ShipStatus int

const (
	// ShipStatus states.
	Undocked ShipStatus = iota
	Docking
	Docked
	Undocking
)

// Entity captures spacial and ownership state for Planets and Ships
type Entity struct {
	id     int
	owner  int
	x      float64
	y      float64
	radius float64
	health float64
}

// Coords returns the current x and y coordinates.
func (e Entity) Coords() (float64, float64) {
	return e.x, e.y
}

// Sweep returns the current radius.
func (e Entity) Sweep() float64 {
	return e.radius
}

// Width returns the current diameter.
func (e Entity) Width() float64 {
	return e.radius * 2
}

// Planet object from which Halite is mined
type Planet struct {
	Entity
	PortCt        float64
	DockedCt      float64
	ProdRate      float64
	Rsrcs         float64
	DockedShipIDs []int
	DockedShips   []Ship
	Owned         float64
	Distance      float64
}

// MakePlanet from a slice of game state tokens
func MakePlanet(tokens []string) (Planet, []string) {
	id, _ := strconv.Atoi(tokens[0])
	x, _ := strconv.ParseFloat(tokens[1], 64)
	y, _ := strconv.ParseFloat(tokens[2], 64)
	health, _ := strconv.ParseFloat(tokens[3], 64)
	radius, _ := strconv.ParseFloat(tokens[4], 64)
	portCt, _ := strconv.ParseFloat(tokens[5], 64)
	prodRate, _ := strconv.ParseFloat(tokens[6], 64)
	rsrcs, _ := strconv.ParseFloat(tokens[7], 64)
	owned, _ := strconv.ParseFloat(tokens[8], 64)
	owner, _ := strconv.Atoi(tokens[9])
	dockedCt, _ := strconv.ParseFloat(tokens[10], 64)

	pEnt := Entity{
		x:      x,
		y:      y,
		radius: radius,
		health: health,
		owner:  owner,
		id:     id,
	}

	p := Planet{
		PortCt:        portCt,
		DockedCt:      dockedCt,
		ProdRate:      prodRate,
		Rsrcs:         rsrcs,
		DockedShipIDs: nil,
		DockedShips:   nil,
		Owned:         owned,
		Entity:        pEnt,
	}

	for i := 0; i < int(dockedCt); i++ {
		dockedShipID, _ := strconv.Atoi(tokens[11+i])
		p.DockedShipIDs = append(p.DockedShipIDs, dockedShipID)
	}

	return p, tokens[11+int(dockedCt):]
}

// Ship is a player controlled Entity made for the purpose of doing combat and mining Halite
type Ship struct {
	Entity
	VelX float64
	VelY float64

	PlanetID        int
	PlanetEnv       Planet
	DockingStatus   ShipStatus
	DockingProgress float64
	WeaponCooldown  float64
}

// MakeShip from a slice of game state tokens
func MakeShip(playerID int, tokens []string) (Ship, []string) {
	shipID, _ := strconv.Atoi(tokens[0])
	shipX, _ := strconv.ParseFloat(tokens[1], 64)
	shipY, _ := strconv.ParseFloat(tokens[2], 64)
	shipHealth, _ := strconv.ParseFloat(tokens[3], 64)
	shipVelX, _ := strconv.ParseFloat(tokens[4], 64)
	shipVelY, _ := strconv.ParseFloat(tokens[5], 64)
	shipDockingStatus, _ := strconv.Atoi(tokens[6])
	shipPlanetID, _ := strconv.Atoi(tokens[7])
	shipDockingProgress, _ := strconv.ParseFloat(tokens[8], 64)
	shipWeaponCooldown, _ := strconv.ParseFloat(tokens[9], 64)

	shipEntity := Entity{
		x:      shipX,
		y:      shipY,
		radius: .5,
		health: shipHealth,
		owner:  playerID,
		id:     shipID,
	}

	ship := Ship{
		PlanetID:        shipPlanetID,
		DockingStatus:   IntToDockingStatus(shipDockingStatus),
		DockingProgress: shipDockingProgress,
		WeaponCooldown:  shipWeaponCooldown,
		VelX:            shipVelX,
		VelY:            shipVelY,
		Entity:          shipEntity,
	}

	return ship, tokens[10:]
}

// IntToDockingStatus converts an int to a DockingStatus
func IntToDockingStatus(i int) ShipStatus {
	statuses := [4]ShipStatus{Undocked, Docking, Docked, Undocking}
	return statuses[i]
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
	distance := Distance(ship, target)
	safeDistance := distance - ship.Entity.radius - target.radius - .1

	angle := Direction(ship, target)
	speed := 7.0
	if distance < 10 {
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
