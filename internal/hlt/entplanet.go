package hlt

import "strconv"

// Planet object from which Halite is mined
type Planet struct {
	Entity
	PortCt   float64
	DockedCt float64
	ProdRate float64
	Rsrcs    float64
	ShipsIDS []int
	Ships    []Ship
	Owned    float64
	Distance float64
}

// MakePlanet from a slice of game state tokens
func MakePlanet(tokens []string) (Planet, []string) {
	id, _ := strconv.Atoi(tokens[0])
	x, _ := strconv.ParseFloat(tokens[1], 64)
	y, _ := strconv.ParseFloat(tokens[2], 64)
	health, _ := strconv.ParseFloat(tokens[3], 64)
	radius, _ := strconv.ParseFloat(tokens[4], 64)
	portCt, _ := strconv.ParseFloat(tokens[5], 64)
	prodRate, _ := strconv.ParseFloat(tokens[6], 64)
	rsrcs, _ := strconv.ParseFloat(tokens[7], 64)
	owned, _ := strconv.ParseFloat(tokens[8], 64)
	owner, _ := strconv.Atoi(tokens[9])
	dockedCt, _ := strconv.ParseFloat(tokens[10], 64)

	pEnt := Entity{
		x:      x,
		y:      y,
		radius: radius,
		health: health,
		owner:  owner,
		id:     id,
	}

	p := Planet{
		PortCt:   portCt,
		DockedCt: dockedCt,
		ProdRate: prodRate,
		Rsrcs:    rsrcs,
		ShipsIDS: nil,
		Ships:    nil,
		Owned:    owned,
		Entity:   pEnt,
	}

	for i := 0; i < int(dockedCt); i++ {
		dockedShipID, _ := strconv.Atoi(tokens[11+i])
		p.ShipsIDS = append(p.ShipsIDS, dockedShipID)
	}

	return p, tokens[11+int(dockedCt):]
}
