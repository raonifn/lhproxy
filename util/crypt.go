package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
)

const CRYPTOR_KEY_SIZE = 32
const CRYPTOR_BLOCK_SIZE = aes.BlockSize

func pkcs5pad(plaintext []byte, blockSize int) []byte {
	plaintextLen := len(plaintext)
	padLen := blockSize - (plaintextLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(plaintext, padText...)
}

func pkcs5trim(ciphertext []byte, blockSize int) []byte {
	ciphertextLen := len(ciphertext)
	paddingLen := int(ciphertext[ciphertextLen-1])
	if paddingLen >= ciphertextLen || paddingLen > blockSize {
		log.Panicf("Wrong padding. blockSize: %d, paddingLen: %d, ciphertextLen: %d", blockSize, paddingLen, ciphertextLen)
	}
	return ciphertext[:ciphertextLen-paddingLen]
}

func CryptGen(size int) []byte {
	key := make([]byte, size)
	n, err := rand.Read(key)
	Check(err)
	if n != size {
		log.Panicf("wrong: %d, expected: %d", n, CRYPTOR_KEY_SIZE)
	}
	return key
}

type Cryptor struct {
	Secret []byte
}

func (c *Cryptor) GenSecret() []byte {
	c.Secret = CryptGen(CRYPTOR_KEY_SIZE)
	return c.Secret
}

func (c *Cryptor) GenIV() []byte {
	return CryptGen(CRYPTOR_BLOCK_SIZE)
}

func (c *Cryptor) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(c.Secret)
	Check(err)
	salt := CryptGen(CRYPTOR_BLOCK_SIZE)
	iv := SHA256(salt)[:CRYPTOR_BLOCK_SIZE]
	plaintext = append(salt, plaintext...)
	encrypter := cipher.NewCBCEncrypter(block, iv)
	padded := pkcs5pad(plaintext, encrypter.BlockSize())
	encrypter.CryptBlocks(padded, padded)
	return padded
}

func (c *Cryptor) Decrypt(ciphertext []byte) []byte {
	block, err := aes.NewCipher(c.Secret)
	Check(err)
	salt := ciphertext[:CRYPTOR_BLOCK_SIZE]
	iv := SHA256(salt)[:CRYPTOR_BLOCK_SIZE]
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(ciphertext, ciphertext)
	trimmed := pkcs5trim(ciphertext, decrypter.BlockSize())
	if !bytes.Equal(trimmed[:CRYPTOR_BLOCK_SIZE], salt) {
		log.Panicf("It is wrong salted")
	}
	trimmed = trimmed[CRYPTOR_BLOCK_SIZE:]
	return trimmed
}
