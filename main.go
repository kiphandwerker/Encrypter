package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: [encrypt|decrypt] [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
		inFile := encryptCmd.String("in", "", "Input file containing the API key")
		outFile := encryptCmd.String("out", "encrypted.bin", "Output encrypted file")
		password := encryptCmd.String("password", "", "Password to encrypt with")

		encryptCmd.Parse(os.Args[2:])

		if *inFile == "" {
			fmt.Println("encrypt: -in is required")
			encryptCmd.Usage()
			os.Exit(1)
		}

		// Prompt for password if not provided
		var pass string
		if *password == "" {
			pass = getPassword("Enter password for encryption: ")
		} else {
			pass = *password
		}

		apiKey, err := os.ReadFile(*inFile)
		ErrorCheck(err)

		err = Encrypt(apiKey, pass, *outFile)
		ErrorCheck(err)

		fmt.Println("✅ Encrypted and saved to", *outFile)

	case "decrypt":
		decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
		inFile := decryptCmd.String("in", "", "Input encrypted file")
		password := decryptCmd.String("password", "", "Password to decrypt with")

		decryptCmd.Parse(os.Args[2:])

		if *inFile == "" {
			fmt.Println("decrypt: -in is required")
			decryptCmd.Usage()
			os.Exit(1)
		}

		// Prompt for password if not provided
		var pass string
		if *password == "" {
			pass = getPassword("Enter password for decryption: ")
		} else {
			pass = *password
		}

		apiKey := Decrypt(*inFile, pass)
		fmt.Println("🔓 Decrypted API File: \n", string(apiKey))

	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: [encrypt|decrypt] [options]")
		os.Exit(1)
	}
}

// getPassword prompts the user for a password without echoing it to the terminal
func getPassword(prompt string) string {
	fmt.Print(prompt)
	// Use term.ReadPassword to hide input
	pass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error reading password:", err)
		os.Exit(1)
	}
	fmt.Println() // Add newline after password input
	return string(pass)
}

func Encrypt(apiKey []byte, password string, outputFile string) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	ErrorCheck(err)

	key := pbkdf2.Key([]byte(password), salt, 100_000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	ErrorCheck(err)

	aesGCM, err := cipher.NewGCM(block)
	ErrorCheck(err)

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	ErrorCheck(err)

	ciphertext := aesGCM.Seal(nil, nonce, bytes.TrimSpace(apiKey), nil)

	f, err := os.Create(outputFile)
	ErrorCheck(err)
	defer f.Close()

	f.Write(salt)
	f.Write(nonce)
	f.Write(ciphertext)

	return nil
}

func Decrypt(filename string, password string) []byte {
	data, err := os.ReadFile(filename)
	ErrorCheck(err)

	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	key := pbkdf2.Key([]byte(password), salt, 100_000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	ErrorCheck(err)

	aesGCM, err := cipher.NewGCM(block)
	ErrorCheck(err)

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	ErrorCheck(err)

	return plaintext
}

func ErrorCheck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error:", err)
		os.Exit(1)
	}
}
