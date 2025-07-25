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

	case keyboard.KeyCtrlQ:
		jumpLines(-10)
	case keyboard.KeyCtrlW:
		jumpLines(10)

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
	origY, origX := ActiveCursor.AbsolutePosition()

	if origY >= len(Lines) {
		origY = len(Lines) - 1
		if origY < 0 {
			origY = 0
			Lines = []string{""}
		}
	}

	currentLine := Lines[origY]

	if origX > len(currentLine) {
		origX = len(currentLine)
	}

	firstHalf := currentLine[:origX]
	secondHalf := currentLine[origX:]

	newLines := make([]string, 0, len(Lines)+1)

	newLines = append(newLines, Lines[:origY]...)

	newLines = append(newLines, firstHalf, secondHalf)

	if origY+1 < len(Lines) {
		newLines = append(newLines, Lines[origY+1:]...)
	}

	Lines = newLines

	ActiveCursor.SetDimensions(ActiveCursor.BodyWidth, ActiveCursor.BodyHeight)

	ActiveCursor.SetAbsY(ActiveCursor.AbsoluteY() + 1)
	ActiveCursor.SetAbsX(0)

	IsDirty = true
}

func insertChars(chars ...rune) {
	origY, origX := ActiveCursor.AbsolutePosition()

	if origY >= len(Lines) {
		origY = len(Lines) - 1
		if origY < 0 {
			origY = 0
			Lines = []string{""}
		}
	}

	line := Lines[origY]

	if origX > len(line) {
		origX = len(line)
	}

	newLine := line[:origX] + string(chars) + line[origX:]
	Lines[origY] = newLine

	ActiveCursor.SetDimensions(ActiveCursor.BodyWidth, ActiveCursor.BodyHeight)

	origX += len(chars)

	wrappedY, wrappedX := originalToWrapped(WrappedLines, origY, origX)

	ActiveCursor.Y = wrappedY - ScrollY
	ActiveCursor.X = wrappedX

	IsDirty = true
}

func deleteChar() {
	origY, origX := ActiveCursor.AbsolutePosition()

	if origX == 0 && origY == 0 {
		return
	}

	if origY >= len(Lines) {
		return
	}

	currentLine := Lines[origY]

	if origX == 0 {
		if origY > 0 {
			previousLine := Lines[origY-1]
			Lines[origY-1] = previousLine + currentLine
			Lines = append(Lines[:origY], Lines[origY+1:]...)

			ActiveCursor.SetDimensions(ActiveCursor.BodyWidth, ActiveCursor.BodyHeight)

			wrappedY, wrappedX := originalToWrapped(WrappedLines, origY-1, len(previousLine))

			ActiveCursor.Y = wrappedY - ScrollY
			ActiveCursor.X = wrappedX
		}
	} else {
		if origX > len(currentLine) {
			origX = len(currentLine)
		}

		newLine := currentLine[:origX-1] + currentLine[origX:]
		Lines[origY] = newLine

		ActiveCursor.SetDimensions(ActiveCursor.BodyWidth, ActiveCursor.BodyHeight)

		wrappedY, wrappedX := originalToWrapped(WrappedLines, origY, origX-1)

		ActiveCursor.Y = wrappedY - ScrollY
		ActiveCursor.X = wrappedX
	}

	IsDirty = true
}

func jumpLines(count int) {
	origY, origX := ActiveCursor.AbsolutePosition()
	var newY int

	if count > 0 {

		newY = origY + count
		if newY >= len(WrappedLines) {
			newY = len(WrappedLines) - 1
		}
	} else {

		newY = origY + count
		if newY < 0 {
			newY = 0
		}
	}

	wrappedY, wrappedX := originalToWrapped(WrappedLines, newY, origX)
	ActiveCursor.SetAbsY(wrappedY)
	ActiveCursor.SetAbsX(wrappedX)

	if wrappedY < ScrollY {
		ScrollY = wrappedY
	} else if wrappedY >= ScrollY+ActiveCursor.BodyHeight {
		ScrollY = wrappedY - ActiveCursor.BodyHeight + 1
	}
}
