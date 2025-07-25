package main

type Cursor struct {
	Disabled bool
	X int
	Y int
	BodyWidth  int
	BodyHeight int
}

var WrappedLines []WrappedLine

func (c *Cursor) SetDimensions(bodyWidth, bodyHeight int) {
	c.BodyWidth = bodyWidth
	c.BodyHeight = bodyHeight
	
	// Recalculate word wrapping when dimensions change
	WrappedLines = make([]WrappedLine, len(Lines))
	for i, line := range Lines {
		WrappedLines[i] = wrapLine(line, bodyWidth-4) // Account for margins
	}
}

func (c *Cursor) AbsolutePosition() (int, int) {
	return wrappedToOriginal(WrappedLines, c.Y+ScrollY, c.X)
}

func (c *Cursor) SetDisabled(disabled bool) {
	c.Disabled = disabled
}

func (c *Cursor) AbsoluteX() int {
	return c.X + ScrollX
}

func (c *Cursor) AbsoluteY() int {
	return c.Y + ScrollY
}

func (c *Cursor) CurrentLine() string {
	origY, _ := c.AbsolutePosition()
	if origY >= len(Lines) {
		return ""
	}
	
	remainingY := c.Y + ScrollY
	
	// Find the current wrapped line
	for i := 0; i <= origY; i++ {
		if remainingY < WrappedLines[i].LineCount {
			return WrappedLines[i].Wrapped[remainingY]
		}
		remainingY -= WrappedLines[i].LineCount
	}
	
	return ""
}

func (c *Cursor) LineLength() int {
	return len([]rune(c.CurrentLine()))
}

func (c *Cursor) CurrentChar() string {
	absX := c.AbsoluteX()
	line := c.CurrentLine()
	runes := []rune(line)

	if absX >= len(runes) {
		return " "
	}

	return string(runes[absX])
}

// Handles vertical scrolling when cursor reaches screen boundaries
func (c *Cursor) updateVerticalScroll() {
	if c.Y < 0 {
		// Scrolling up
		ScrollY += c.Y
		if ScrollY < 0 {
			ScrollY = 0
		}
		c.Y = 0
	} else if c.Y >= c.BodyHeight {
		// Scrolling down
		ScrollY += c.Y - (c.BodyHeight - 1)
		c.Y = c.BodyHeight - 1
		
		// Calculate total wrapped lines
		totalWrappedLines := 0
		for _, wl := range WrappedLines {
			totalWrappedLines += wl.LineCount
		}
		
		// Prevent scrolling past the last line
		maxScroll := totalWrappedLines - c.BodyHeight
		if ScrollY > maxScroll {
			ScrollY = maxScroll
			if ScrollY < 0 {
				ScrollY = 0
			}
		}
	}
}

// Moves cursor up one line, handles scrolling if needed
func (c *Cursor) Up() {
	if c.Disabled {
		return
	}

	// Don't move up if at the very top of the document
	if c.AbsoluteY() <= 0 {
		return
	}

	c.Y--
	c.updateVerticalScroll()
}

// Moves cursor down one line, handles scrolling if needed
func (c *Cursor) Down() {
	if c.Disabled {
		return
	}

	// Calculate total wrapped lines
	totalWrappedLines := 0
	for _, wl := range WrappedLines {
		totalWrappedLines += wl.LineCount
	}

	// Don't move down if at the very bottom of the document
	if c.AbsoluteY() >= totalWrappedLines-1 {
		return
	}

	c.Y++
	c.updateVerticalScroll()
}

// SetAbsY now only handles setting the Y position and triggers scroll updates
func (c *Cursor) SetAbsY(y int) {
	// Ensure y is not negative
	if y < 0 {
		y = 0
	}

	origY, origX := wrappedToOriginal(WrappedLines, y, c.X)
	
	// Ensure we don't go beyond the last line
	if origY >= len(Lines) {
		origY = len(Lines) - 1
		if origY < 0 {
			origY = 0
		}
		origX = 0
	}
	
	// Convert back to wrapped coordinates
	wrappedY, wrappedX := originalToWrapped(WrappedLines, origY, origX)
	
	c.Y = wrappedY - ScrollY
	c.X = wrappedX
	
	c.updateVerticalScroll()
}

// Existing horizontal movement functions remain unchanged
func (c *Cursor) SetAbsX(x int) {
	ScrollX = 0
	c.X = x

	if c.X < 0 {
		ScrollX += c.X
		if ScrollX < 0 {
			ScrollX = 0
		}
		c.X = 0
	} else {
		maxLine := c.LineLength()
		maxX := maxLine

		if maxX >= c.BodyWidth {
			maxX = c.BodyWidth - 1
		}

		tooFar := c.X - maxX
		if tooFar > 0 {
			ScrollX += tooFar
			c.X = maxX

			if ScrollX > maxLine-maxX {
				ScrollX = maxLine - maxX
				c.X = maxX
			}
		}
	}
}

func (c *Cursor) Left() {
	if c.Disabled {
		return
	}

	if c.AbsoluteX() == 0 {
		ay := c.AbsoluteY()
		if ay > 0 {
			c.Up()
			c.MoveEndOfLine()
		}
		return
	}

	c.SetAbsX(c.AbsoluteX() - 1)
}

func (c *Cursor) Right() {
	if c.Disabled {
		return
	}

	if c.AbsoluteX() == c.LineLength() {
		if c.AbsoluteY() < len(Lines)-1 {
			c.MoveStartOfLine()
			c.Down()
		}
		return
	}

	c.SetAbsX(c.AbsoluteX() + 1)
}

func (c *Cursor) MoveEndOfLine() {
	c.SetAbsX(c.LineLength())
}

func (c *Cursor) MoveStartOfLine() {
	c.SetAbsX(0)
}
