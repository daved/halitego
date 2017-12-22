package ops

import "github.com/daved/halitego/geom"

// Planet object from which Halite is mined
type Planet struct {
	Entity
	PortCt   float64
	DockedCt float64
	ProdRate float64
	Rsrcs    float64
	ShipIDs  []int
	Owned    float64
}

// makePlanet from a slice of game state tokens
func makePlanet(tokens []string) (Planet, []string) {
	p := Planet{
		Entity: Entity{
			Location: geom.MakeLocation(
				readTokenFloat(tokens, 1),
				readTokenFloat(tokens, 2),
				readTokenFloat(tokens, 4),
			),
			id:     readTokenInt(tokens, 0),
			health: readTokenFloat(tokens, 3),
			owner:  readTokenInt(tokens, 9),
		},
		PortCt:   readTokenFloat(tokens, 5),
		DockedCt: readTokenFloat(tokens, 10),
		ProdRate: readTokenFloat(tokens, 6),
		Rsrcs:    readTokenFloat(tokens, 7),
		Owned:    readTokenFloat(tokens, 8),
	}

	shipCt := int(p.DockedCt)

	for i := 0; i < shipCt; i++ {
		shipID := readTokenInt(tokens, 11+i)

		p.ShipIDs = append(p.ShipIDs, shipID)
	}

	return p, tokens[11+shipCt:]
}
