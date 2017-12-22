package fred

import (
	"math"
	"sort"

	"github.com/daved/halitego/ops"
)

type field struct {
	ops.Board
}

func makeField(b ops.Board) field {
	return field{b}
}

// ObstaclesBetween demonstrates how the player might determine if the path
// between two enitities is clear
func (f field) ObstaclesBetween(a, b ops.Marker) bool {
	x1, y1 := a.Coords()
	x2, y2 := b.Coords()
	dx := x2 - x1
	dy := y2 - y1
	ptA := dx*dx + dy*dy + 1e-8
	crossterms := x1*x1 - x1*x2 + y1*y1 - y1*y2

	var es []ops.Marker
	for _, v := range f.Planets() {
		es = append(es, v)
	}
	for _, v := range f.Ships() {
		for _, y := range v {
			es = append(es, y)
		}
	}

	for _, e := range es {
		x0, y0 := e.Coords()
		if x0 == x1 || x0 == x2 {
			continue
		}

		closestDistance := e.Distance(b)
		if closestDistance < e.Radius()+1 {
			return true
		}

		ptB := -2 * (crossterms + x0*dx + y0*dy)
		t := -ptB / (2 * ptA)

		if t <= 0 || t >= 1 {
			continue
		}

		closestX := x1 + dx*t
		closestY := x2 + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-x0, 2) * +math.Pow(closestY-y0, 2))

		if closestDistance <= e.Radius()+a.Radius()+1 {
			return true
		}
	}
	return false
}

// PlanetsByProximity orders all planets based on their proximity
// to a given ship from nearest for farthest
func (f field) PlanetsByProximity(l ops.Locator) []ops.Planet {
	pscs := makePlanetScans(f.Planets(), l)
	pscs.sortByDistance()

	return pscs.planets()
}

type planetScan struct {
	ops.Planet

	dist float64
}

type planetScans []planetScan

func makePlanetScans(ps []ops.Planet, l ops.Locator) planetScans {
	var pscs []planetScan
	for _, p := range ps {
		psc := planetScan{
			Planet: p,
			dist:   p.Distance(l),
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
