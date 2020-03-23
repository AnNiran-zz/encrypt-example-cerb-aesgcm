package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func createCipherKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return key, nil
}

func SaveCipherKey(cipherKey []byte, path, version string) error {

	cipherKeyPath := filepath.Join(path, version)
	if err := os.MkdirAll(cipherKeyPath, 0755); err != nil {
		return err
	}

	file, err := os.Create(cipherKeyPath + "/cipher")
	if err != nil {
		return err
	}

	if _, err = file.Write(cipherKey); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func DeleteCipherKeyFile(path, version string) error {

	cipherKeyPath := filepath.Join(path, version, "cipher")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("Path to cipher key location does not exist")
	}

	if _, err = os.Stat(cipherKeyPath); os.IsNotExist(err) {
		return nil
	}

	directory, err := ioutil.ReadDir(cipherKeyPath)
	if err != nil {
		return err
	}

	for _, item := range directory {
		os.RemoveAll(filepath.Join([]string{cipherKeyPath, item.Name()}...))
	}

	os.Remove(cipherKeyPath)

	return nil
}

func ReadCipherKey(cipherPath string) ([]byte, error) {

	filename := cipherPath + "/cipher"
	cipherKey, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return cipherKey, nil
}
