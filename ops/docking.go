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
	KnownError() bool
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

// IsSet ...
func (e *DockingErr) IsSet() bool {
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

	if s == "" {
		return "unknown"
	}

	return s[1 : len(s)-1]
}
