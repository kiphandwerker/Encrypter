# 🔐 Go File Encryptor
A simple CLI tool written in Go for securely encrypting and decrypting a file using AES-GCM encryption with password-based key derivation (PBKDF2 with SHA-256).

## Motivation
I am managing a few API keys and various other connection strings and storing them in plain text in various places for convenience is obviously not a great idea. I wrote this to store everything in one place but to only be accessible by certain parties.

## ✨ Features
- Secure encryption with AES-256-GCM
- Password-based key derivation (PBKDF2)
- Salt and nonce generation for each encryption
- Minimal and easy-to-use CLI interface

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

```bash
./encryptor encrypt -in api.txt -out encrypted.bin -password "yourPassword123"
```

- in: Path to the file containing the API key.
- out: (Optional) Output file for the encrypted data (default: encrypted.bin).
- password: Password used to encrypt the data (required).

## Decrypting a file
```bash
./encryptor decrypt -in encrypted.bin -password "yourPassword123"
```

- in: Path to the encrypted file.
- password: Password used during encryption (required).

## How It Works
- PBKDF2 is used with 100,000 iterations and SHA-256 to derive a 256-bit AES key from your password and a randomly generated salt.

- AES-GCM is used for authenticated encryption with a randomly generated nonce.

- The final encrypted file format: [salt (16 bytes)] + [nonce (12 bytes)] + [ciphertext].