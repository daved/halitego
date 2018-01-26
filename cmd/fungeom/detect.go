package main

import (
	"math"

	"github.com/daved/halitego/geom"
)

func (c *graphicContext) obstacles(ms []geom.Marker, b, a geom.Marker) bool {
	pfl, pfr, pbr, pbl := pathPolygon(b, a)
	pmaxX, pminX, pmaxY, pminY := minmax(pfl, pfr, pbr, pbl)

	for _, po := range ms {
		if !potentialObstacle(po, pmaxX, pminX, pmaxY, pminY) {
			continue
		}

		if c.obstacle(po, pfl, pfr, pbr, pbl) {
			return true
		}
	}

	cs := makeCoordsFromGeomLocators(pfl, pfr, pbr, pbl) // can skip
	c.addDrawers(makePoly(cs, purt))                     // can skip

	return false
}

func pathPolygon(b, a geom.Marker) (h, i, j, k geom.Locator) {
	pfl := geom.PerpindicularLocation(-b.Radius()+a.Radius(), geom.Left, b, a)
	pfr := geom.PerpindicularLocation(-b.Radius()+a.Radius(), geom.Right, b, a)
	pbr := geom.PerpindicularLocation(0, geom.Left, a, b)
	pbl := geom.PerpindicularLocation(0, geom.Right, a, b)

	return pfl, pfr, pbr, pbl
}

func minmax(coords ...geom.Locator) (maxX, minX, maxY, minY float64) {
	maxX, minX, maxY, minY = 0.0, 1000000000.0, 0.0, 1000000000.0
	for _, c := range coords {
		x, y := c.Coords()
		maxX = math.Max(x, maxX)
		minX = math.Min(x, minX)
		maxY = math.Max(y, maxY)
		minY = math.Min(y, minY)
	}

	return maxX, minX, maxY, minY
}

func potentialObstacle(ob geom.Marker, maxX, minX, maxY, minY float64) bool {
	x, y := ob.Coords()
	r := ob.Radius()
	maxX, minX, maxY, minY = maxX+r, minX-r, maxY+r, minY+r

	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func (c *graphicContext) obstacle(ob geom.Marker, coords ...geom.Locator) bool {
	if c.inPolygon(false, ob, coords...) {
		return true
	}

	for _, pi := range potentialIntrusions(ob, coords...) {
		if c.inPolygon(true, pi, coords...) {
			return true
		}
	}

	return false
}

func (c *graphicContext) inPolygon(t bool, m geom.Marker, coords ...geom.Locator) bool {
	if t {
		destMarker := makeEntityFromMarker(m, redt)
		c.addDrawers(destMarker)
	}

	if len(coords) < 3 {
		return false
	}

	ct := 0
	x, y := m.Coords()
	a := coords[0]
	b := coords[1]

	for k, cz := range coords {
		if m.Radius() > 0 && geom.CenterDistance(m, cz) <= m.Radius() {
			return true
		}

		if k < 2 {
			continue
		}

		if inTriangle(x, y, a, b, cz) {
			ct++
		}

		b = cz
	}

	return ct%2 != 0
}

func potentialIntrusions(gm geom.Marker, coords ...geom.Locator) []geom.Marker {
	var ms []geom.Marker
	if len(coords) < 2 {
		return ms
	}

	p1x, p1y := coords[len(coords)-1].Coords()
	p3x, p3y := gm.Coords()

	for _, c := range coords {
		p2x, p2y := c.Coords()

		m := (p2y - p1y) / (p2x - p1x)
		b1 := m*-p1x + p1y

		// intersecting point
		b2 := 1/m*p3x + p3y
		iy := m*(b2-b1)/(m+1/m) + b1
		ix := (iy - b1) / m
		pi := geom.MakeLocation(ix, iy, 0)

		p1x, p1y = p2x, p2y

		if geom.CenterDistance(gm, pi) > gm.Radius() {
			continue
		}

		ms = append(ms, pi)
		ms = append(ms, geom.MakeLocation(ix+.0001, iy+.0001, 0)) // rough fix for overlap errors
		ms = append(ms, geom.MakeLocation(ix-.0001, iy-.0001, 0)) // rough fix for overlap errors
	}

	return ms
}

// https://math.stackexchange.com/questions/51326
func inTriangle(x, y float64, a, b, c geom.Locator) bool {
	ax, ay := a.Coords()
	bx, by := b.Coords()
	cx, cy := c.Coords()

	x, y = shift(x, y, ax, ay)
	bx, by = shift(bx, by, ax, ay)
	cx, cy = shift(cx, cy, ax, ay)

	d := bx*cy - cx*by
	wa := (x*(by-cy) + y*(cx-bx) + bx*cy - cx*by) / d
	wb := (x*cy - y*cx) / d
	wc := (y*bx - x*by) / d

	return ge0le1(wa) && ge0le1(wb) && ge0le1(wc)
}

func shift(x, y, sx, sy float64) (float64, float64) {
	return x - sx, y - sy
}

func ge0le1(n float64) bool {
	return n >= 0 && n <= 1
}
