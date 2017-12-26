package fred

import (
	"math"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
)

// Logger describes the halitego logging behavior.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Fred ...
type Fred struct {
	l    Logger
	iniB ops.Board
}

// New ...
func New(l Logger, initialBoard ops.Board) *Fred {
	return &Fred{
		l:    l,
		iniB: initialBoard,
	}
}

// Command ...
func (bot *Fred) Command(b ops.Board, id int) ops.CommandMessengers {
	ss := b.Ships()[id]
	var ms ops.CommandMessengers

	for _, s := range ss {
		ms = append(ms, bot.messenger(b, s))
	}

	return ms
}

// messenger demonstrates how the player might direct their ships
// in achieving victory
func (bot *Fred) messenger(b ops.Board, s ops.Ship) ops.CommandMessenger {
	if s.SDStatus != ops.Undocked {
		return s.NoOp()
	}

	ps := planetsByProximity(b, s)
	for _, p := range ps {
		msg, err := s.Dock(p)
		if err == nil {
			return msg
		}

		derr, ok := err.(ops.DockingError)
		if ok && (derr.NoRights() || derr.NoPorts()) {
			continue
		}
		if ok && derr.NoJuncture() {
			return bot.navTo(b, geom.BufferedLocation(0, p, s), s)
		}

		panic("unexpected error while docking")
	}

	return s.NoOp()
}

// navTo demonstrates how the player might negotiate obsticles between
// a ship and its target
func (bot *Fred) navTo(b ops.Board, target geom.Marker, s ops.Ship) ops.CommandMessenger {
	ms := b.Markers()
	ob := geom.Obstacles(ms, target, s.Entity)

	if !ob {
		return s.Navigate(target)
	}

	tx, ty := target.Coords()
	cx, cy := s.Coords()

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
			ob1 := geom.Obstacles(ms, s.Entity, intermediateTarget)
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

	return s.Navigate(bestTarget)
}
