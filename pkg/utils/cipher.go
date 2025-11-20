package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MD5Sum(data []byte) []byte {
	result := md5.Sum(data)
	return result[:]
}

func MD5Hex(data []byte) string {
	return hex.EncodeToString(MD5Sum(data))
}

func SHA256Sum(data []byte) []byte {
	result := sha256.Sum256(data)
	return result[:]
}

func SHA256Hex(data []byte) string {
	return hex.EncodeToString(SHA256Sum(data))
}

func AES_256_CFBEncrypt(key []byte, dataIn []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV length")
	}
	ciphertext := make([]byte, len(dataIn))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, dataIn)
	return ciphertext, nil
}

func AES_256_CFBDecrypt(key []byte, dataIn []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV length")
	}
	plaintext := make([]byte, len(dataIn))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, dataIn)
	return plaintext, nil
}
