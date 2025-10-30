# 🔐 Go File Encryptor
A simple command line interface (CLI) tool written in Go for securely encrypting and decrypting a file using AES-GCM encryption with password-based key derivation (PBKDF2 with SHA-256).

# Table of Contents
- [Motivation](#motivation)
- [How It Works](#how-it-works)
- [Installation](#-installation)
- [Encrypt](#encrypting-a-file)
- [Decrypt](#decrypting-a-file)
- [Helpful tip](#helpful-tip)

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
    go build -o encryptor.exe main.go
    ```

## Encrypting a file

```shell
encryptor.exe encrypt -in /path/to/some/file.txt -out encrypted.bin
```

```
Enter password for encryption: 
```

- in: Path to the file containing the API key.
- out: (Optional) Output file for the encrypted data (default: encrypted.bin).

## Decrypting a file
```shell
encryptor.exe decrypt -in /path/to/some/encrypted.bin
```

```
Enter password for decryption: 
```

- in: Path to the encrypted bin.

# Helpful tip

It is annoying to have to save and move the .exe and .bin around and keep up with it. So one thing that can make this less of a hassle to is put it in one place and create an alias or similar to run it.

I use a simple call from my .bashrc to do just that.
<ul>
  <li>dkeys:       Specifies the name to call in the terminal</li>
  <li>exe:         Location of the .exe</li>
  <li>infile:      Location of the encrypted file</li>
  <li>decrypt -in: Specifies we will be decrypting the target file</li>
</ul>

```bash
dkeys() {
  local exe="/mnt/c/Users/...../encryptor.exe"

  local infile="C:\\Users\\.....\\encrypted.bin"

  "$exe" decrypt -in "$infile"
}
```
```

