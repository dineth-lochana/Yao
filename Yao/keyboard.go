package main

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

func readSingleKey() (rune, keyboard.Key) {
	char, key, err := keyboard.GetKey()
	if err != nil {
		fmt.Println("Failed to read key:", err)

		os.Exit(1)
	}

	return char, key
}

func readAndHandleKey() {
	char, key := readSingleKey()

	switch key {
	case keyboard.KeyCtrlS:
		save()
	case keyboard.KeyCtrlC:
		save()
		exit()
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		deleteChar()
	case keyboard.KeyDelete:
		ActiveCursor.Right()
		deleteChar()

	case keyboard.KeyArrowUp:
		ActiveCursor.Up()
	case keyboard.KeyArrowDown:
		ActiveCursor.Down()
	case keyboard.KeyArrowLeft:
		ActiveCursor.Left()
	case keyboard.KeyArrowRight:
		ActiveCursor.Right()
	case keyboard.KeyHome:
		ActiveCursor.MoveStartOfLine()
	case keyboard.KeyEnd:
		ActiveCursor.MoveEndOfLine()

	case keyboard.KeyCtrlR:
		w, h := screen.Size()
		ActiveCursor.SetDimensions(w-4, h-2)

	case keyboard.KeyEnter:
		insertNewLine()
	case keyboard.KeySpace:
		insertChars(' ')
	case keyboard.KeyTab:
		insertChars(' ', ' ', ' ', ' ')
	default:
		if char != 0 {
			insertChars(char)
		}
	}
}

func insertNewLine() {
	x := ActiveCursor.AbsoluteX()
	y := ActiveCursor.AbsoluteY()

	currentLine := Lines[y]

	linesBefore := Lines[:y]
	linesAfter := Lines[y+1:]

	newLines := append([]string{}, linesBefore...)

	if x == 0 {
		newLines = append(newLines, "", currentLine)
	} else if x == len(currentLine) {
		newLines = append(newLines, currentLine, "")
	} else {
		before := currentLine[:x]
		after := currentLine[x:]

		newLines = append(newLines, before, after)
	}

	Lines = append(newLines, linesAfter...)

	ActiveCursor.MoveStartOfLine()
	ActiveCursor.Down()

	IsDirty = true
}

func insertChars(chars ...rune) {
	if len(chars) == 0 {
		return
	}

	x := ActiveCursor.AbsoluteX()
	y := ActiveCursor.AbsoluteY()

	if y >= len(Lines) {
		newLines := make([]string, y+1)
		copy(newLines, Lines)
		Lines = newLines
	}

	line := Lines[y]
	lineRunes := []rune(line)

	if x > len(lineRunes) {
		x = len(lineRunes)
	}

	newLineRunes := make([]rune, 0, len(lineRunes)+len(chars))
	newLineRunes = append(newLineRunes, lineRunes[:x]...)
	newLineRunes = append(newLineRunes, chars...)
	newLineRunes = append(newLineRunes, lineRunes[x:]...)

	Lines[y] = string(newLineRunes)
	ActiveCursor.SetAbsX(x + len(chars))

	IsDirty = true
}

func deleteChar() {
	x := ActiveCursor.AbsoluteX()
	y := ActiveCursor.AbsoluteY()

	if x == 0 && y == 0 {
		return
	}

	current := Lines[y]

	if x == 0 {
		previous := Lines[y-1]

		Lines[y-1] = previous + current

		Lines = append(Lines[:y], Lines[y+1:]...)

		ActiveCursor.Up()
		ActiveCursor.MoveEndOfLine()
	} else {
		before := current[:x-1]
		after := current[x:]

		Lines[y] = before + after

		ActiveCursor.Left()
	}

	IsDirty = true
}
