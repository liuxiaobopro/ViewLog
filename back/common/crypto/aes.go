package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type AES struct {
	Key string // 密钥
	IV  string // 偏移量
}

// AES 加密
func (th *AES) Encrypt(plainText string) (string, error) {
	// 创建加密算法aes
	block, err := aes.NewCipher([]byte(th.Key))
	if err != nil {
		return "", err
	}

	// 补全码
	blockSize := block.BlockSize()
	plaintext := PKCS5Padding([]byte(plainText), blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(th.IV))
	// 创建密文长度的数组
	ciphertext := make([]byte, len(plaintext))
	// 加密明文
	blockMode.CryptBlocks(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES 解密
func (th *AES) Decrypt(ciphertext string) (string, error) {
	// 创建加密算法aes
	block, err := aes.NewCipher([]byte(th.Key))
	if err != nil {
		return "", err
	}

	// 解密模式
	blockMode := cipher.NewCBCDecrypter(block, []byte(th.IV))
	// base64解密
	ciphertextByte, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	// 创建密文长度的数组
	plaintext := make([]byte, len(ciphertextByte))
	// 解密密文
	blockMode.CryptBlocks(plaintext, ciphertextByte)
	// 去除补全码
	plaintext, err = PKCS5UnPadding(plaintext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// PKCS5Padding 补全码
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// PKCS5UnPadding 去除补全码
func PKCS5UnPadding(plaintext []byte) ([]byte, error) {
	length := len(plaintext)
	// 最后一个字节为 padding 的长度
	unpadding := int(plaintext[length-1])
	if unpadding > length {
		return nil, errors.New("PKCS5 unpadding failed: padding size too big")
	}
	return plaintext[:(length - unpadding)], nil
}
