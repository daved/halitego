package ops

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// Board describes the current state of the game
type Board struct {
	xLen int
	yLen int
	pCt  int
	ps   []Planet
	ss   [][]Ship
}

// MakeBoard from a slice of game state tokens
func MakeBoard(xLen, yLen int, gameString string) Board {
	tokens := strings.Split(gameString, " ")
	pCt, err := strconv.Atoi(tokens[0])
	if err != nil {
		panic(err)
	}
	tokens = tokens[1:]

	b := Board{
		xLen: xLen,
		yLen: yLen,
		ss:   make([][]Ship, pCt),
	}

	for k := range b.ss {
		curID, err := strconv.Atoi(tokens[0])
		if err != nil {
			panic(err)
		}
		if curID != k {
			panic("curID is not iteration when making board")
		}

		shipCt, err := strconv.Atoi(tokens[1])
		if err != nil {
			panic(err)
		}
		tokens = tokens[2:]

		for i := 0; i < shipCt; i++ {
			s, trimmedTokens := MakeShip(k, tokens)
			tokens = trimmedTokens

			b.ss[k] = append(b.ss[k], s)
		}
	}

	plntCt, err := strconv.Atoi(tokens[0])
	if err != nil {
		panic(err)
	}
	tokens = tokens[1:]

	for i := 0; i < plntCt; i++ {
		p, trimmedTokens := MakePlanet(tokens)
		tokens = trimmedTokens

		b.ps = append(b.ps, p)
	}

	return b
}

// Dimensions ...
func (b *Board) Dimensions() (int, int) {
	return b.xLen, b.yLen
}

// PlayerCt ...
func (b *Board) PlayerCt() int {
	return b.pCt
}

// Planets ...
func (b *Board) Planets() []Planet {
	return b.ps
}

// Ships ...
func (b *Board) Ships() [][]Ship {
	return b.ss
}

// ObstaclesBetween demonstrates how the player might determine if the path
// between two enitities is clear
func (b *Board) ObstaclesBetween(start Entity, end Entity) bool {
	x1 := start.X
	y1 := start.Y
	x2 := end.X
	y2 := end.Y
	dx := x2 - x1
	dy := y2 - y1
	a := dx*dx + dy*dy + 1e-8
	crossterms := x1*x1 - x1*x2 + y1*y1 - y1*y2

	var entities []Entity
	for _, v := range b.ps {
		entities = append(entities, v.Entity)
	}
	for _, v := range b.ss {
		for _, y := range v {
			entities = append(entities, y.Entity)
		}
	}

	for i := 0; i < len(entities); i++ {
		entity := entities[i]
		if entity.X == start.X || entity.X == end.X {
			continue
		}

		x0 := entity.X
		y0 := entity.Y

		closestDistance := Distance(end, entity)
		if closestDistance < entity.Radius+1 {
			return true
		}

		bz := -2 * (crossterms + x0*dx + y0*dy)
		t := -bz / (2 * a)

		if t <= 0 || t >= 1 {
			continue
		}

		closestX := start.X + dx*t
		closestY := start.Y + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-x0, 2) * +math.Pow(closestY-y0, 2))

		if closestDistance <= entity.Radius+start.Radius+1 {
			return true
		}
	}
	return false
}

// NearestPlanetsByDistance orders all planets based on their proximity
// to a given ship from nearest for farthest
func (b *Board) NearestPlanetsByDistance(ship Ship) []Planet {
	planets := b.ps

	for i := 0; i < len(planets); i++ {
		planets[i].Distance = Distance(ship, planets[i])
	}

	sort.Sort(byDist(planets))

	return planets
}

type byDist []Planet

func (a byDist) Len() int           { return len(a) }
func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
