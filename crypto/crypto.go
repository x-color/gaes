package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func Encrypt(key, plainText []byte) ([]byte, error) {
	plainText = padding(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func Decrypt(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipher text must be longer than blocksize")
	} else if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("cipher text must be multiple of blocksize(128bit)")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	plainText := make([]byte, len(cipherText))

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plainText, cipherText)

	return plainText, nil
}

func padding(b []byte) []byte {
	size := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(size)}, size)
	return append(b, pad...)
}

func unpadding(b []byte) []byte {
	size := int(b[len(b)-1])
	return b[:len(b)-size]
}

func ReadPwd() ([]byte, error) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	originalTerminalState, err := terminal.GetState(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	go func() {
		<-signalChan
		terminal.Restore(int(syscall.Stdin), originalTerminalState)
		os.Exit(1)
	}()

	return terminal.ReadPassword(syscall.Stdin)
}

func PwdToKey(pw []byte) []byte {
	h := sha256.New()
	h.Write(pw)
	return h.Sum(nil)
}
