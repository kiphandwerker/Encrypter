package main

func main(){
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

      if *inFile == "" || *password == "" {
        fmt.Println("encrypt: -in and -password are required")
        encryptCmd.Usage()
        os.Exit(1)
      }

      apiKey, err := os.ReadFile(*inFile)
      check(err)

      err = Encrypt(apiKey, *password, *outFile)
      check(err)

      fmt.Println("✅ Encrypted and saved to", *outFile)

    case "encrypt":
      decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
      inFile := decryptCmd.String("in", "", "Input encrypted file")
      password := decryptCmd.String("password", "", "Password to decrypt with")

      decryptCmd.Parse(os.Args[2:])

      if *inFile == "" || *password == "" {
        fmt.Println("decrypt: -in and -password are required")
        decryptCmd.Usage()
        os.Exit(1)
      }

      apiKey := decryptFile(*inFile, *password)
      fmt.Println("🔓 Decrypted API Key:", string(apiKey))

    default:
      fmt.Println("Unknown command:", os.Args[1])
      fmt.Println("Usage: [encrypt|decrypt] [options]")
      os.Exit(1)
	}

}

func Encrypt(apiKey []byte, password string, outputFile string){
  salt := make([]byte, 16)
  _, err := rand.Read(salt)
  ErrorCheck(err)

  key := pbkdf2.Key([]byte(password), salt, 100_100, 32, sha256.New)
  block, err := aes.NewCipher(key)
  ErrorCheck(err)

  aes.GCM,err := cipher.NewGCM(block)
  ErrorCheck(err)

  nonce := make([]byte, aesGCM.NonceSize())
  _, err := io.ReadFull(rand.Reader, nonce)
  ErrorCheck(err)

  ciphertext := aesGCM.Seal(nil, nonce, bytes.TrimSpace(apiKey), nil)

  f, err := os.Create(outputFile)
  ErrorCheck(err)

  defer f.close()

  f.Write(salt)
  f.Write(nonce)
  f.Write(ciphertext)

  return nil

}

func Decrypt(){
	data, err := os.ReadFile(filename)
	check(err)

	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	key := pbkdf2.Key([]byte(password), salt, 100_000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	check(err)

	aesGCM, err := cipher.NewGCM(block)
	check(err)

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	check(err)

	return plaintext

}

func ErrorCheck(err error){
  if err != nil {
    fmt.FPrint(os.Stderr, "Error: ", err)
    os.Exit(1)
  }

}
