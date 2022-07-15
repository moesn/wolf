package passwd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/sirupsen/logrus"
	"os"
)

type RsaEncryptor struct {
	PrivateKey []byte
	PublickKey []byte
	Base64     bool
}

func (r *RsaEncryptor) GenerateRSAKey(bits int) {
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

func (r *RsaEncryptor) RsaEncrypt(plainText string) string {
	block, _ := pem.Decode(r.PublickKey)

	defer func() {
		if err := recover(); err != nil {
			logrus.Error("RSA加密错误:", err)
		}
	}()

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))

	if err != nil {
		return ""
	}

	if r.Base64 {
		return base64.StdEncoding.EncodeToString(cipherText)
	} else {
		return hex.EncodeToString(cipherText)
	}
}

func (r *RsaEncryptor) RsaDecrypt(cryptText string, errMsg ...string) string {
	block, _ := pem.Decode(r.PrivateKey)
	crypted := make([]byte, 0)

	defer func() {
		if err := recover(); err != nil {
			logrus.Error("RSA解密错误:", errMsg, err)
		}
	}()

	if r.Base64 {
		crypted, _ = base64.StdEncoding.DecodeString(cryptText)
	} else {
		crypted, _ = hex.DecodeString(cryptText)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, crypted)

	if err != nil {
		return ""
	}
	return string(plainText)
}
