package ops

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Logger describes the halitego logging behavior.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Commander ...
type Commander interface {
	Cmds(Board) []fmt.Stringer
}

// Operations ...
type Operations struct {
	id   int
	xLen int
	yLen int
	r    *bufio.Reader
	w    io.Writer
}

// New ...
func New(botName string) *Operations {
	c := &Operations{
		r: bufio.NewReader(os.Stdin),
		w: os.Stdout,
	}

	c.id = c.readLineInt()
	c.xLen, c.yLen = c.readLineInts()

	c.send(botName)

	return c
}

// ID ...
func (ctrl *Operations) ID() int {
	return ctrl.id
}

// Run gathers and submits game commands to the GameCommunicator.
func (ctrl *Operations) Run(done <-chan struct{}, l Logger, c Commander) {
	for i := 1; ; i++ {
		select {
		case <-done:
			return
		default:
			ctrl.runIteration(l, c, i)
		}
	}
}

func (ctrl *Operations) runIteration(l Logger, c Commander, iter int) {
	l.Printf("--- Turn %v\n", iter)

	b := ctrl.board()
	l.Printf("   Parsed Board")

	cmds := c.Cmds(b)

	m := cmdsToMsg(cmds)
	l.Printf("   Sending Message: %s\n", m)
	ctrl.send(m)
}

func (ctrl *Operations) send(msg string) {
	fmt.Fprintf(ctrl.w, "%s\n", msg)
}

func (ctrl *Operations) readLine() []byte {
	bs, err := ctrl.r.ReadBytes('\n')
	if err != nil {
		panic(err)
	}

	return bytes.TrimSpace(bs)
}

func (ctrl *Operations) readLineString() string {
	return string(ctrl.readLine())
}

func (ctrl *Operations) readLineInt() int {
	i, err := strconv.Atoi(ctrl.readLineString())
	if err != nil {
		panic(err)
	}

	return i
}

func (ctrl *Operations) readLineInts() (int, int) {
	xy := strings.Split(ctrl.readLineString(), " ")

	x, err := strconv.Atoi(xy[0])
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		panic(err)
	}

	return x, y
}

func (ctrl *Operations) board() Board {
	s := ctrl.readLineString()

	return MakeBoard(ctrl.xLen, ctrl.yLen, s)
}

func cmdsToMsg(cmds []fmt.Stringer) string {
	s := ""
	for _, v := range cmds {
		s += v.String() + " "
	}

	if len(s) < 1 {
		return ""
	}

	return s[:len(s)-1]
}
