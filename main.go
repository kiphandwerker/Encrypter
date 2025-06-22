package main

func main(){

}

func Encrypt(apiKey []byte, password string, outputFile string){
  salt := make([]byte, 16)
  _, err := rand.Read(salt)
  if err != nil{
    return err
  }

  key := pbkdf2.Key([]byte(pasword), salt, 100_100, 32, sha256.New)
  block, err := aes.NewCipher(key)
  if err != nil{
    return err
  }

  aes.GCM,err := cipher.NewGCM(block)
  if err != nil{
    return err
  }

  nonce := make([]byte, aesGCM.NonceSize())
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil{
    return err
  }

  ciphertext := aesGCM.Seal(nil, nonce, bytes.TrimSpace(apiKey), nil)

  f, err := os.Create(outputFile)
  if err != nil{
    return err
  }

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
