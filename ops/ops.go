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
	Command(Board) []Messenger
}

// Operations ...
type Operations struct {
	id   int
	xLen int
	yLen int
	r    *bufio.Reader
	w    io.Writer
	done chan struct{}
}

// New ...
func New(botName string) *Operations {
	o := &Operations{
		r:    bufio.NewReader(os.Stdin),
		w:    os.Stdout,
		done: make(chan struct{}),
	}

	o.id = o.readLineInt()
	o.xLen, o.yLen = o.readLineInts()

	o.send(botName)

	return o
}

// ID ...
func (o *Operations) ID() int {
	return o.id
}

// Stop ...
func (o *Operations) Stop() {
	select {
	case <-o.done:
	default:
		close(o.done)
	}
}

// Wait ...
func (o *Operations) Wait() {
	<-o.done
}

// Run gathers and submits game commands to the GameCommunicator.
func (o *Operations) Run(l Logger, c Commander) {
	for i := 1; ; i++ {
		select {
		case <-o.done:
			return
		default:
			o.runIteration(l, i, c)
		}
	}
}

func (o *Operations) runIteration(l Logger, iter int, c Commander) {
	l.Printf("--- Turn %v\n", iter)

	b := MakeBoard(o.xLen, o.yLen, o.readLineString())
	l.Printf("   Parsed Board")

	ms := c.Command(b)

	m := messengersToSysMsg(ms)
	l.Printf("   System Message: %s\n", m)

	o.send(m)
}

func (o *Operations) send(msg string) {
	fmt.Fprintf(o.w, "%s\n", msg)
}

func (o *Operations) readLine() []byte {
	bs, err := o.r.ReadBytes('\n')
	if err != nil {
		panic(err)
	}

	return bytes.TrimSpace(bs)
}

func (o *Operations) readLineString() string {
	return string(o.readLine())
}

func (o *Operations) readLineInt() int {
	i, err := strconv.Atoi(o.readLineString())
	if err != nil {
		panic(err)
	}

	return i
}

func (o *Operations) readLineInts() (int, int) {
	xy := strings.Split(o.readLineString(), " ")

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
