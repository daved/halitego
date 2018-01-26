package geom

import (
	"math"
)

// Direction ...
type Direction int

// Direction constants.
const (
	Left Direction = iota
	Right
)

// Locator describes types which are able to show their coordinates.
type Locator interface {
	Coords() (float64, float64)
}

// Sizer describes types which are able to show their size.
type Sizer interface {
	Radius() float64
}

// Marker describes types which are able to be located and sized.
type Marker interface {
	Locator
	Sizer
}

// CenterDistance returns the distance between two instances of Locator types.
func CenterDistance(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	d := distanceBetween(bx, by, ax, ay)

	return d
}

// EdgeDistance returns the distance between two instances of Locator types.
func EdgeDistance(b, a Marker) float64 {
	return CenterDistance(b, a) - b.Radius() - a.Radius()
}

// Degrees returns the angle in degrees between two instances of Locator types.
func Degrees(b, a Locator) float64 {
	return radiansToDegrees(Radians(b, a))
}

// BoundDegrees ...
func BoundDegrees(b, a Locator) int {
	d := Degrees(b, a)

	bd := int(math.Ceil(d - .5))
	if d > 0.0 {
		bd = int(math.Floor(d + .5))
	}

	return ((bd % 360) + 360) % 360
}

// Radians returns the angle in radians between two instances of Locator types.
func Radians(b, a Locator) float64 {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	return radiansBetween(bx, by, ax, ay)
}

// BufferedLocation returns the closest point from Marker "a" to Marker "b" that is at
// least a distance of "min" from Marker "b".
func BufferedLocation(buffer float64, b, a Marker) Location {
	d := EdgeDistance(b, a) - buffer
	r := Radians(b, a)

	ax, ay := a.Coords()
	x := ax + (d * math.Cos(r))
	y := ay + (d * math.Sin(r))

	return MakeLocation(x, y, 0)
}

// PerpindicularLocation ...
func PerpindicularLocation(buffer float64, dir Direction, b, a Marker) Location {
	d := b.Radius() + buffer
	r := Radians(b, a)
	bx, by := b.Coords()

	dMult := 1.0
	if dir == Left {
		dMult = -1.0
	}

	x := bx + ((d * math.Cos((math.Pi/2.0)+r)) * dMult)
	y := by + ((d * math.Sin((math.Pi/2.0)-r)) * dMult)

	return MakeLocation(x, y, 0)
}

// Obstacles demonstrates how the player might determine if the path
// between two enitities is clear
func Obstacles(ms []Marker, b, a Marker) bool {
	pfl, pfr, pbr, pbl := pathPolygon(b, a)
	pmaxX, pminX, pmaxY, pminY := minmax(pfl, pfr, pbr, pbl)

	for _, po := range ms {
		if !potentialObstacle(po, pmaxX, pminX, pmaxY, pminY) {
			continue
		}

		if obstacle(po, pfl, pfr, pbr, pbl) {
			return true
		}
	}

	return false
}

func pathPolygon(b, a Marker) (h, i, j, k Locator) {
	pfl := PerpindicularLocation(-b.Radius()+a.Radius(), Left, b, a)
	pfr := PerpindicularLocation(-b.Radius()+a.Radius(), Right, b, a)
	pbr := PerpindicularLocation(0, Left, a, b)
	pbl := PerpindicularLocation(0, Right, a, b)

	return pfl, pfr, pbr, pbl
}

func minmax(coords ...Locator) (maxX, minX, maxY, minY float64) {
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

func potentialObstacle(ob Marker, maxX, minX, maxY, minY float64) bool {
	x, y := ob.Coords()
	r := ob.Radius()
	maxX, minX, maxY, minY = maxX+r, minX-r, maxY+r, minY+r

	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func obstacle(ob Marker, coords ...Locator) bool {
	if inPolygon(ob, coords...) {
		return true
	}

	for _, pi := range potentialIntrusions(ob, coords...) {
		if inPolygon(pi, coords...) {
			return true
		}
	}

	return false
}

func inPolygon(m Marker, coords ...Locator) bool {
	if len(coords) < 3 {
		return false
	}

	ct := 0
	x, y := m.Coords()
	a := coords[0]
	b := coords[1]

	for k, c := range coords {
		if m.Radius() > 0 && CenterDistance(m, c) <= m.Radius() {
			return true
		}

		if k < 2 {
			continue
		}

		if inTriangle(x, y, a, b, c) {
			ct++
		}

		b = c
	}

	return ct%2 != 0
}

func potentialIntrusions(gm Marker, coords ...Locator) []Marker {
	var ms []Marker
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
		pi := MakeLocation(ix, iy, 0)

		p1x, p1y = p2x, p2y

		if CenterDistance(gm, pi) > gm.Radius() {
			continue
		}

		ms = append(ms, pi)
		ms = append(ms, MakeLocation(ix+.0001, iy+.0001, 0)) // rough fix for overlap errors
		ms = append(ms, MakeLocation(ix-.0001, iy-.0001, 0)) // rough fix for overlap errors
	}

	return ms
}

// https://math.stackexchange.com/questions/51326
func inTriangle(x, y float64, a, b, c Locator) bool {
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

func distanceBetween(bx, by, ax, ay float64) float64 {
	dx := ax - bx
	dy := ay - by

	return math.Sqrt(dx*dx + dy*dy)
}

func radiansBetween(bx, by, ax, ay float64) float64 {
	dx := bx - ax
	dy := by - ay

	return math.Atan2(dy, dx)
}

func radiansToDegrees(r float64) float64 {
	return r / math.Pi * 180
}
