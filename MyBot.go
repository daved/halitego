package main

import (
	"log"
	"os"
	"strconv"

	"github.com/daved/halitefred/internal/hlt"
)

func main() {
	botName := "Fred The SpaceGopher"
	logging := true
	logFileSuffix := "_gamelog.log"

	conn := hlt.NewConnection(botName)

	if logging {
		fname := strconv.Itoa(conn.PlayerTag) + logFileSuffix

		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		log.SetOutput(f)
	}

	for i := 1; ; i++ {
		gmap := conn.UpdateMap()
		env := gmap.Players[gmap.MyID]

		cmds := []string{}
		for k := range env.Ships {
			s := env.Ships[k]

			if s.DockingStatus == hlt.Undocked {
				cmds = append(cmds, hlt.StrategyBasicBot(s, gmap))
			}
		}

		log.Printf("Turn %v\n", i)

		conn.SubmitCommands(cmds)
	}
}
