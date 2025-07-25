
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "io"
    "os"
    "strings"
    "crypto/rand"
    "github.com/inancgumus/screen"
)


func encrypt(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    return gcm.Seal(nonce, nonce, data, nil), nil
}

func save(password string) {
    contents := strings.Join(Lines, "\n")
    var data []byte
    var err error

    if password != "" {
        key := deriveKey(password)
        data, err = encrypt([]byte(contents), key)
        throw("Unable to encrypt data.", err)
    } else {
        data = []byte(contents)
    }

    err = File.Truncate(0)
    throw("Unable to truncate file.", err)

    _, err = File.Seek(0, 0)
    throw("Unable to seek to beginning of file.", err)

    _, err = File.Write(data)
    throw("Unable to write to file.", err)

    IsDirty = false
}

func exit() {
    ShowCursor()
    screen.MoveTopLeft()
    screen.Clear()
    _ = File.Close()
    os.Exit(0)
}
