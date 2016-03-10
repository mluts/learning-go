package main

import (
	"bufio"
	"fmt"
	"github.com/mluts/learning-go/progress_bar"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Printf(`
  USAGE:
  %s DELAY COMMAND [...ARGS]

  DESCRIPTION:
  Executes COMMAND after DELAY time

  DELAY format:
  1s - one second
  1m - one minute
  1.5m - 90 seconds
`, os.Args[0])
	os.Exit(0)
}

const (
	refreshTime = 500 // milliseconds
	exitChar    = 'x'
	pauseChar   = 'c'
)

var (
	pause = false
)

func prepareTerminal() {
	prepareTerminal := exec.Command("/bin/stty", "cbreak", "-echo")
	prepareTerminal.Stdin = os.Stdin
	prepareTerminal.Run()
}

func restoreTerminal() {
	restoreTerminal := exec.Command("/bin/stty", "sane")
	restoreTerminal.Stdin = os.Stdin
	restoreTerminal.Run()
}

func exit() {
	restoreTerminal()
	fmt.Print("\n")
	os.Exit(0)
}

func handleInterrupts() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		exit()
	}
}

func readCommands() {
	prepareTerminal()

	r := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := r.ReadRune()
		if err == nil {
			switch ch {
			case pauseChar:
				pause = !pause
			case exitChar:
				exit()
			}
		}
	}
}

func clearLine() {
	fmt.Print("\r")
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return
	}
	widthStr := strings.Split(string(out), " ")[1]
	widthStr = strings.Trim(widthStr, " \t\n")
	width, err := strconv.ParseInt(widthStr, 10, 16)
	if err != nil {
		return
	}
	for i := 0; i < int(width); i++ {
		fmt.Print(" ")
	}
}

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

	bar := pbar.New(float32(duration))

	ticker := time.NewTicker(tickTime)

	go handleInterrupts()
	go readCommands()

	for range ticker.C {
		str := []string{fmt.Sprintf("Duration: %s", duration)}
		str = append(str, fmt.Sprintf("(%c - PAUSE, %c - EXIT)", pauseChar, exitChar))
		if !pause {
			timePassed += tickTime
			bar.Advance(float32(tickTime))
		}
		if pause {
			str = append(str, "[paused]")
		}
		if bar.Percent() >= 100 {
			ticker.Stop()
			break
		}
		clearLine()
		bar.Print(strings.Join(str, " "))
	}

	exec.Command("/bin/sh", "-c", strings.Join(cmdArgs, " ")).Run()

	exit()
}
