package ops

import (
	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops/internal/msg"
)

// makeShipStatus converts an int to a ShipStatus.
func makeShipStatus(i int) ShipDockingStatus {
	ss := [4]ShipDockingStatus{Undocked, Docking, Docked, Undocking}

	return ss[i]
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
func (s Ship) DockingStatus() ShipDockingStatus {
	return s.SDStatus
}

// NoOp ...
func (s Ship) NoOp() msg.NoOp {
	return msg.MakeNoOp()
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (s Ship) Dock(p Planet) (msg.Dock, error) {
	msg := msg.MakeDock(s.id, p.id)
	err := &DockingErr{
		junct: geom.CenterDistance(p, s)-p.Radius()-4.0 > 0,
		right: p.Owned != 0 && p.Owner() != s.Owner(),
		ports: p.DockedCt >= p.PortCt,
	}

	if err.IsError() {
		return msg, err
	}

	return msg, nil
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (s Ship) Undock() msg.Undock {
	return msg.MakeUndock(s.id)
}

// Navigate demonstrates how the player might move ships through space
func (s Ship) Navigate(l geom.Locator) msg.Thrust {
	sp := 7
	a := geom.BoundDegrees(l, s)

	d := geom.CenterDistance(l, s)
	id := s.id
	if d < 7 {
		sp = int(d)
	}

	return msg.MakeThrust(id, sp, a)
}
