package ipfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var personAccountsPath = os.Getenv("GOPATH") + "/src/cerberus/ipfs/personAccounts"
var institutionAccountsTempPath = os.Getenv("GOPATH") + "/src/cerberus/ipfs/institutionAccounts"

func CreateAccountIpfsTempDirectory(path, directoryName string) (string, error) {

	accountPath := filepath.Join(path, directoryName)

	err := os.MkdirAll(accountPath, 0755)

	if err != nil {
		return "", err
	}

	return accountPath, nil
}

func GetAccountIpfsTempDirectory(path, directoryName string) (string, error) {

	accountPath := filepath.Join(path, directoryName)

	_, err := os.Stat(accountPath)

	if os.IsNotExist(err) {
		path, err := CreateAccountIpfsTempDirectory(path, directoryName)

		if err != nil {
			return "", err
		}

		accountPath = path
	}

	return accountPath, nil
}

func CreateDocumentIpfsTempDiretory(path, documentName string) (string, error) {

	documentPath := filepath.Join(path, documentName)

	err := os.MkdirAll(documentPath, 0755)

	if err != nil {
		return "", err
	}

	rsaPath := filepath.Join(documentPath, "rsa")
	err = os.MkdirAll(rsaPath, 0755)

	if err != nil {
		return "", err
	}

	return documentPath, nil
}

func GetDocumentIpfsTempDirectory(path, account, documentName string) (string, error) {

	accountPath, err := GetAccountIpfsTempDirectory(path, account)

	if err != nil {
		return "", err
	}

	documentPath := filepath.Join(accountPath, documentName)

	_, err = os.Stat(documentPath)

	if os.IsNotExist(err) {
		path, err := CreateDocumentIpfsTempDiretory(accountPath, documentName)

		if err != nil {
			return "", err
		}

		documentPath = path
	}

	return documentPath, nil
}

func deleteFile(fileName string) error {

	err := os.Remove(fileName)

	if err != nil {
		return err
	}

	return nil
}

func DeleteDocumentIpfsTempDirectory(path, account, documentName string) (string, error) {

	accountPath, err := GetAccountIpfsTempDirectory(path, account)

	if err != nil {
		return "", nil
	}

	documentPath := filepath.Join(accountPath, documentName)

	_, err = os.Stat(documentPath)

	if os.IsNotExist(err) {

		if err != nil {
			return "", nil
		}
	}

	directory, err := ioutil.ReadDir(documentPath)

	if err != nil {
		return "", err
	}

	for _, item := range directory {
		os.RemoveAll(filepath.Join([]string{documentPath, item.Name()}...))
	}

	os.Remove(documentPath)

	return "", nil
}
