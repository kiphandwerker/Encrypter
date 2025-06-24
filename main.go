package main

func main(){

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
