package fred

import (
	"sort"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
)

// planetsByProximity orders all planets based on their proximity
// to a given ship from nearest for farthest
func planetsByProximity(b ops.Board, l geom.Marker) []ops.Planet {
	pscs := makePlanetScans(b.Planets(), l)
	pscs.sortByDistance()

	return pscs.planets()
}

type planetScan struct {
	ops.Planet

	dist float64
}

type planetScans []planetScan

func makePlanetScans(ps []ops.Planet, l geom.Marker) planetScans {
	var pscs []planetScan
	for _, p := range ps {
		psc := planetScan{
			Planet: p,
			dist:   geom.EdgeDistance(l, p),
		}
		pscs = append(pscs, psc)
	}

	return planetScans(pscs)
}

func (ss planetScans) sortByDistance() {
	sort.Sort(planetScansDistanceSort(ss))
}

func (ss planetScans) planets() []ops.Planet {
	var ps []ops.Planet
	for _, s := range ss {
		ps = append(ps, s.Planet)
	}

	return ps
}

type planetScansDistanceSort planetScans

func (s planetScansDistanceSort) Len() int {
	return len(s)
}

func (s planetScansDistanceSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s planetScansDistanceSort) Less(i, j int) bool {
	return s[i].dist < s[j].dist
}
