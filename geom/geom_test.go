package geom

import (
	"testing"
)

func TestInTriangle(t *testing.T) {
	ds := []struct {
		x, y    float64
		a, b, c Location
		res     bool
	}{
		{
			2.0, 3.0,
			MakeLocation(1.0, 1.0, 0),
			MakeLocation(4.0, 2.0, 0),
			MakeLocation(2.0, 7.0, 0),
			true,
		},
		{
			1.5, 5.0,
			MakeLocation(1.0, 1.0, 0),
			MakeLocation(4.0, 2.0, 0),
			MakeLocation(2.0, 7.0, 0),
			false,
		},
		{
			10.0, 10.0,
			MakeLocation(0.0, 10.0, 0),
			MakeLocation(10.0, 20.0, 0),
			MakeLocation(10.0, 0.0, 0),
			true,
		},
		{
			10.0, 10.0,
			MakeLocation(0.0, 10.0, 0),
			MakeLocation(20.0, 10.0, 0),
			MakeLocation(10.0, 0.0, 0),
			true,
		},
		{
			10.1, 10.1,
			MakeLocation(0.0, 10.0, 0),
			MakeLocation(20.0, 10.0, 0),
			MakeLocation(10.0, 0.0, 0),
			false,
		},
	}

	for _, d := range ds {
		want := d.res
		got := inTriangle(d.x, d.y, d.a, d.b, d.c)
		if got != want {
			t.Errorf("got %v, want %v - %v, %v in %v, %v, %v", got, want, d.x, d.y, d.a, d.b, d.c)
		}
	}
}
