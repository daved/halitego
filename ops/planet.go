package ops

import "github.com/daved/halitego/geom"

// Planet object from which Halite is mined
type Planet struct {
	Entity
	portCt   float64
	dockedCt float64
	prodRate float64
	rsrcs    float64
	shipIDs  []int
	owned    float64
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
		portCt:   readTokenFloat(tokens, 5),
		dockedCt: readTokenFloat(tokens, 10),
		prodRate: readTokenFloat(tokens, 6),
		rsrcs:    readTokenFloat(tokens, 7),
		owned:    readTokenFloat(tokens, 8),
	}

	shipCt := int(p.dockedCt)

	for i := 0; i < shipCt; i++ {
		shipID := readTokenInt(tokens, 11+i)

		p.shipIDs = append(p.shipIDs, shipID)
	}

	return p, tokens[11+shipCt:]
}

func (p Planet) Owned() bool {
	return p.owned > 0
}
