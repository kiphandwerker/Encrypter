package main

func main(){

}

func Encrypt(apiKey []byte, password string, outputFile string){
  salt := make([]byte, 16)
  _ := rand.Read(salt)

  key := pbkdf2.Key([]byte(pasword), salt, 100_100, 32, sha256.New)
  block := aes.NewCipher(key)

  aes.GCM := cipher.NewGCM(block)

  nonce := make([]byte, aesGCM.NonceSize())

  ciphertext := aesGCM.Seal(nil, nonce, bytes.TrimSpace(apiKey), nil)

  f := os.Create(outputFile)

  defer f.close()

  f.Write(salt)
  f.Write(nonce)
  f.Write(ciphertext)

  return nil

}

func Decrypt(){


}

func ErrorCheck(){


}
