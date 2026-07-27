package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/caio/ThunderUpdaterGO/internal/app"
	"github.com/caio/ThunderUpdaterGO/internal/consoleui"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			ui := consoleui.New()
			msg := fmt.Sprintf("PANIC: %v\n\nStack trace:\n%s", r, string(debug.Stack()))
			ui.PrintFatalError(msg)
			waitForEnter()
		}
	}()

	err := app.Run()
	if err != nil {
		ui := consoleui.New()
		ui.PrintFatalError(err.Error())
		waitForEnter()
	}
}

func waitForEnter() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
