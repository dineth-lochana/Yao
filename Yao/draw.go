package main

import (
	"fmt"
	"strings"

	"github.com/inancgumus/screen"
)

func draw(lines []string) {
	screen.MoveTopLeft()

	title := fmt.Sprintf("Yao - %s%s ", FileName, DirtyIndicator())
	footer := fmt.Sprintf("Ln %d, Col %d ", ActiveCursor.AbsoluteY()+1, ActiveCursor.AbsoluteX()+1)

	headerSuffix := fmt.Sprintf(" %s ─┐", Version)
	fmt.Println(renderRow(-1, "┌ ", title, "─", headerSuffix))

	for y := 0; y < ActiveCursor.BodyHeight; y++ {
		absoluteY := y + ScrollY

		line := ""

		if absoluteY < len(lines) {
			line = lines[absoluteY]
		}

		fmt.Println(renderRow(y, "│ ", line, " ", " │"))
	}

	fmt.Print(renderRow(-1, "└ ", footer, "─", " ^S Save ─ ^R Refresh ─ ^C Save & Exit ┘"))

	screen.MoveTopLeft()
}

func renderRow(y int, prefix, text, padding, suffix string) string {
	var sb strings.Builder

	sb.WriteString(prefix)

	textRunes := []rune(text)
	textLen := len(textRunes)

	prefixLen := len([]rune(prefix)) - 2
	suffixLen := len([]rune(suffix)) - 2
	paddingLen := ActiveCursor.BodyWidth - textLen - prefixLen - suffixLen

	if paddingLen > 0 {
		text += strings.Repeat(padding, paddingLen)
		textRunes = []rune(text)
		textLen = len(textRunes)
	}

	if y != -1 {
		if textLen > 0 {
			from := ScrollX
			to := ScrollX + ActiveCursor.BodyWidth

			if from > textLen {
				from = textLen
			}
			if to > textLen {
				to = textLen
			}

			if from < to {
				text = string(textRunes[from:to])
			} else {
				text = ""
			}
		}

		currentLen := len([]rune(text))
		if currentLen < ActiveCursor.BodyWidth {
			text += strings.Repeat(" ", ActiveCursor.BodyWidth-currentLen)
		}

		if !ActiveCursor.Disabled && ActiveCursor.Y == y {
			char := ActiveCursor.CurrentChar()
			cursorX := ActiveCursor.X

			if cursorX < ActiveCursor.BodyWidth {
				text = ReplaceCharacterAt(text, InvertColors(char), cursorX)
			}
		}
	}

	sb.WriteString(text)
	sb.WriteString(suffix)

	return sb.String()
}
