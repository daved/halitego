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
	alpha = uint8(0xFF)
	red   = icolor.RGBA{0xFF, 0x00, 0x00, alpha}
	pur   = icolor.RGBA{0xFF, 0x00, 0x99, alpha}
	blu   = icolor.RGBA{0x00, 0x00, 0xFF, alpha}
	ora   = icolor.RGBA{0xFF, 0xA5, 0x00, alpha}
	yel   = icolor.RGBA{0xFF, 0xFF, 0x00, alpha}

	pSm, pMd, pLg = 15.0, 20.0, 25.0
	sMd           = 2.0
	tMd           = 2.0

	rng = rand.New(rand.NewSource(time.Now().Unix()))

	xaxis = 0
	yaxis = 1
)

func init() {
	_ = xaxis
	_ = yaxis
}

func main() {
	c := newGraphicContext(400, 300)

	destPlanet := makeEntity(100, 250, pSm, pur)
	planets := []entity{
		makeEntity(100, 40, pSm, red),
		makeEntity(200, 80, pMd, red),
		makeEntity(150, 120, pLg, red),
		makeEntity(160, 160, pLg, red),
		makeEntity(175, 200, pLg, red),
		makeEntity(185, 210, pLg, red),
		makeEntity(165, 190, pLg, red),
		makeEntity(340, 240, pSm, red),
		makeEntity(240, 240, pSm, red),
		makeEntity(40, 220, pMd, red),
		makeEntity(40, 260, pMd, red),
		makeEntity(40, 300, pMd, red),
		destPlanet,
	}
	c.addEntities(planets...)

	ship := makeEntity(160, 130, sMd, blu)
	c.addEntities(ship)

	dest := geom.BufferedLocation(10, destPlanet, ship)
	destMarker := makeEntityFromMarker(dest, yel)
	c.addEntities(destMarker)

	var adj func(int, []entity, geom.Marker, geom.Marker)
	adj = func(count int, ents []entity, t, s geom.Marker) {
		count++
		if count > 256 {
			return
		}

		ob := geom.Obstacles(entsToMarkers(ents), t, s)
		if !ob {
			return
		}

		buf := float64(rng.Intn(23) + 17)
		dir := geom.Left
		if time.Now().Nanosecond()%2 == 0 {
			dir = geom.Right
		}

		t = geom.PerpindicularLocation(buf, dir, t, ship)
		destMarker = makeEntityFromMarker(t, ora)
		c.addEntities(destMarker)

		adj(count, ents, t, s)
	}

	adj(0, planets, dest, ship)

	_ = c.save("out.png")
}

type graphicContext struct {
	dest *image.RGBA
	*draw2dimg.GraphicContext
	ents  []entity
	lines []line
}

func newGraphicContext(x, y int) *graphicContext {
	dest := image.NewRGBA(image.Rect(0, 0, x, y))

	gctx := &graphicContext{
		dest:           dest,
		GraphicContext: draw2dimg.NewGraphicContext(dest),
	}

	gctx.SetLineWidth(1)

	return gctx
}

/*func (ctx *graphicContext) obstacles(ms []geom.Marker, b, a geom.Marker) bool {
	bx, by := b.Coords()
	ax, ay := a.Coords()

	pfl := PerpindicularLocation(-b.Radius()+a.Radius(), Left, b, a)
	pfr := PerpindicularLocation(-b.Radius()+a.Radius(), Right, b, a)
	pbl := PerpindicularLocation(0, Right, a, b)
	pbr := PerpindicularLocation(0, Left, a, b)

	for _, po := range ms {
		pox, poy := po.Coords()
		if (pox == bx && poy == by) || (pox == ax && poy == ay) {
			continue
		}

		poxa, poxb := pox-po.Radius(), pox+po.Radius()
		poya, poyb := poy+po.Radius(), poy-po.Radius()

	}

	return false
}*/

func (ctx *graphicContext) save(filename string) error {
	for _, v := range ctx.ents {
		v.draw(ctx.GraphicContext)
	}

	ctx.Close()

	return draw2dimg.SaveToPngFile(filename, ctx.dest)
}

func (ctx *graphicContext) addEntities(ents ...entity) {
	ctx.ents = append(ctx.ents, ents...)
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
		r = tMd
	}

	return makeEntity(x, y, r, color)
}

func (e *entity) draw(gctx *draw2dimg.GraphicContext) {
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

func (l *line) draw(gctx *draw2dimg.GraphicContext) {
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
