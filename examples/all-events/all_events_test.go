package main

import (
	"os"
	"testing"
)

// Just make sure the example runs
func TestScores(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}

	os.Args = []string{"cmd", "-demo", "../../cs-demos/test/out-2019-08-30T16:23:47+01:00.dem"}

	main()
}
