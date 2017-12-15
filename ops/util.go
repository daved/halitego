package ops

import (
	"math"
	"strconv"
)

// DegToRad converts degrees to radians.
func DegToRad(d float64) float64 {
	return d / 180 * math.Pi
}

// RadToDeg converts radians to degrees.
func RadToDeg(r float64) float64 {
	return r / math.Pi * 180
}

func readTokenString(tokens []string, k int) string {
	if k >= len(tokens) {
		panic("index out of token range")
	}

	return tokens[k]
}

func readTokenInt(tokens []string, k int) int {
	n, err := strconv.Atoi(readTokenString(tokens, k))
	if err != nil {
		panic(err)
	}

	return n
}

func readTokenFloat(tokens []string, k int) float64 {
	n, err := strconv.ParseFloat(readTokenString(tokens, k), 64)
	if err != nil {
		panic(err)
	}

	return n
}
