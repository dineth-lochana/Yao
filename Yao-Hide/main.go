package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
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

type Position struct {
	X, Y int
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func promptPassword(prompt string) string {
	fmt.Print(prompt)
	var pass string
	fmt.Scanln(&pass)
	return pass
}

var password string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: Cellel-Grimoire <file> [pass=password]")
		return
	}

	path, err := filepath.Abs(os.Args[1])
	throw("Unable to resolve file path.", err)

	if len(os.Args) >= 3 && strings.HasPrefix(os.Args[2], "pass=") {
		password = strings.TrimPrefix(os.Args[2], "pass=")
	}

	FileName = filepath.Base(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If a password was provided and we're creating a new file
		if password != "" {
			// Prompt for password verification
			verifyPass := promptPassword("Verify password: ")
			if verifyPass != password {
				fmt.Println("Passwords do not match. Exiting...")
				os.Exit(1)
			}
		}

		File, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
		throw("Unable to create file.", err)
		Lines = []string{""}
	} else {
		File, err = os.OpenFile(path, os.O_RDWR, 0755)
		throw("Unable to open file.", err)

		b, err := io.ReadAll(File)
		throw("Unable to read file.", err)

		if password != "" {
			key := deriveKey(password)
			b, err = decrypt(b, key)
			throw("Unable to decrypt file contents.", err)
		}

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