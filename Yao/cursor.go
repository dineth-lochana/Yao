package main

type Cursor struct {
	Disabled bool

	X int
	Y int

	BodyWidth  int
	BodyHeight int
}

func (c *Cursor) SetDimensions(bodyWidth, bodyHeight int) {
    // Enforce minimum dimensions
    if bodyWidth < 10 {
        bodyWidth = 10
    }
    if bodyHeight < 3 {
        bodyHeight = 3
    }
    c.BodyWidth = bodyWidth
    c.BodyHeight = bodyHeight
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
    // Check if the absolute Y position is within bounds
    if y := c.AbsoluteY(); y >= 0 && y < len(Lines) {
        return Lines[y]
    }
    return ""
}

func (c *Cursor) LineLength() int {
    line := c.CurrentLine()
    return len([]rune(line))
}

func (c *Cursor) CurrentChar() string {
    absX := c.AbsoluteX()
    line := c.CurrentLine()
    runes := []rune(line)

    if absX < 0 || absX >= len(runes) {
        return " "
    }

    return string(runes[absX])
}

func (c *Cursor) SetAbsX(x int) {
    ScrollX = 0
    c.X = x

    if c.X < 0 {
        ScrollX += c.X
        if ScrollX < 0 {
            ScrollX = 0
        }
        c.X = 0
        return
    }

    maxLine := c.LineLength()
    maxX := maxLine

    if maxX >= c.BodyWidth {
        maxX = c.BodyWidth - 1
    }

    // Ensure cursor doesn't go beyond line length
    if c.X > maxLine {
        c.X = maxLine
    }

    tooFar := c.X - maxX
    if tooFar > 0 {
        ScrollX += tooFar
        c.X = maxX

        if maxLine > maxX && ScrollX > maxLine-maxX {
            ScrollX = maxLine - maxX
            c.X = maxX
        }
    }
}

func (c *Cursor) SetAbsY(y int) {
    // Only move the cursor's Y position within allowed bounds
    maxLines := len(Lines) - 1
    if maxLines < 0 {
        maxLines = 0
    }
    if y < 0 {
        c.Y = 0
    } else if y > maxLines {
        c.Y = maxLines
    } else {
        c.Y = y
    }
}

func (c *Cursor) scrollIfNeeded() {
    // If cursor is at the top of the screen and more lines can be scrolled up
    if c.Y == 0 && ScrollY > 0 {
        ScrollY -= 1 // Scroll up by one line
    } else if c.Y == c.BodyHeight - 1 && ScrollY < len(Lines) - c.BodyHeight {
        // If cursor is at the bottom of the screen and more lines can be scrolled down
        ScrollY += 1 // Scroll down by one line
    }
}

func (c *Cursor) Up() {
    if c.Disabled {
        return
    }
    // Only move the cursor up if not at the top of the visible screen
    if c.Y > 0 {
        c.SetAbsY(c.Y - 1)
    } else if ScrollY > 0 {
        c.scrollIfNeeded()  // Keep scrolling when at the top
    }
}

func (c *Cursor) Down() {
    if c.Disabled {
        return
    }
    // Only move the cursor down if not at the bottom of the visible screen
    if c.Y < c.BodyHeight - 1 {
        c.SetAbsY(c.Y + 1)
    } else if ScrollY < len(Lines) - c.BodyHeight {
        c.scrollIfNeeded()  // Keep scrolling when at the bottom
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

    currentLineLength := c.LineLength()
    if c.AbsoluteX() >= currentLineLength {
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
