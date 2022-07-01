package main

import (
	"encoding/hex"
	"fmt"
	"github.com/moesn/wolf/common/passwd"
)

var privateKey = `-----BEGIN 私钥-----
MIIBOgIBAAJBAIBnCBr3eKZ4IckfTYnw7Zy3LunnQHHk6hwO3llcOL+4dyliulcc
NywKuqqaz3UXlk7N4HAK55i4w+HA64Mini0CAwEAAQJAOtn+dw2yOOqa4an9Shop
bkjvsbmSaiF3xAB0nRSN1iE0QTsAmgQzA2tYxjrK5kRzpqX+wpnWzOcLTekLnICp
gQIhAOww9ZzNU+5Ml3dwtjsOP8lLy6iueWvPzHI8fQ1SqzpxAiEAiyvXK9ju0JbI
tTorV1qyAzkHC+2bj/FfVlviET5A5X0CIChU5+wBp95eV9L+fS+ybDP8bGd+e59s
WXX8XEJFJpYRAiAlINEKtZkgtNle3+HoMB4FlhU3LXCu+ioBFAo/p2gQ0QIhAL4E
vjx88Ivu8xhrC6cVHoEOuuehDbTJP7TCDg6UFC1b
-----END 私钥-----
`

var publicKey = `-----BEGIN 公钥-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAIBnCBr3eKZ4IckfTYnw7Zy3LunnQHHk
6hwO3llcOL+4dyliulccNywKuqqaz3UXlk7N4HAK55i4w+HA64Mini0CAwEAAQ==
-----END 公钥-----
`

func main() {

	privateKey := []byte(privateKey)

	message, _ := hex.DecodeString("554aad4a89cf21670d170b2a90ca5befa2dc5dc0e6c1c0bbdfb9e5f1e59451ca64799c23fdc50aaa097a84ea914d319b0baf40bfeec7bc6cdb1164a75c75ade1")

	plainText, _ := passwd.RsaDecrypt(message, privateKey)
	fmt.Println("解密后为：", string(plainText))
}
