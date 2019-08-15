package main

import (
	"bytes"
	"fmt"
	dem "github.com/faceit/demoinfocs-golang"
	"github.com/faceit/demoinfocs-golang/common"
	ex "github.com/faceit/demoinfocs-golang/examples"
	"io/ioutil"
	"os"
)

// Run like this: go run print_events.go -demo /path/to/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	defer f.Close()
	checkError(err)

	var buf bytes.Buffer
	mr := common.NewStoppableReader(f, &buf)
	mr.Begin()

	p := dem.NewParser(mr)

	// Parse header
	header, err := p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Parse to end
	err = p.ParseToEnd()
	checkError(err)

	fmt.Printf("got %d bytes", buf.Len())
	err = ioutil.WriteFile("../../cs-demos/round.dem", buf.Bytes(), 777)
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
