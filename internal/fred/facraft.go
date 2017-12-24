package fred

import (
	"math"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
)

type faCraft struct {
	ops.Ship
}

func makeFACraft(s ops.Ship) faCraft {
	return faCraft{s}
}

// order demonstrates how the player might negotiate obsticles between
// a ship and its target
func (c faCraft) order(target geom.Marker, b ops.Board) ops.CommandMessenger {
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

// messenger demonstrates how the player might direct their ships
// in achieving victory
func (c faCraft) messenger(b ops.Board) ops.CommandMessenger {
	if c.SDStatus != ops.Undocked {
		return c.NoOp()
	}

	ps := planetsByProximity(b, c)

	for _, p := range ps {
		if (p.Owned == 0 || p.Owner() == c.Owner()) && p.DockedCt < p.PortCt && p.ID()%2 == c.ID()%2 {
			if msg, err := c.Dock(p); err == nil {
				return msg
			}

			return c.order(geom.Nearest(3, p, c), b)
		}
	}

	return c.NoOp()
}
