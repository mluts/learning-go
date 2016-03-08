package main

import (
	"fmt"
	"github.com/mluts/learning-go/alrm/progress_bar"
	"os"
	"os/exec"
	"time"
)

func usage() {
	fmt.Printf(`
  USAGE:
  %s DELAY COMMAND [...ARGS]

  DESCRIPTION:
  Executes COMMAND after DELAY time

  DELAY format:
  1 - one second
  1m - one minute
`, os.Args[0])
	os.Exit(0)
}

const (
	refreshTime = 500 // milliseconds
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	duration, err := time.ParseDuration(os.Args[1])

	if err != nil {
		fmt.Printf("  Can't parse duration! %v\n", err)
		usage()
	}

	cmdArgs := os.Args[2:]

	tickTime := refreshTime * time.Millisecond
	timePassed := 0 * time.Millisecond
	tickTimeFloat := float32(tickTime)

	bar := pbar.New(float32(duration))

	ticker := time.NewTicker(tickTime)

	for range ticker.C {
		timePassed += tickTime
		bar.Advance(tickTimeFloat)
		bar.Print(fmt.Sprintf("%6s / %6s", timePassed, duration))
		if bar.Percent() >= 100 {
			ticker.Stop()
			break
		}
	}

	exec.Command(cmdArgs[0], cmdArgs[1:]...).Run()

	fmt.Print("\n")
}
