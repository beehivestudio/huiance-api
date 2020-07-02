package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

/******************************************************************************
 **函数名称: AesEncrypt
 **功    能: AES加密
 **输入参数:
 **     orig: 被加密处理的字串
 **     key: 加密秘钥
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.09.28 16:53:28 #
 ******************************************************************************/
func AesEncrypt(orig, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	blockSize := block.BlockSize()
	orig = PKCS5Padding(orig, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(orig))
	blockMode.CryptBlocks(crypted, orig)

	return crypted, nil
}

/******************************************************************************
 **函数名称: AesEncrypt
 **功    能: AES解密
 **输入参数:
 **     crypted: 被加密处理的字串
 **     key: 加密秘钥
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.09.28 16:55:15 #
 ******************************************************************************/
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	defer func() {
		recover()
	}()

	block, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	orig := make([]byte, len(crypted))
	blockMode.CryptBlocks(orig, crypted)
	orig = PKCS5UnPadding(orig)

	return orig, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(orig []byte) []byte {
	length := len(orig)
	unpadding := int(orig[length-1])
	return orig[:(length - unpadding)]
}
