package main

import (
	"bufio"
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
	"unicode"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
)

const minPasswordLength = 8

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
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
			pass = getPasswordWithConfirmation("Enter password for encryption: ")
		} else {
			pass = *password
		}

		if len(pass) < minPasswordLength {
			fmt.Fprintf(os.Stderr, "❌ Error: password must be at least %d characters\n", minPasswordLength)
			os.Exit(1)
		}

		label, color := passwordStrength(pass)
		fmt.Println("Password strength:", colorize(label, color))

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

// getPassword prompts for a password, echoing '*' per character typed.
// Backspace/Delete edits the buffer; Ctrl+C aborts the program.
func getPassword(prompt string) string {
	fmt.Print(prompt)

	fd := int(syscall.Stdin)
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error reading password:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	var buf []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			term.Restore(fd, oldState)
			fmt.Println()
			os.Exit(1)
		}

		switch b {
		case '\r', '\n':
			term.Restore(fd, oldState)
			fmt.Println()
			return string(buf)
		case 3: // Ctrl+C
			term.Restore(fd, oldState)
			fmt.Println()
			fmt.Fprintln(os.Stderr, "❌ Cancelled")
			os.Exit(1)
		case 127, 8: // Backspace / Delete
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
		default:
			buf = append(buf, b)
			fmt.Print("*")
		}
	}
}

// colorize wraps text in an ANSI color code, but only if stdout is a terminal
// (so redirected/piped output stays plain).
func colorize(text, color string) string {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return text
	}
	return color + text + colorReset
}

// passwordStrength scores pw by length tiers and character-class variety.
// It assumes pw already meets minPasswordLength.
func passwordStrength(pw string) (label string, color string) {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, c := range pw {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	score := 1 // meets the minimum length requirement
	switch {
	case len(pw) >= 16:
		score += 2
	case len(pw) >= 12:
		score++
	}
	if hasLower && hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	switch {
	case score >= 5:
		return "Strong", colorGreen
	case score == 4:
		return "Good", colorCyan
	case score == 3:
		return "Fair", colorYellow
	default:
		return "Weak", colorRed
	}
}

// getPasswordWithConfirmation prompts for a password twice and exits if they don't match
func getPasswordWithConfirmation(prompt string) string {
	pass := getPassword(prompt)
	confirm := getPassword("Confirm password: ")
	if pass != confirm {
		fmt.Fprintln(os.Stderr, "❌ Error: passwords do not match")
		os.Exit(1)
	}
	return pass
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
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Incorrect password")
		os.Exit(1)
	}

	return plaintext
}

func ErrorCheck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error:", err)
		os.Exit(1)
	}
}
