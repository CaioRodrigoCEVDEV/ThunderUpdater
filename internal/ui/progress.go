package ui

import (
	"fmt"
	"math"
	"strings"
)

const barWidth = 40

var fullBlock = "█"
var emptyBlock = "░"

type Progress struct {
	contentLines int
}

func NewProgress(contentLines int) *Progress {
	return &Progress{contentLines: contentLines}
}

func (p *Progress) Start(title string) {
	fmt.Printf("  %s\n", title)
	for i := 0; i < p.contentLines; i++ {
		fmt.Println()
	}
}

func (p *Progress) Update(percent float64, detailLines ...string) {
	moveUp(p.contentLines)

	bar := renderBar(percent)
	clearLine()
	fmt.Printf("  [%s] %3.0f%%\n", bar, percent)

	for _, line := range detailLines {
		clearLine()
		fmt.Printf("  %s\n", line)
	}
}

func (p *Progress) Done(message string) {
	moveUp(p.contentLines)
	fmt.Print("\033[J")
	bar := renderBar(100)
	fmt.Printf("  [%s] 100%%\n", bar)
	fmt.Println()
	fmt.Printf("  %s\n", message)
	fmt.Println()
}

func moveUp(lines int) {
	if lines > 0 {
		fmt.Printf("\033[%dA", lines)
	}
}

func clearLine() {
	fmt.Print("\033[2K\r")
}

func renderBar(percent float64) string {
	filled := int(math.Round(percent / 100 * barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat(fullBlock, filled) + strings.Repeat(emptyBlock, barWidth-filled)
}
