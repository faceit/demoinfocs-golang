package main

import (
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	ex "github.com/markus-wa/demoinfocs-golang/examples"
	"os"
)

// Run like this: go run print_events.go -demo /path/to/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	checkError(err)
	defer f.Close()

	p := dem.NewParser(f)

	p.RegisterEventHandler(func(e interface{}) {
		round := p.GameState().TotalRoundsPlayed()
		switch e.(type) {
		case events.FrameDone, events.TickDone:
			break
		default:
			fmt.Printf("Round %d %T\n", round, e)
		}

	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		ingameTime := p.CurrentTime()
		progressPercent := p.Progress() * 100
		round := p.GameState().TotalRoundsPlayed()

		fmt.Printf("Round %d started: ingameTime=%s, progress=%f\n",
			round, ingameTime, progressPercent)
	})

	// Parse header
	header, err := p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	checkError(p.ParseToEnd())
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
