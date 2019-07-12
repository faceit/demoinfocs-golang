package main

import (
	"bytes"
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/common"
	"github.com/markus-wa/demoinfocs-golang/events"
	ex "github.com/markus-wa/demoinfocs-golang/examples"
	"io/ioutil"
	"os"
)

// Run like this: go run print_events.go -demo /path/to/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	defer f.Close()
	checkError(err)

	var buf bytes.Buffer
	mr := &teeReader{r: f, w: &buf}
	mr.Begin()

	p := dem.NewParser(mr)

	// Parse header
	header, err := p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	first := true

	// Register handler on round end to figure out who won
	p.RegisterEventHandler(func(event events.RoundEnd) {
		if first {
			gs := p.GameState()
			fmt.Printf("got event roundendofficial at tick %d", gs.IngameTick())
			fmt.Println()
			mr.End()
			first = false
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	checkError(err)

	fmt.Printf("got %d bytes", buf.Len())
	err = ioutil.WriteFile("singleround.dem", buf.Bytes(), 777)
	checkError(err)
}

func formatPlayer(p *common.Player) string {
	if p == nil {
		return "?"
	}

	switch p.Team {
	case common.TeamTerrorists:
		return "[T]" + p.Name
	case common.TeamCounterTerrorists:
		return "[CT]" + p.Name
	}
	return p.Name
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
