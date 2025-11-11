package tmhelper

import (
	"crypto/aes"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

var _salt = []byte{37, 112, 39, 97, 86, 35, 118, 22, 43, 78, 111, 123, 17, 48, 19, 29}

func GenKey(key[] byte,size int)([] byte){
    return pbkdf2.Key(key, _salt, 15, size, sha256.New)
}
func GenKeyX(key string,size int)([] byte){
    return GenKey([]byte(key),size)
}

func AesEnc(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(key)
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func AesDec(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(key)
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
