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

	p := dem.NewCaptureParser(f)
	startRound := 3
	endRound := 5

	p.RegisterEventHandler(func(e events.MatchStart) {
		p.EndCapture()
	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		ingameTime := p.CurrentTime()
		progressPercent := p.Progress() * 100
		round := p.GameState().TotalRoundsPlayed()

		fmt.Printf("Round %d started: ingameTime=%s, progress=%f\n",
			round, ingameTime, progressPercent)

		if round == startRound {
			p.BeginCapture()
		}
	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		ingameTime := p.CurrentTime()
		progressPercent := p.Progress() * 100
		round := p.GameState().TotalRoundsPlayed()

		fmt.Printf("Round %d finished: ingameTime=%s, progress=%f\n",
			round, ingameTime, progressPercent)

		if round == endRound {
			p.EndCapture()
		}
	})

	// Parse header
	header, err := p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	checkError(p.ParseToEnd())

	fmt.Printf("got %d bytes", p.Out.Len())
	checkError(p.WriteOut("../../cs-demos/round.dem"))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
