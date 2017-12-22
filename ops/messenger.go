package ops

import "fmt"

// Messenger ...
type Messenger interface {
	Message() string
}

// NoOpMsg ...
type NoOpMsg struct{}

func makeNoOpMsg() NoOpMsg {
	return NoOpMsg{}
}

// Message ...
func (m NoOpMsg) Message() string {
	return ""
}

// ThrustMsg ...
type ThrustMsg struct {
	id        int
	magnitude int
	direction int
}

func makeThrustMsg(id, magnitude, direction int) ThrustMsg {
	return ThrustMsg{id, magnitude, direction}
}

// Message ...
func (m ThrustMsg) Message() string {
	return fmt.Sprintf("t %d %d %d", m.id, m.magnitude, m.direction)
}

// DockMsg ...
type DockMsg struct {
	id       int
	planetID int
}

func makeDockMsg(id, planetID int) DockMsg {
	return DockMsg{id, planetID}
}

// Message ...
func (m DockMsg) Message() string {
	return fmt.Sprintf("d %d %d", m.id, m.planetID)
}

// UndockMsg ...
type UndockMsg struct {
	id int
}

func makeUndockMsg(id int) UndockMsg {
	return UndockMsg{id}
}

// Message ...
func (m UndockMsg) Message() string {
	return fmt.Sprintf("u %d", m.id)
}

func messengersToSysMsg(ms []Messenger) string {
	m := ""
	for _, v := range ms {
		m += v.Message() + " "
	}

	if len(m) == 0 {
		return ""
	}

	return m[:len(m)-1]
}
