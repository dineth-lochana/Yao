package main

type WrappedLine struct {
	Original  string 
	Wrapped   []string 
	StartIdx  []int    
	LineCount int      
}

type WrappedPosition struct {
	OriginalLine  int 
	WrappedLine   int 
	OriginalCol   int 
	WrappedCol    int 
}

func wrapLine(line string, width int) WrappedLine {
	if width <= 0 {
		return WrappedLine{
			Original:  line,
			Wrapped:   []string{line},
			StartIdx:  []int{0},
			LineCount: 1,
		}
	}

	runes := []rune(line)
	var wrapped []string
	var startIdx []int
	lineStart := 0
	lastSpace := -1
	lineLength := 0

	for i, r := range runes {
		if r == ' ' {
			lastSpace = i
		}
		lineLength++

		if lineLength >= width {
			if lastSpace != -1 {
				wrapped = append(wrapped, string(runes[lineStart:lastSpace]))
				startIdx = append(startIdx, lineStart)
				lineStart = lastSpace + 1
				lineLength = i - lastSpace
				lastSpace = -1
			} else {
				wrapped = append(wrapped, string(runes[lineStart:i]))
				startIdx = append(startIdx, lineStart)
				lineStart = i
				lineLength = 0
			}
		}
	}

	if lineStart < len(runes) {
		wrapped = append(wrapped, string(runes[lineStart:]))
		startIdx = append(startIdx, lineStart)
	}

	if len(wrapped) == 0 {
		wrapped = append(wrapped, line)
		startIdx = append(startIdx, 0)
	}

	return WrappedLine{
		Original:  line,
		Wrapped:   wrapped,
		StartIdx:  startIdx,
		LineCount: len(wrapped),
	}
}

func wrappedToOriginal(wrappedLines []WrappedLine, wrappedY, wrappedX int) (int, int) {
	if len(wrappedLines) == 0 {
		return 0, 0
	}

	currentLine := 0
	remainingWrapped := wrappedY

	for i, wl := range wrappedLines {
		if remainingWrapped < wl.LineCount {
			currentLine = i
			break
		}
		remainingWrapped -= wl.LineCount
		if i == len(wrappedLines)-1 {
			currentLine = i
			remainingWrapped = wl.LineCount - 1
			if remainingWrapped < 0 {
				remainingWrapped = 0
			}
		}
	}

	if currentLine >= len(wrappedLines) {
		lastLine := len(wrappedLines) - 1
		if lastLine < 0 {
			return 0, 0
		}
		return lastLine, len([]rune(wrappedLines[lastLine].Original))
	}

	wl := wrappedLines[currentLine]
	if len(wl.StartIdx) == 0 {
		return currentLine, wrappedX
	}

	if remainingWrapped >= len(wl.StartIdx) {
		remainingWrapped = len(wl.StartIdx) - 1
	}

	origX := wl.StartIdx[remainingWrapped] + wrappedX

	if origX > len([]rune(wl.Original)) {
		origX = len([]rune(wl.Original))
	}

	return currentLine, origX
}


func originalToWrapped(wrappedLines []WrappedLine, origY, origX int) (int, int) {
	if len(wrappedLines) == 0 || origY >= len(wrappedLines) {
		return 0, 0
	}

	wrappedY := 0
	for i := 0; i < origY; i++ {
		wrappedY += wrappedLines[i].LineCount
	}

	wl := wrappedLines[origY]
	
	if len(wl.StartIdx) <= 1 {
		return wrappedY, origX
	}

	for i := 0; i < len(wl.StartIdx); i++ {
		nextIdx := len(wl.Original)
		if i+1 < len(wl.StartIdx) {
			nextIdx = wl.StartIdx[i+1]
		}

		if origX >= wl.StartIdx[i] && origX < nextIdx {
			return wrappedY + i, origX - wl.StartIdx[i]
		}
	}

	lastWrapIdx := len(wl.StartIdx) - 1
	return wrappedY + lastWrapIdx, len([]rune(wl.Wrapped[lastWrapIdx]))
}
