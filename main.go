package main

func main(){

}

func Encrypt(apiKey []byte, password string, outputFile string){
  salt := make([]byte, 16)
  _, err := rand.Read(salt)
  ErrorCheck(err)

  key := pbkdf2.Key([]byte(pasword), salt, 100_100, 32, sha256.New)
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


}

func ErrorCheck(err error){
  if err != nil {
    fmt.FPrint(os.Stderr, "Error: ", err)
    os.Exit(1)
  }

}
