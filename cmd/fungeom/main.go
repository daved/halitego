package main

import (
	"image"
	icolor "image/color"
	"math/rand"
	"time"

	"github.com/daved/halitego/geom"
	"github.com/daved/halitego/ops"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	opaq = uint8(0xFF)
	tran = uint8(0x80)
	redo = icolor.RGBA{0xFF, 0x00, 0x00, opaq}
	redt = altAlpha(redo, tran)
	puro = icolor.RGBA{0xFF, 0x00, 0x99, opaq}
	purt = altAlpha(puro, tran)
	bluo = icolor.RGBA{0x00, 0x00, 0xFF, opaq}
	blut = altAlpha(bluo, tran)
	orao = icolor.RGBA{0xFF, 0xA5, 0x00, opaq}
	orat = altAlpha(orao, tran)
	yelo = icolor.RGBA{0xFF, 0xFF, 0x00, opaq}
	yelt = altAlpha(yelo, tran)

	pSm, pMd, pLg = 15.0, 20.0, 25.0
	sMd           = 2.0
	lMd           = 2.0

	rng = rand.New(rand.NewSource(time.Now().Unix()))

	xaxis, yaxis = 0, 1
)

func altAlpha(rgba icolor.RGBA, alpha uint8) icolor.RGBA {
	rgba.A = alpha
	return rgba
}

func main() {
	c := newGraphicContext(400, 300)

	destPlanet := makeEntity(100, 250, pSm, puro)
	planets := []entity{
		makeEntity(100, 40, pSm, redo),
		makeEntity(200, 80, pMd, redo),
		makeEntity(150, 120, pLg, redo),
		makeEntity(175, 200, pLg, redo),
		makeEntity(340, 240, pSm, redo),
		makeEntity(240, 240, pSm, redo),
		makeEntity(40, 260, pMd, redo),
		destPlanet,
	}
	c.addDrawers(entsToDrawers(planets)...)

	ship := makeEntity(260, 140, sMd, bluo)
	c.addDrawers(ship)

	dest := geom.BufferedLocation(10, destPlanet, ship)
	destMarker := makeEntityFromMarker(dest, yelo)
	c.addDrawers(destMarker)

	var adj func(int, []entity, geom.Marker, geom.Marker)
	adj = func(count int, ents []entity, t, s geom.Marker) {
		count++
		if count > 256 {
			return
		}

		ob := c.obstacles(entsToMarkers(ents), t, s)
		if !ob {
			return
		}

		buf := float64(rng.Intn(23) + 17)
		dir := geom.Left
		if time.Now().Nanosecond()%2 == 0 {
			dir = geom.Right
		}

		t = geom.PerpindicularLocation(buf, dir, t, ship)
		destMarker = makeEntityFromMarker(t, orao)
		c.addDrawers(destMarker)

		adj(count, ents, t, s)
	}

	adj(0, planets, dest, ship)

	_ = c.save("out.png")
}

type drawer interface {
	draw(*draw2dimg.GraphicContext)
}

type graphicContext struct {
	dest *image.RGBA
	*draw2dimg.GraphicContext
	drawers []drawer
}

func newGraphicContext(x, y int) *graphicContext {
	dest := image.NewRGBA(image.Rect(0, 0, x, y))

	c := &graphicContext{
		dest:           dest,
		GraphicContext: draw2dimg.NewGraphicContext(dest),
	}

	c.SetLineWidth(1)

	return c
}

func (c *graphicContext) obstacles(ms []geom.Marker, b, a geom.Marker) bool {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	pfl := geom.PerpindicularLocation(-b.Radius()+a.Radius(), geom.Left, b, a)
	pfr := geom.PerpindicularLocation(-b.Radius()+a.Radius(), geom.Right, b, a)
	pbr := geom.PerpindicularLocation(0, geom.Left, a, b)
	pbl := geom.PerpindicularLocation(0, geom.Right, a, b)

	cs := makeCoordsFromGeomLocators(pfl, pfr, pbr, pbl)
	c.addDrawers(makePoly(cs, purt))

	for _, po := range ms {
		pox, poy := po.Coords()
		if (pox == bx && poy == by) || (pox == ax && poy == ay) {
			continue
		}

		poxa, poxb := pox-po.Radius(), pox+po.Radius()
		poya, poyb := poy+po.Radius(), poy-po.Radius()
		_, _, _, _ = poxa, poxb, poya, poyb
	}

	return false
}

func (c *graphicContext) save(filename string) error {
	for _, v := range c.drawers {
		v.draw(c.GraphicContext)
	}

	return draw2dimg.SaveToPngFile(filename, c.dest)
}

func (c *graphicContext) addDrawers(ents ...drawer) {
	c.drawers = append(c.drawers, ents...)
}

type entity struct {
	ops.Entity
	icolor.Color
}

func makeEntity(x, y, radius float64, color icolor.Color) entity {
	return entity{
		Entity: ops.MakeEntity(x, y, radius, 0, 0, 0),
		Color:  color,
	}
}

func makeEntityFromMarker(m geom.Marker, color icolor.Color) entity {
	x, y := m.Coords()
	r := m.Radius()
	if r == 0 {
		r = lMd
	}

	return makeEntity(x, y, r, color)
}

func (e entity) draw(gctx *draw2dimg.GraphicContext) {
	gctx.SetFillColor(e.Color)

	cx, cy := e.Coords()
	draw2dkit.Circle(gctx, cx, cy, e.Radius())
	gctx.FillStroke()
}

func entsToMarkers(ents []entity) []geom.Marker {
	var ms []geom.Marker
	for _, v := range ents {
		ms = append(ms, v)
	}

	return ms
}

func entsToDrawers(ents []entity) []drawer {
	var ds []drawer
	for _, v := range ents {
		ds = append(ds, v)
	}

	return ds
}

type line struct {
	a, b float64
	axis int
	icolor.Color
}

func makeLine(a, b float64, axis int, color icolor.Color) line {
	return line{
		a:     a,
		b:     b,
		axis:  axis,
		Color: color,
	}
}

func (l line) draw(gctx *draw2dimg.GraphicContext) {
	gctx.SetFillColor(l.Color)

	x1, y1 := l.a, 0.0
	x2, y2 := l.b, 4.0

	if l.axis == yaxis {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}

	draw2dkit.Rectangle(gctx, x1, y1, x2, y2)
	gctx.FillStroke()
}

type coord struct {
	x, y float64
}

func makeCoordsFromGeomLocators(ls ...geom.Locator) []coord {
	var cs []coord
	for _, v := range ls {
		x, y := v.Coords()
		cs = append(cs, coord{x, y})
	}

	return cs
}

type poly struct {
	cs []coord
	icolor.Color
}

func makePoly(cs []coord, color icolor.Color) poly {
	return poly{
		cs:    cs,
		Color: color,
	}
}

func (p poly) draw(gctx *draw2dimg.GraphicContext) {
	gctx.SetFillColor(p.Color)

	if len(p.cs) < 3 {
		return
	}

	gctx.MoveTo(p.cs[0].x, p.cs[0].y)
	for _, v := range p.cs[1:] {
		gctx.LineTo(v.x, v.y)
	}
	gctx.Close()

	gctx.FillStroke()
}
