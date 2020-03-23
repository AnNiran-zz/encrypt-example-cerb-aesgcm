package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// Create a Merke-Damgard MD5 checksum, hex encoded
// the output is not going to be stored so we do not care about the MSD5 insecurity
func hashMD5(key []byte) string {
	hasher := md5.New()
	hasher.Write(key)

	// 128-bit long
	return hex.EncodeToString(hasher.Sum(nil))
}

// Create a new 32-bit key used for further encryption
func Key32byt() ([]byte, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

// Encrypt bytes stream: account data, image, document data
// create 32-bit encryption key
func EncrAESGCM(data []byte, key) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		retutn nil, err
	}

	gcmCipher := gcm.Seal(nil, nonce, data, nil)
	return gcmCipher, nil
}

// Decrypt data using AESGCM
func DecrAESGCM(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
