package person

import (
	"cerberus/blockchain/persaccntschannel"
	"cerberus/services/crypto"
	"cerberus/services/ipfs"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func GetAccountById(accountId, key string) (string, error) {

	if accountId == "" {
		return "", errors.New("Account Id value cannot be an empty string")
	}

	if key == "" {
		return "", errors.New(" Key value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountId)
	if err != nil {
		return "", err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountData, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	record := &personAccount{}
	if err = json.Unmarshal(decrRecord, record); err != nil {
		return nil, nil, "", err
	}

	return string(record), nil
}

// only for administration use
func GetAccountsByEmail(email string) (string, error) {

	if email == "" {
		return "", errors.New("Email value cannot be an empty string")
	}

	selectorKey := "email"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryRecords(selectorKey, email)
	if err != nil {
		return "", err
	}

	return string(accountData), nil
}

// only for administration use
func GetAccountsByFirstName(firstName string) (string, error) {

	if firstName == "" {
		return "", errors.New("First name value cannot be an empty string")
	}

	selectorKey := "firstName"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryRecords(selectorKey, firstName)
	if err != nil {
		return "", err
	}

	return string(accountData), nil
}

// only for administration use
func GetAccountsByLastName(lastName string) (string, error) {

	if lastName == "" {
		return "", errors.New("Last name value cannot be an empty string")
	}

	selectorKey := "lastName"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryRecords(selectorKey, lastName)
	if err != nil {
		return "", err
	}

	return string(accountData), nil
}

func GetAccountHistory(accountId, key string) (string, error) {

	if accountId == "" {
		return "", errors.New("Account Id value cannot be an empty string")
	}

	if key == "" {
		return "", errors.New(" Key value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryAccountData("getAccountHistory", accountId)
	if err != nil {
		return "", err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountData, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	record := &personAccount{}
	if err = json.Unmarshal(decrRecord, record); err != nil {
		return nil, nil, "", err
	}

	return string(record), nil
}

// selectors
// only for administration use

/*
- email
- firstName
- lastName
*/
func GetAccountsBySelector(selectorKey, selectorValue string) (string, error) {

	if selectorKey == "" {
		return "", errors.New("Selector key value cannot be an empty string")
	}

	if selectorValue == "" {
		return "", errors.New("Selector value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryRecords(selectorKey, selectorValue)
	if err != nil {
		return "", err
	}

	return string(accountData), nil
}

func GetAccountDocument(accountId, key documentName string) (string, error) {

	if accountId == "" {
		return "", errors.New("Account Id value cannot be an empty string")
	}

	if documentName == "" {
		return "", errors.New("Document name value cannot be an empty string")
	}

	if key == "" {
		return "", errors.New(" Key value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountId)
	if err != nil {
		return "", err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountData, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	record := &personAccount{}
	if err = json.Unmarshal(decrRecord, record); err != nil {
		return nil, nil, "", err
	}

	if _, ok := record.Documents[documentName]; ok {
		fmt.Println(record.Documents[documentName])
	} else {
		return "", errors.New("Document with name " + documentName + " does not exist")
	}

	documentDataAsBytes, err := json.Marshal(record.Documents[documentName])
	if err != nil {
		return "", err
	}

	return string(documentDataAsBytes), nil
}

func GetAccountDocumentVersion(accountId, key, documentName, documentVersion string) ([]string, error) {

	if accountId == "" {
		return nil, errors.New("Account Id value cannot be an empty string")
	}

	if documentName == "" {
		return nil, errors.New("Document name value cannot be an empty value")
	}

	if documentVersion == "" {
		return nil, errors.New("Document version value cannot be an empty string")
	}

	if key == "" {
		return "", errors.New(" Key value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountId)
	if err != nil {
		return nil, err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountData, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	record := &personAccount{}
	if err = json.Unmarshal(decrRecord, record); err != nil {
		return nil, nil, "", err
	}

	if _, ok := record.Documents[documentName]; !ok {
		return nil, errors.New("Document with name " + documentName + " does not exist")
	}

	ver, err := strconv.Atoi(documentVersion)
	if err != nil {
		return nil, err
	}

	document := record.Documents[documentName]
	version := document.IpfsDocumentVersionsData[ver]

	if _, ok := record.Documents[documentName].IpfsDocumentVersionsData[ver]; !ok {
		return nil, errors.New("Document version " + documentVersion + " for document " + documentName + " does not exist")
	}

	// get content from ipfs
	ipfsTempDocumentPath, err := ipfs.GetDocumentIpfsTempDirectory(personAccountsIpfsTempPath, record.PublicId, documentName)
	if err != nil {
		return nil, err
	}

	fmt.Println(ipfsTempDocumentPath)

	// get cipher key
	cipherKeyPath := filepath.Join(ipfsTempDocumentPath, "cipher", strconv.Itoa(version.Name))

	cipherKey, err := crypto.ReadCipherKey(cipherKeyPath)
	if err != nil {
		return nil, err
	}

	filename, err := ipfs.ExportFileFromIpfs(version.IpfsData.ObjectHash, strconv.Itoa(version.Name), ipfsTempDocumentPath, cipherKey)
	if err != nil {
		return nil, err
	}

	versionAsBytes, err := json.Marshal(version)
	if err != nil {
		return nil, err
	}

	return []string{string(versionAsBytes), filename}, nil
}

func GetAccountDocumentVersions(accountId, key documentName string) ([]string, error) {

	if accountId == "" {
		return nil, errors.New("Account Id value cannot be an empty string")
	}

	if documentName == "" {
		return nil, errors.New("Document name value cannot be an empty string")
	}

	if key == "" {
		return "", errors.New(" Key value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountData, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountId)
	if err != nil {
		return nil, err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountData, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	record := &personAccount{}
	if err = json.Unmarshal(decrRecord, record); err != nil {
		return nil, nil, "", err
	}

	if _, ok := record.Documents[documentName]; !ok {
		return nil, errors.New("Document with name " + documentName + " does not exist")
	}

	// get content from ipfs
	ipfsTempDocumentPath, err := ipfs.GetDocumentIpfsTempDirectory(personAccountsIpfsTempPath, record.PublicId, documentName)
	if err != nil {
		return nil, err
	}

	var versions []string
	for versionNumber, version := range record.Documents[documentName].IpfsDocumentVersionsData {

		// get cipher key
		cipherKeyPath := filepath.Join(ipfsTempDocumentPath, strconv.Itoa(versionNumber), "cipher")

		cipherKey, err := crypto.ReadCipherKey(cipherKeyPath)
		if err != nil {
			return nil, err
		}

		_, err = ipfs.ExportFileFromIpfs(version.IpfsData.ObjectHash, strconv.Itoa(version.Name), ipfsTempDocumentPath, cipherKey)
		if err != nil {
			return nil, err
		}

		versionAsBytes, err := json.Marshal(version)
		if err != nil {
			return nil, err
		}

		versions = append(versions, string(versionAsBytes))
	}

	return versions, nil
}
