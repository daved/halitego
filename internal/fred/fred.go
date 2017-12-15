package fred

import (
	"fmt"

	"github.com/daved/halitego/ops"
)

// Fred ...
type Fred struct {
	id int
}

// NewFred ...
func New(id int) *Fred {
	return &Fred{
		id: id,
	}
}

type tmpStringer struct {
	s string
}

func (s *tmpStringer) String() string {
	return s.s
}

// Cmds ...
func (f *Fred) Cmds(b ops.Board) []fmt.Stringer {
	ships := b.Ships()[f.id]

	var cmds []fmt.Stringer
	for k := range ships {
		s := ships[k]

		if s.Status == ops.Undocked {
			cmds = append(cmds, &tmpStringer{f.StrategyBasicBot(s, b)})
		}
	}

	return cmds
}

// StrategyBasicBot demonstrates how the player might direct their ships
// in achieving victory
func (f *Fred) StrategyBasicBot(ship ops.Ship, gameMap ops.Board) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.Owner == f.id) && planet.DockedCt < planet.PortCt && planet.ID%2 == ship.ID%2 {
			if msg, err := ship.Dock(planet); err == nil {
				return msg
			}

			return ship.Navigate(ops.Nearest(ship, planet, 3), gameMap)
		}
	}

	return ""
}
