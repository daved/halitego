package ops

// ShipDockingStatus represents possible ship docking states.
type ShipDockingStatus int

// ShipDockingStatus states.
const (
	Undocked ShipDockingStatus = iota
	Docking
	Docked
	Undocking
)

// DockingError ...
type DockingError interface {
	error
	NoJuncture() bool
	NoRights() bool
	NoPorts() bool
}

// DockingErr ...
type DockingErr struct {
	junct bool
	right bool
	ports bool
}

// Error ...
func (e *DockingErr) Error() string {
	return "cannot dock: " + e.reason()
}

// IsError ...
func (e *DockingErr) IsError() bool {
	return e.junct || e.right || e.ports
}

// NoJuncture ...
func (e *DockingErr) NoJuncture() bool {
	return e.junct
}

// NoRights ...
func (e *DockingErr) NoRights() bool {
	return e.right
}

// NoPorts ...
func (e *DockingErr) NoPorts() bool {
	return e.ports
}

func (e *DockingErr) reason() string {
	if !e.IsError() {
		return "unknown"
	}

	s := ""
	if e.junct {
		s += " no proximity,"
	}
	if e.right {
		s += " no permission,"
	}
	if e.ports {
		s += " no available port,"
	}

	return s[1 : len(s)-1]
}
