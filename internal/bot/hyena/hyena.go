package hyena

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

// Hyena ...
type Hyena struct {
	l    Logger
	iniB ops.Board
	rng  *rand.Rand
}

// New ...
func New(l Logger, initialBoard ops.Board) *Hyena {
	return &Hyena{
		l:    l,
		iniB: initialBoard,
		rng:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// Command ...
func (bot *Hyena) Command(b ops.Board, id int) ops.CommandMessengers {
	ss := b.Ships()[id]
	var ms ops.CommandMessengers

	for _, s := range ss {
		ms = append(ms, bot.messenger(b, s))
	}

	return ms
}

// messenger demonstrates how the player might direct their ships
// in achieving victory
func (bot *Hyena) messenger(b ops.Board, s ops.Ship) ops.CommandMessenger {
	if s.DockingStatus() != ops.Undocked {
		return s.NoOp()
	}

	ps := ops.PlanetsByProximity(b, s)
	striking := bot.allOwned(ps)

	for _, p := range ps {
		msg, err := s.Dock(p)
		if err == nil {
			return msg
		}

		if aMsg, ok := bot.altMsg(err, striking, bpsLoad{b, p, s}); ok {
			return aMsg
		}
	}

	return s.NoOp()
}

type bpsLoad struct {
	b ops.Board
	p ops.Planet
	s ops.Ship
}

func (bot *Hyena) altMsg(err error, striking bool, bps bpsLoad) (ops.CommandMessenger, bool) {
	derr, ok := err.(ops.DockingError)
	if !ok {
		return bps.s.NoOp(), true
	}

	if !striking && (derr.NoPorts() || derr.NoRights()) {
		return bps.s.NoOp(), false
	}

	if striking && !derr.NoRights() {
		return bps.s.NoOp(), false
	}

	if derr.NoJuncture() {
		return bot.nav(bps.b, geom.BufferedLocation(2, bps.p, bps.s), bps.s), true
	}

	return bps.s.NoOp(), true
}

// nav demonstrates how the player might negotiate obsticles between
// a ship and its target
func (bot *Hyena) nav(b ops.Board, target geom.Marker, s ops.Ship) ops.CommandMessenger {
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

func (bot *Hyena) allOwned(ps []ops.Planet) bool {
	for _, p := range ps {
		if !p.Owned() {
			return false
		}
	}

	return true
}
