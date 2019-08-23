package main

import (
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	ex "github.com/markus-wa/demoinfocs-golang/examples"
	"os"
	"time"
)

// Run like this: go run print_events.go -demo /path/to/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	checkError(err)
	defer f.Close()

	cp := dem.NewCaptureParser(f)

	cp.RegisterEventHandler(func(e events.DataTablesParsed) {
		cp.EndCapture()
	})

	cp.RegisterEventHandler(func(e events.RoundEndOfficial) {
		ingameTime := cp.CurrentTime()
		round := cp.GameState().TotalRoundsPlayed()
		fmt.Printf("Round %d ingameTime=%s\n", round, ingameTime)

		clipRound(cp, round)
	})

	// Parse header
	header, err := cp.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	checkError(cp.ParseToEnd())

	path := fmt.Sprintf("../../cs-demos/test/out-%s.dem", time.Now().Format(time.RFC3339))
	checkError(cp.WriteOut(path))

	f, err = os.Open(path)
	checkError(err)
	defer f.Close()

	p := dem.NewParser(f)

	p.RegisterEventHandler(func(e interface{}) {
		round := p.GameState().TotalRoundsPlayed()
		switch event := e.(type) {
		case events.FrameDone, events.TickDone:
			break
		case events.ParserWarn:
			fmt.Printf("Round %d Parser Warn: %s\n", round, event.Message)
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
	header, err = p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	checkError(p.ParseToEnd())
}

func clipRound(cp *dem.CaptureParser, round int) {
	startRound := 2
	endRound := 2

	if round == startRound-1 {
		cp.BeginCapture()
	}

	if round == endRound {
		cp.EndCapture()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
