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
	velX     float64
	velY     float64
	planetID int
	sdStatus ShipDockingStatus
	docking  float64
	cooldown float64
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
		velX:     readTokenFloat(tokens, 4),
		velY:     readTokenFloat(tokens, 5),
		planetID: readTokenInt(tokens, 7),
		sdStatus: makeShipStatus(readTokenInt(tokens, 6)),
		docking:  readTokenFloat(tokens, 8),
		cooldown: readTokenFloat(tokens, 9),
	}

	return s, tokens[10:]
}

// DockingStatus ...
func (s Ship) DockingStatus() ShipDockingStatus {
	return s.sdStatus
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
		right: p.owned != 0 && p.Owner() != s.Owner(),
		ports: p.dockedCt >= p.portCt,
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
