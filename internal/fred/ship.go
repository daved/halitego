package fred

import (
	"math"

	"github.com/daved/halitego/ops"
)

type fredShip struct {
	ops.Ship
}

// Navigate demonstrates how the player might negotiate obsticles between
// a ship and its target
func (s fredShip) Navigate(target ops.Entity, gameMap ops.Board) string {
	ob := gameMap.ObstaclesBetween(s.Entity, target)

	if !ob {
		return s.NavigateTo(target, gameMap)
	}

	x0 := math.Min(s.X, target.X)
	x2 := math.Max(s.X, target.X)
	y0 := math.Min(s.Y, target.Y)
	y2 := math.Max(s.Y, target.Y)

	dx := (x2 - x0) / 5
	dy := (y2 - y0) / 5
	bestdist := 1000.0
	bestTarget := target

	for x1 := x0; x1 <= x2; x1 += dx {
		for y1 := y0; y1 <= y2; y1 += dy {
			intermediateTarget := ops.Entity{
				X:      x1,
				Y:      y1,
				Radius: 0,
				Health: 0,
				Owner:  0,
				ID:     -1,
			}
			ob1 := gameMap.ObstaclesBetween(s.Entity, intermediateTarget)
			if !ob1 {
				ob2 := gameMap.ObstaclesBetween(intermediateTarget, target)
				if !ob2 {
					totdist := math.Sqrt(math.Pow(x1-x0, 2)+math.Pow(y1-y0, 2)) + math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2))
					if totdist < bestdist {
						bestdist = totdist
						bestTarget = intermediateTarget

					}
				}
			}
		}
	}

	return s.NavigateTo(bestTarget, gameMap)
}
