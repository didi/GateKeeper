package public

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"github.com/e421083458/golang_common/lib"
)

//PKCS5Padding padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding unpadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AesEncrypt encode
func AesEncrypt(c context.Context, origData, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			ContextWarning(c, lib.DLTagUndefind, map[string]interface{}{"msg": "AesEncrypt.recover", "err": err})
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AesDecrypt aes decode
func AesDecrypt(c context.Context, crypted, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			ContextWarning(c, lib.DLTagUndefind, map[string]interface{}{"msg": "AesEncrypt.recover", "err": err})
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
