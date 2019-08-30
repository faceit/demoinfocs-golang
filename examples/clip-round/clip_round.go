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

	cp.RegisterEventHandler(func(e events.RoundEndOfficial) {
		ingameTime := cp.CurrentTime()
		round := cp.GameState().TotalRoundsPlayed()
		fmt.Printf("Round %d ingameTime=%s\n", round, ingameTime)

		if round == 1 {
			cp.EndCapture()
		}

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

	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		ingameTime := p.CurrentTime()
		round := p.GameState().TotalRoundsPlayed()
		fmt.Printf("Round %d ingameTime=%s\n", round, ingameTime)

	})

	// Parse header
	header, err = p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	checkError(p.ParseToEnd())
}

func clipRound(cp *dem.CaptureParser, round int) {
	startRound := 5
	endRound := 5

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
