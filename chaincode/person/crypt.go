package main

// Create a Merke-Damgard MD5 checksum, hex encoded
// the output is not going to be stored so we do not care about the MSD5 insecurity
func hashMD5(key []byte) string {
	hasher := md5.New()
	hasher.Write(key)

	// 128-bit long
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt bytes stream: account data, image, document data
func encrAESGCM(data []byte, key[]byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		retutn nil, err
	}

	gcmCipher := gcm.Seal(nonce, nonce, data, nil)
	return gcmCipher, nil
}

func decrAESGCM(data []byte, key []byte) ([]byte, error) {
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