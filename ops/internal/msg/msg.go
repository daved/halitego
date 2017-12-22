package msg

import "fmt"

// Messenger ...
type Messenger interface {
	Message() string
}

// Messengers ...
type Messengers []Messenger

// Message ...
func (ms Messengers) Message() string {
	m := ""
	for _, v := range ms {
		m += v.Message() + " "
	}

	if len(m) == 0 {
		return ""
	}

	return m[:len(m)-1]
}

// NoOp ...
type NoOp struct{}

// MakeNoOp ...
func MakeNoOp() NoOp {
	return NoOp{}
}

// Message ...
func (m NoOp) Message() string {
	return ""
}

// Thrust ...
type Thrust struct {
	id        int
	magnitude int
	direction int
}

// MakeThrust ...
func MakeThrust(id, magnitude, direction int) Thrust {
	return Thrust{id, magnitude, direction}
}

// Message ...
func (m Thrust) Message() string {
	return fmt.Sprintf("t %d %d %d", m.id, m.magnitude, m.direction)
}

// Dock ...
type Dock struct {
	id       int
	planetID int
}

// MakeDock ...
func MakeDock(id, planetID int) Dock {
	return Dock{id, planetID}
}

// Message ...
func (m Dock) Message() string {
	return fmt.Sprintf("d %d %d", m.id, m.planetID)
}

// Undock ...
type Undock struct {
	id int
}

// MakeUndock ...
func MakeUndock(id int) Undock {
	return Undock{id}
}

// Message ...
func (m Undock) Message() string {
	return fmt.Sprintf("u %d", m.id)
}
