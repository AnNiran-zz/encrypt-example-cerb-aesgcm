package crypto

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ipfsPath string = os.Getenv("GOPATH") + "/src/cerberus/ipfs"

func encryptWithPublicKey(cipherKey []byte, rsaPath string) ([]byte, error) {

	rsaKeysPath := rsaPath + "/rsa_key.pem"
	privateKey, err := parseRsaPrivateKeyFromPemFile(rsaKeysPath)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	cipherKeyEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, cipherKey)
	if err != nil {
		return nil, err
	}

	return cipherKeyEncrypted, nil
}

func decryptWitPrivateKey(encryptedCipherKey []byte, rsaKeysPath string) ([]byte, error) {

	privateKey, err := parseRsaPrivateKeyFromPemFile(rsaKeysPath)
	if err != nil {
		return nil, err
	}

	decryptedCipherKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedCipherKey)
	if err != nil {
		return nil, err
	}

	return decryptedCipherKey, nil
}

func GenerateRSAKeyPair(rsaPath string) (string, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}

	privateFilePem, err := os.Create(rsaPath + "/rsa_key.pem")
	if err != nil {
		return "", err
	}

	privateKeyPem := exportRsaPrivateBlock(privateKey)
	if err = pem.Encode(privateFilePem, privateKeyPem); err != nil {
		return "", err
	}

	if err = privateFilePem.Close(); err != nil {
		return "", err
	}

	return filepath.Join(rsaPath, "rsa_key.pem"), nil
}

func exportRsaPrivateBlock(privateKey *rsa.PrivateKey) *pem.Block {

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	return privateBlock
}

func parseRsaPrivateKeyFromPemFile(rsaKeysPath string) (*rsa.PrivateKey, error) {

	// read from the pubkey file
	privateKeyFile, err := os.Open(rsaKeysPath)
	if err != nil {
		return nil, err
	}

	pemFileinfo, _ := privateKeyFile.Stat()

	var size int64 = pemFileinfo.Size()

	pemBytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)

	if _, err = buffer.Read(pemBytes); err != nil {
		return nil, err
	}

	privateKeyData, _ := pem.Decode([]byte(pemBytes))

	if err = privateKeyFile.Close(); err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyData.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func exportRsaPublicKeyAsPemString(publicKey *rsa.PublicKey) ([]byte, error) {

	publicKeyASN1, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyASN1,
		},
	)

	return publicKeyPem, nil
}

func parseRsaPublicKeyFromPemString(publicKeyPem string) (*rsa.PublicKey, error) {

	if block, _ := pem.Decode([]byte(publicKeyPem)); block == nil {
		return nil, errors.New("Failed to parse PEM block containing public key.")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return publicKey, nil

	default:
		break
	}

	return nil, errors.New("Key type is not RSA.")
}

func saveHashcode(cipherKey []byte, hashPath, documentName string) {

	// create hashcode
	hashCode := fnv.New128()
	hashCode.Write(cipherKey)

	if err := ioutil.WriteFile(hashPath, []byte(hashCode.Sum([]byte(cipherKey))), 0644); err != nil {
		fmt.Println(err)
	}
}

func getRsaKeyFromUpload(rsaPath, filename string) error {

	// check if path exists
	if _, err := os.Stat(rsaPath); os.IsNotExist(err) {
		if err := os.MkdirAll(rsaPath, 0755); err != nil {
			return err
		}
	}

	// copy file to path
	source, err := os.Open(filename)
	if err != nil {
		return err
	}

	destination, err := os.OpenFile(rsaPath+"/rsa_key.pem", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}

func DeleteRsaDirectory(path, version string) error {

	rsaPath := filepath.Join(path, version, "rsa")

	if _, err := os.Stat(path); os.IsNotExist(err); err != nil {
		return errors.New("Path to rsa key pair does not exist")
	}

	if _, err = os.Stat(rsaPath); os.IsNotExist(err) {
		fmt.Println(err) // ? 
		return nil
	}

	directory, err := ioutil.ReadDir(rsaPath)
	if err != nil {
		return err
	}

	for _, item := range directory {
		os.RemoveAll(filepath.Join([]string{rsaPath, item.Name()}...))
	}

	os.Remove(rsaPath)

	return nil
}
