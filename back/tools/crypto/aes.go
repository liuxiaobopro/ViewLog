package crypto

import (
	commonCrypto "ViewLog/back/common/crypto"
	"ViewLog/back/global"
)

// AesEncrypt AES加密
func AesEncrypt(s string) (string, error) {
	aes := commonCrypto.AES{
		Key: global.Conf.Aes.Key,
		IV:  global.Conf.Aes.IV,
	}
	return aes.Encrypt(s)
}

// AesDecrypt AES解密
func AesDecrypt(s string) (string, error) {
	aes := commonCrypto.AES{
		Key: global.Conf.Aes.Key,
		IV:  global.Conf.Aes.IV,
	}
	return aes.Decrypt(s)
}
