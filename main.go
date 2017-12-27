package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codemodus/sigmon"
	"github.com/daved/halitego/bot/hyena"
	"github.com/daved/halitego/ops"
)

func main() {
	sm := sigmon.New(func(*sigmon.SignalMonitor) {
		panic("startup interrupted")
	})
	sm.Run()

	l := log.New(ioutil.Discard, "", 0)
	o := ops.New("Hyena")
	c := hyena.New(l, o.InitialBoard())

	if false {
		fn := fmt.Sprintf("%d_%s", o.ID(), "game.log")
		defer setLoggerOutput(l, fn)()
	}

	sm.Set(func(*sigmon.SignalMonitor) {
		o.Stop()
	})

	o.Run(l, c)
	sm.Stop()
}

func setLoggerOutput(l *log.Logger, filename string) func() {
	var (
		lFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
		lPerms = os.FileMode(0664)
	)

	f, err := os.OpenFile(filename, lFlags, lPerms)
	if err != nil {
		panic(err)
	}

	l.SetOutput(f)

	return func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
