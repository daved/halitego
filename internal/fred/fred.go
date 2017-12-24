package fred

import (
	"github.com/daved/halitego/ops"
)

// Fred ...
type Fred struct {
	id int
}

// New ...
func New(id int) *Fred {
	return &Fred{
		id: id,
	}
}

// Command ...
func (bot *Fred) Command(b ops.Board) ops.CommandMessengers {
	if bot.id >= len(b.Ships()) {
		return nil
	}

	ss := b.Ships()[bot.id]
	var ms ops.CommandMessengers

	for _, s := range ss {
		c := makeFACraft(s)
		ms = append(ms, c.messenger(b))
	}

	return ms
}
