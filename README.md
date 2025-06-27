# 🔐 Go File Encryptor
A simple command line interface (CLI) tool written in Go for securely encrypting and decrypting a file using AES-GCM encryption with password-based key derivation (PBKDF2 with SHA-256).

# Table of Contents
- [Motivation](#motivation)
- [How It Works](#how-it-works)
- [Installation](#-installation)
- [Encrypt](#encrypting-a-file)
- [Decrypt](#decrypting-a-file)

## Motivation
I am managing a few API keys and various other connection strings. Storing them in plain text in various places for convenience is obviously not a great idea. I wrote this to store everything in a place of my choosing, and to prevent unwanted parties from accessing the information.

While the original conception of this idea was to encrypt 1 API key, it can be used for an entire file. The results of decryption displays the results of the encrypted.bin file in the terminal.

## How It Works
- [PBKDF2](https://en.wikipedia.org/wiki/PBKDF2) is used with 100,000 iterations and [SHA-256](https://en.wikipedia.org/wiki/SHA-2) to derive a 256-bit AES key from your password and a randomly generated [salt](https://en.wikipedia.org/wiki/Salt_%28cryptography%29).

- [AES-GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode) is used for authenticated encryption with a randomly generated [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce).

- The final encrypted file format: [salt (16 bytes)] + [nonce (12 bytes)] + [ciphertext].

## 📦 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/kiphandwerker/Encrypter.git
    ```

2. Build the binary:
    ```bash
    go build -o encryptor main.go
    ```

## Encrypting a file

```shell
./encryptor encrypt -in /path/to/some/file.txt -out encrypted.bin -password "yourPassword123"
```

- in: Path to the file containing the API key.
- out: (Optional) Output file for the encrypted data (default: encrypted.bin).
- password: Password used to encrypt the data (required).

## Decrypting a file
```shell
./encryptor decrypt -in /path/to/some/encrypted.bin -password "yourPassword123"
```

- in: Path to the encrypted file.
- password: Password used during encryption (required).

