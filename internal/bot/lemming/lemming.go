package lemming

import (
	"math/rand"
	"time"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
)

// Logger describes the halitego logging behavior.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Lemming ...
type Lemming struct {
	l    Logger
	iniB ops.Board
	rng  *rand.Rand
}

// New ...
func New(l Logger, initialBoard ops.Board) *Lemming {
	return &Lemming{
		l:    l,
		iniB: initialBoard,
		rng:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// Command ...
func (bot *Lemming) Command(b ops.Board, id int) ops.CommandMessengers {
	ss := b.Ships()[id]
	var ms ops.CommandMessengers

	for _, s := range ss {
		ms = append(ms, bot.messenger(b, s))
	}

	return ms
}

// messenger demonstrates how the player might direct their ships
// in achieving victory
func (bot *Lemming) messenger(b ops.Board, s ops.Ship) ops.CommandMessenger {
	if s.DockingStatus() != ops.Undocked {
		return s.NoOp()
	}

	ps := ops.PlanetsByProximity(b, s)

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
			return bot.nav(b, geom.BufferedLocation(2, p, s), s)

		}
	}

	return s.NoOp()
}

// nav demonstrates how the player might negotiate obsticles between
// a ship and its target
func (bot *Lemming) nav(b ops.Board, target geom.Marker, s ops.Ship) ops.CommandMessenger {
	ms := b.Markers()
	ob := geom.Obstacles(ms, target, s)
	if !ob {
		return s.Navigate(target)
	}

	buf := float64(bot.rng.Intn(23) + 17)
	if time.Now().Nanosecond()%2 == 0 {
		buf *= -1
	}
	pl := geom.PerpindicularLocation(buf, target, s)

	return bot.nav(b, pl, s)
}
