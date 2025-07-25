package main

import (
	"fmt"
	"strings"

	"github.com/inancgumus/screen"
)

func draw(lines []string) {
	screen.MoveTopLeft()

	title := fmt.Sprintf("Yao-Write - %s%s ", FileName, DirtyIndicator())
	origY, origX := ActiveCursor.AbsolutePosition()
	footer := fmt.Sprintf("Ln %d, Col %d ", origY+1, origX+1)

	headerSuffix := fmt.Sprintf(" %s ─┐", Version)
	fmt.Println(renderRow(-1, "┌ ", title, "─", headerSuffix))

	visibleLine := 0
	currentWrappedLine := 0

	for i := 0; i < len(WrappedLines) && visibleLine < ActiveCursor.BodyHeight; i++ {
		wl := WrappedLines[i]

		for j := 0; j < wl.LineCount && visibleLine < ActiveCursor.BodyHeight; j++ {
			if currentWrappedLine >= ScrollY {
				line := wl.Wrapped[j]
				fmt.Println(renderRow(visibleLine, "│ ", line, " ", " │"))
				visibleLine++
			}
			currentWrappedLine++
		}
	}

	for visibleLine < ActiveCursor.BodyHeight {
		fmt.Println(renderRow(visibleLine, "│ ", "", " ", " │"))
		visibleLine++
	}

	fmt.Print(renderRow(-1, "└ ", footer, "─", " ^S Save ─ ^R Refresh ─ ^C Save & Exit ┘"))

	screen.MoveTopLeft()
}

func renderRow(y int, prefix, text, padding, suffix string) string {
	var sb strings.Builder

	sb.WriteString(prefix)

	textLen := len([]rune(text))

	prefixLen := len([]rune(prefix)) - 2
	suffixLen := len([]rune(suffix)) - 2

	paddingLen := ActiveCursor.BodyWidth - textLen - prefixLen - suffixLen

	if paddingLen > 0 {
		text += strings.Repeat(padding, paddingLen)
	}

	if y != -1 {
		if textLen > 0 {
			from := ScrollX
			to := ScrollX + ActiveCursor.BodyWidth

			for textLen+1 < to {
				text += " "
				textLen++
			}

			runes := []rune(text)
			if from < 0 {
				from = 0
			}
			if to > len(runes) {
				to = len(runes)
			}
			text = string(runes[from:to])
		}

		if !ActiveCursor.Disabled && ActiveCursor.Y == y {
			char := ActiveCursor.CurrentChar()

			if ActiveCursor.X == ActiveCursor.BodyWidth {
				text = text[1:]
			}

			text = ReplaceCharacterAt(text, InvertColors(char), ActiveCursor.X)
		}
	}

	sb.WriteString(text)

	sb.WriteString(suffix)

	return sb.String()
}
