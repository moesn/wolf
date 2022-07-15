package passwd

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"github.com/sirupsen/logrus"
)

type AesEncryptor struct {
	Key    []byte
	Iv     []byte
	Block  cipher.Block
	Base64 bool
}

func (a *AesEncryptor) AesEncrypt(plainText string) string {
	origData := []byte(plainText)
	origData = PKCS5Padding(origData, a.Block.BlockSize())
	crypted := make([]byte, len(origData))

	defer func() {
		if err := recover(); err != nil {
			logrus.Error("AES加密错误:", err)
		}
	}()

	bm := cipher.NewCBCEncrypter(a.Block, a.Iv)
	bm.CryptBlocks(crypted, origData)

	if a.Base64 {
		return base64.StdEncoding.EncodeToString(crypted)
	} else {
		return hex.EncodeToString(crypted)
	}
}

func (a *AesEncryptor) AesDecrypt(cryptText string, errMsg ...string) string {
	crypted := make([]byte, 0)

	if a.Base64 {
		crypted, _ = base64.StdEncoding.DecodeString(cryptText)
	} else {
		crypted, _ = hex.DecodeString(cryptText)
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Error("AES解密错误:", errMsg, err)
		}
	}()

	origData := make([]byte, len(crypted))
	bm := cipher.NewCBCDecrypter(a.Block, a.Iv)
	bm.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)

	return string(origData)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
