package ops

import (
	"strconv"
)

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
