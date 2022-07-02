package passwd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"runtime"
)

func GenerateRSAKey(bits int) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}

	defer privateFile.Close()
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	pem.Encode(privateFile, &privateBlock)

	publicKey := privateKey.PublicKey
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}

	defer publicFile.Close()
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	pem.Encode(publicFile, &publicBlock)
}

func RsaEncrypt(plainText, key []byte) (cryptText []byte, err error) {
	block, _ := pem.Decode(key)
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "Check that the key is correct")
			default:
				log.Println("error:", err)
			}
		}
	}()
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

func RsaDecrypt(cryptText, key []byte) (plainText []byte, err error) {
	block, _ := pem.Decode(key)

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "Check that the key is correct")
			default:
				log.Println("error:", err)
			}
		}
	}()

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}

	plainText, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, cryptText)
	if err != nil {
		return []byte{}, err
	}
	return plainText, nil
}
