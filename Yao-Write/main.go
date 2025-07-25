package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

const Version = "v1.0.0"

var (
	File *os.File
	IsCRLF   = false
	IsDirty  = false
	FileName string
	Lines    []string

	ActiveCursor = Cursor{X: 0, Y: 0}
	ScrollY      = 0
	ScrollX      = 0
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: pico <file>")
		return
	}

	path, err := filepath.Abs(os.Args[1])
	throw("Unable to resolve file path.", err)

	FileName = filepath.Base(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		File, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
		throw("Unable to create file.", err)
		Lines = []string{""}
	} else {
		File, err = os.OpenFile(path, os.O_RDWR, 0755)
		throw("Unable to open file.", err)

		b, err := io.ReadAll(File)
		throw("Unable to read file.", err)

		if bytes.Contains(b, []byte("\r\n")) {
			IsCRLF = true
			b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
		}

		b = bytes.ReplaceAll(b, []byte("\t"), []byte("    "))

		Lines = strings.Split(string(b), "\n")
	}

	loop()
}

func loop() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Cannot setup the keyboard:", err)
		return
	}

	defer keyboard.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		exit()
	}()

	HideCursor()
	screen.Clear()

	for {
		screenWidth, screenHeight := screen.Size()
		ActiveCursor.SetDimensions(screenWidth-4, screenHeight-2)
		draw(Lines)
		readAndHandleKey()
	}
}
