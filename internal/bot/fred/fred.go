package fred

import (
	"math"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
)

// Fred ...
type Fred struct{}

// New ...
func New() *Fred {
	return &Fred{}
}

// Command ...
func (bot *Fred) Command(b ops.Board, id int) ops.CommandMessengers {
	if id >= len(b.Ships()) {
		return nil
	}

	ss := b.Ships()[id]
	var ms ops.CommandMessengers

	for _, s := range ss {
		ms = append(ms, messenger(b, s))
	}

	return ms
}

// messenger demonstrates how the player might direct their ships
// in achieving victory
func messenger(b ops.Board, c ops.Ship) ops.CommandMessenger {
	if c.SDStatus != ops.Undocked {
		return c.NoOp()
	}

	ps := planetsByProximity(b, c)

	for _, p := range ps {
		if (p.Owned == 0 || p.Owner() == c.Owner()) && p.DockedCt < p.PortCt && p.ID()%2 == c.ID()%2 {
			if msg, err := c.Dock(p); err == nil {
				return msg
			}

			return navTo(b, geom.Nearest(3, p, c), c)
		}
	}

	return c.NoOp()
}

// navTo demonstrates how the player might negotiate obsticles between
// a ship and its target
func navTo(b ops.Board, target geom.Marker, c ops.Ship) ops.CommandMessenger {
	ms := b.Markers()
	ob := geom.Obstacles(ms, c.Entity, target)

	if !ob {
		return c.Navigate(target, b)
	}

	tx, ty := target.Coords()
	cx, cy := c.Coords()

	x0 := math.Min(cx, tx)
	x2 := math.Max(cx, tx)
	y0 := math.Min(cy, ty)
	y2 := math.Max(cy, ty)

	dx := (x2 - x0) / 5
	dy := (y2 - y0) / 5
	bestdist := 1000.0
	bestTarget := target

	for x1 := x0; x1 <= x2; x1 += dx {
		for y1 := y0; y1 <= y2; y1 += dy {
			intermediateTarget := geom.MakeLocation(x1, x2, 0)
			ob1 := geom.Obstacles(ms, c.Entity, intermediateTarget)
			if !ob1 {
				ob2 := geom.Obstacles(ms, intermediateTarget, target)
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

	return c.Navigate(bestTarget, b)
}
