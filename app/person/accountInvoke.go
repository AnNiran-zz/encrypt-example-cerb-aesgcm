package person

import (
	"cerberus/blockchain/persaccntschannel"
	"cerberus/services/crypto"
	"cerberus/services/ipfs"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func CreateAccount(firstName, lastName, email, phone string) ([]string, []string, error) {

	if firstName == "" {
		return nil, "", errors.New("First name value cannot be an empty string")
	}

	if lastName == "" {
		return nil, "", errors.New("Last name value cannot be an empty string")
	}

	if email == "" {
		return nil, "", errors.New("Email value cannot be an empty string")
	}

	if phone == "" {
		return nil, "", errors.New("Phone value cannot be an empty string")
	}

	// create object
	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)
	email = strings.ToLower(email)

	id := bson.NewObjectId().Hex()
	publicID := crypto.HashMD5(id)

	documents := make(map[string]*documentDirectory)

	accountData := &accountData{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		CreatedAt: getTime(),
	}

	accountObject := &personAccount{
		ID:          id,
		PublicID:    publicID,
		ObjectType:  "person",
		AccountData: accountData,
		Documents:   documents,
	}

	// create personAccount folder in ipfs
	// linkReference := "/personAccounts/" + email
	ipfsData, _, err := ipfs.CreateIpfsAccountDirectory(publicID, personAccountsIpfsHash)
	if err != nil {
		return nil, "", err
	}

	accountObject.IpfsAccountData = ipfsData
	accountObjectAsBytes, err := json.Marshal(accountObject)
	if err != nil {
		return nil, "", err
	}

	// encrypt account object AES-GCM using the provided passphrase
	// key is returned rom this function and must be provided for data decryption
	key := crypto.Key32byt()
	aesGCM, err := crypto.EncAESCGM(accountObjectAsBytes, key)
	if err != nil {
		return nil, "", err
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.CreateAccount(publicID, accountObjectAsBytes)

	if err != nil {
		ipfs.DeleteDirectoryFromIpfs(ipfsData.ObjectHash, ipfsData.LinkObjectHash)
		return nil, "", err
	}

	return response, []string{string(accountObjectAsBytes), string(key)}, nil
}

// Update account:
/*
Selectors:
- FirstName
- LastName
- Phone
- etc
*/
func UpdateAccountBySelector(accountPublicID, key, selectorName, selectorValue string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account Public ID value cannot be an empty string")
	}

	if selectorName == "" {
		return nil, nil, errors.New("SelectorName value cannot be an empty string")
	}

	if selectorValue == "" {
		return nil, nil, errors.New("SelectorValue cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Passphrase cannot be an empty string")
	}

	selector := strings.ToLower(selectorName)

	switch selector {
	case "firstname":
		selectorName = "FirstName"

	case "lastname":
		selectorName = "LastName"

	case "email":
		selectorName = "Email"

	case "phone":
		selectorName = "Phone"
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.UpdateRecords("updateAccount", []string{accountPublicId, key, selectorName, strings.ToLower(selectorValue)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(newAccountData)}, response, nil
}

func UpdateAccountFirstName(accountPublicID, key, firstName string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account Public Id cannot be an empty string")
	}

	if firstName == "" {
		return nil, nil, errors.New("First name value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	dataField := "FirstName"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.UpdateRecords("updateAccount", []string{accountPublicID, key, dataField, strings.ToLower(firstName)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(newAccountData)}, response, nil
}

func UpdateAccountLastName(accountPublicID, key, lastName string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account Public Id cannot be an empty string")
	}

	if lastName == "" {
		return nil, nil, errors.New("Last name value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	dataField := "LastName"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.UpdateRecords("updateAccount", []string{accountPublicID, key, dataField, strings.ToLower(lastName)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(newAccountData)}, response, nil
}

func UpdateAccountPhone(accountPublicID, key, phone string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account Public Id cannot be an empty string")
	}

	if phone == "" {
		return nil, nil, errors.New("Phome name value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	dataField := "Phone"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.UpdateRecords("updateAccount", []string{accountPublicID, key, dataField, phone})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(newAccountData)}, response, nil
}

func UpdateAccountEmail(accountPublicID, key, email string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account Public ID cannot be an empty string")
	}

	if email == "" {
		return nil, nil, errors.New("Phome name value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	dataField := "Email"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, newAccountData, err := persAccntsChannelClient.UpdateRecords("updateAccount", []string{accountPublicID, key, dataField, email})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(newAccountData)}, response, nil
}

func DeleteAccount(accountPublicID string) ([]string, error) {

	if accountPublicID == "" {
		return nil, errors.New("Account Public ID cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, deletedRecord, err := persAccntsChannelClient.DeleteAccount(accountPublicID)

	if err != nil {
		return nil, err
	}

	record := &personAccount{}
	if err = json.Unmarshal([]byte(deletedRecord), record); err != nil {}
		return nil, err
	}

	// delete records from ipfs
	ipfs.DeleteDirectoryFromIpfs(record.IpfsAccountData.ObjectHash, record.IpfsAccountData.LinkObjectHash)

	for _, directory := range record.Documents {
		_, response, err = DeleteDocument(accountPublicID, directory.DocumentData.DocumentName)

		if err != nil {
			return nil, err
		}
	}

	// ...
	return response, nil
}

func CreateNewDocument(accountPublicID, key, documentName, holderName, countryIssue, filename string) ([]string, []string, string, error) {

	if accountPublicID == "" {
		return nil, nil, "", errors.New("Id value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, "", errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, "", errors.New("Document name value cannot be an empty string")
	}

	if holderName == "" {
		return nil, nil, "", errors.New("Holder name value cannot be an empty string")
	}

	if countryIssue == "" {
		return nil, nil, "", errors.New("Country issue value cannot be an empty string")
	}

	if filename == "" {
		return nil, nil, "", errors.New("Filename value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)
	holderName = strings.ToLower(holderName)
	countryIssue = strings.ToLower(countryIssue)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountPublicID)
	if err != nil {
		return nil, nil, "", err
	}

	// Decrypt account data from the Database using the account key
	decrRecord, err := DecrAESGCM(accountRecords, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	recordUpdate := &personAccount{}
	if err = json.Unmarshal(decrRecord, recordUpdate); err != nil {
		return nil, nil, "", err
	}

	// check if document folder already exists in account record
	if _, ok := recordUpdate.Documents[documentName]; ok {
		return nil, nil, "", errors.New("Document with name " + documentName + " already exists. ")
	}

	// create ipfs temp directory
	newDocumentIpfsDirectory, newDocumentVersion, updatedAccountIpfsLinks, rsaLink, err := createFirstDocumentVersion(filename, documentName, recordUpdate.PublicId, recordUpdate.IpfsAccountData.ObjectHash, recordUpdate.IpfsAccountData.LinkObjectHash)
	if err != nil {
		return nil, nil, "", err
	}

	// add new document directory to record
	newDocument := &documentDirectory{
		ID:         bson.NewObjectId().Hex(),
		ObjectType: "docType",
		DocumentData: &documentData{
			DocumentName: documentName,
			Holder:       holderName,
			CountryIssue: countryIssue,
			CreatedAt:    getTime(),
		},
		IpfsDocumentDirectoryData: newDocumentIpfsDirectory,
		IpfsDocumentVersionsData:  make(map[int]*documentVersion),
	}

	// add new document version to the folder
	newDocument.IpfsDocumentVersionsData[newDocumentVersion.Name] = newDocumentVersion

	recordUpdate.Documents[documentName] = newDocument
	recordUpdate.IpfsAccountData.LinkObjectHash = updatedAccountIpfsLinks
	recordUpdate.AccountData.CreatedAt = getTime()

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, "", err
	}

	// Encrypt the new record
	encrRecord, err := crypto.EncAESCGM(recordUpdateAsBytes, key)
	if err != nil {
		return nil, nil, "", err
	}

	response, updatedAccount, err := persAccntsChannelClient.UpdateRecords("updateDocumentRecords", []string{[]string{accountPublicId, "", string(encrRecord)})
	if err != nil {
		return nil, nil, "", err
	}

	return []string{string(updatedAccount)}, response, rsaLink, nil
}

func CreateDocumentVersion(accountPublicID, key, documentName, filename, rsaLink string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("ID value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, errors.New("Document name value cannot be an empty string")
	}

	if filename == "" {
		return nil, nil, errors.New("Filename value cannot be an empty string")
	}

	if rsaLink == "" {
		return nil, nil, errors.New("Rsa link value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountPublicID)
	if err != nil {
		return nil, nil, err
	}

	// Decrypt data from the Database
	decrRecord, err := DecrAESGCM(accountRecords, key)
	if err != nil {
		return nil, nil, "", err
	}

	recordUpdate := &personAccount{}
	if err = json.Unmarshal([]byte(decrRecord), recordUpdate); err != nil {
		return nil, nil, err
	}

	// check if document folder already exists in account record
	if _, ok := recordUpdate.Documents[documentName]; !ok {
		return nil, nil, errors.New("Document with name " + documentName + " does not exist")
	}

	// get existing document directory
	document := recordUpdate.Documents[documentName] // type documentDirectory

	// create document next version name
	newVersionNumber := getNextDocumentVersion(document.IpfsDocumentVersionsData)
	newDocumentVersion, updatedDirectoryLinks, err := createNewDocumentVersion(newVersionNumber, filename, rsaLink, recordUpdate.PublicID, documentName, document.IpfsDocumentDirectoryData.ObjectHash, document.IpfsDocumentDirectoryData.LinkObjectHash)
	if err != nil {
		return nil, nil, err
	}

	document.IpfsDocumentVersionsData[newDocumentVersion.Name] = newDocumentVersion
	document.IpfsDocumentDirectoryData.LinkObjectHash = updatedDirectoryLinks
	document.UpdatedAt = getTime()

	// add new version to the record
	recordUpdate.Documents[documentName].IpfsDocumentVersionsData[newDocumentVersion.Name] = newDocumentVersion

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the new record
	encrRecord, err := crypto.EncAESCGM(recordUpdateAsBytes, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	response, updatedAccount, err := persAccntsChannelClient.UpdateRecords("updateDocumentRecords", []string{[]string{accountPublicId, "", string(encrRecord)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(updatedAccount)}, response, nil
}

func UpdateDocumentCountryIssue(accountPublicID, key, documentName, countryIssueUpdate string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("ID value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, errors.New("Document name value cannot be an empty string")
	}

	if countryIssueUpdate == "" {
		return nil, nil, errors.New("Country issue update value cannot be an empty string")
	}

	persAccntsChanelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChanelClient.QueryAccountData("getAccountRecords", accountPublicID)
	if err != nil {
		return nil, nil, err
	}

	// Decrypt account data from the Database
	decrRecord, err := crypto.DecrAESGCM(accountRecords, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}
	
	recordUpdate := &personAccount{}
	if err = json.Unmarshal(decrRecord, recordUpdate); err != nil {
		return nil, nil, "", err
	}

	if _, ok := recordUpdate.Documents[documentName]; !ok {
		return nil, nil, errors.New("Document with name " + documentName + " does not exist")
	}

	recordUpdate.Documents[documentName].DocumentData.CountryIssue = countryIssueUpdate
	recordUpdate.Documents[documentName].UpdatedAt = getTime()

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, err
	}

	encrRecord, err := crypto.EncAESCGM(accountObjectAsBytes, []byte(key))
	if err != nil {
		return nil, "", err
	}

	response, updatedRecord, err := persAccntsChanelClient.UpdateRecords("updateDocumentRecords", []string{[]string{accountPublicId, "", string(encrRecord)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(updatedRecord)}, response, nil
}

func UpdateDocumentHolderName(accountPublicID, key, documentName, personNameUpdate string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("ID value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, errors.New("Document name value cannot be an empty string")
	}

	if personNameUpdate == "" {
		return nil, nil, errors.New("Person name update value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)
	personNameUpdate = strings.ToLower(personNameUpdate)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountPublicID)
	if err != nil {
		return nil, nil, err
	}

	// Decrypt account data from the Database
	decrRecord, err := crypto.DecrAESGCM(accountRecords, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	recordUpdate := &personAccount{}
	if err = json.Unmarshal(decrRecord, recordUpdate); err != nil {
		return nil, nil, err
	}

	if _, ok := recordUpdate.Documents[documentName]; !ok {
		return nil, nil, errors.New("Document with name " + documentName + " does not exist")
	}

	recordUpdate.Documents[documentName].DocumentData.Holder = personNameUpdate
	recordUpdate.Documents[documentName].UpdatedAt = getTime()

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, err
	}

	encrRecord, err := crypto.EncAESCGM(recordUpdateABytes, []byte(key))
	if err != nil {
		return nil, "", err
	}

	response, updatedRecord, err := persAccntsChannelClient.UpdateRecords("updateDocumentRecords", []string{[]string{accountPublicId, "", string(encrRecord)})
	if err != nil {
		return nil, nil, err
	}

	return []string{string(updatedRecord)}, response, nil
}

func DeleteDocument(accountPublicID, key, documentName string) ([]string, []string, error) {

	if accountPublicID == "" {
		return nil, nil, errors.New("Account ID value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, errors.New("Document name value cannot be an empty string")
	}

	documentName = strings.ToLower(documentName)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountPublicID)
	if err != nil {
		return nil, nil, err
	}

	// Decrypt account data from the Database
	decrRecord, err := crypto.DecrAESGCM(accountRecords, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	recordUpdate := &personAccount{}
	if err = json.Unmarshal(decrRecord, recordUpdate); err != nil {
		return nil, nil, err
	}

	// check if document exists
	documents := recordUpdate.Documents

	if _, ok := documents[documentName]; !ok {
		return nil, nil, errors.New("Document with name: " + documentName + " does not exist.")
	}

	documentToDelete := documents[documentName]

	delete(recordUpdate.Documents, documentName)

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, err
	}

	encrRecord, err := crypto.EncAESCGM(recordUpdateABytes, []byte(key))
	if err != nil {
		return nil, "", err
	}

	response, updatedRecord, err := persAccntsChannelClient.UpdateRecords("updateDocumentRecords", []string{accountPublicID, "", string(encrRecord)})
	if err != nil {
		return nil, nil, err
	}

	// delete records from ipfs
	ipfs.DeleteDirectoryFromIpfs(documentToDelete.IpfsDocumentDirectoryData.ObjectHash, documentToDelete.IpfsDocumentDirectoryData.LinkObjectHash)

	for _, version := range documentToDelete.IpfsDocumentVersionsData {
		ipfs.DeleteDocumentObjectFromIpfs(version.IpfsData)

		ipfsTempDocumentPath := filepath.Join(personAccountsIpfsTempPath, accountPublicID, documentName)

		if err = crypto.DeleteCipherKeyFile(ipfsTempDocumentPath, strconv.Itoa(version.Name)); err != nil {
			return nil, nil, err
		}

		if err = crypto.DeleteRsaDirectory(ipfsTempDocumentPath); err != nil {
			return nil, nil, err
		}
	}

	// delete folder from ipfs temp directory
	documentIpfsTempPath := filepath.Join(personAccountsIpfsTempPath)
	if _, err = ipfs.DeleteDocumentIpfsTempDirectory(documentIpfsTempPath, accountPublicId, documentName); err != nil {
		return nil, nil, err
	}

	return []string{string(updatedRecord)}, response, nil
}

func DeleteDocumentVersion(accountPublicId, key, documentName string, documentVersion int) ([]string, []string, error) {

	if accountPublicId == "" {
		return nil, nil, errors.New("Account Id value cannot be an empty string")
	}

	if key == "" {
		return nil, nil, errors.New("Key value cannot be an empty string")
	}

	if documentName == "" {
		return nil, nil, errors.New("Document name value cannot be an empty string")
	}

	if documentVersion < 1 {
		return nil, nil, errors.New("Document version value must be a valid version number")
	}

	documentName = strings.ToLower(documentName)

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	accountRecords, err := persAccntsChannelClient.QueryAccountData("getAccountRecords", accountPublicId)
	if err != nil {
		return nil, nil, err
	}

	// Decrypt account data from the Database
	decrRecord, err := crypto.DecrAESGCM(accountRecords, []byte(key))
	if err != nil {
		return nil, nil, "", err
	}

	recordUpdate := &personAccount{}
	if err = json.Unmarshal(decrRecord, recordUpdate); err != nil {
		return nil, nil, err
	}

	// check if document exists
	documents := recordUpdate.Documents

	if _, ok := documents[documentName]; !ok {
		return nil, nil, errors.New("Document with name: " + documentName + " does not exist.")
	}

	if _, ok := documents[documentName].IpfsDocumentVersionsData[documentVersion]; !ok {
		return nil, nil, errors.New("Document version " + strconv.Itoa(documentVersion) + " does not exist for " + documentName)
	}

	documentVersionToDelete := documents[documentName].IpfsDocumentVersionsData[documentVersion]

	delete(recordUpdate.Documents[documentName].IpfsDocumentVersionsData, documentVersion)
	recordUpdate.Documents[documentName].UpdatedAt = getTime()

	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return nil, nil, err
	}

	encrRecord, err := crypto.EncAESCGM(recordUpdateABytes, []byte(key))
	if err != nil {
		return nil, "", err
	}

	response, updatedRecord, err := persAccntsChannelClient.UpdateRecords("updateDocumentRecords", []string{accountPublicId, "", string(encrRecord)})
	if err != nil {
		return nil, nil, err
	}

	// delete records from ipfs
	ipfs.DeleteDocumentObjectFromIpfs(documentVersionToDelete.IpfsData)

	// delete cipher key
	ipfsTempDocumentPath := filepath.Join(personAccountsIpfsTempPath, accountPublicId, documentName)
	if err = crypto.DeleteCipherKeyFile(ipfsTempDocumentPath, strconv.Itoa(documentVersion)); err != nil {
		return nil, nil, err
	}

	return []string{string(updatedRecord)}, response, nil
}

//
func createNewDocumentVersion(newVersion int, filename, accountPublicId, documentName, parentDirHash, parentDirObjectLinkHash string) (*documentVersion, string, error) {

	nextDocumentVersionString := strconv.Itoa(newVersion)

	// get Ipfs temporary document directory
	ipfsTempDocumentPath, err := ipfs.GetDocumentIpfsTempDirectory(personAccountsIpfsTempPath, accountPublicId, documentName)
	if err != nil {
		return nil, "", err
	}

	// RSA key is saved inside /ipfs/accountID/documentName/version/rsa
	// temporary solution for the current implementation - keys are supposed to be
	// sent to the client
	rsaPath := filepath.Join(ipfsTempDocumentPath, nextDocumentVersionString, "rsa")
	rsaLink, err := crypto.GenerateRSAKeyPair(rsaPath)
	if err != nil {
		return nil, nil, "", "", err
	}

	// filename is the location of the scanned image before encryption
	// create 32-bit key for the document data encryption and encrypt the document data with it 
	// return signed key with the public key
	key := crypto.Key32byt()
	encryptedDocument, cipherKey, err := crypto.EncryptDocument(filename, key, rsaPath, rsaLink)
	if err != nil {
		return nil, nil, "", "", err
	}

	// save encrypted cipher key to a file - temporary solution
	// /ipfs/accountID/documentName/version/cipher
	if err = crypto.SaveCipherKey(cipherKey, ipfsTempDocumentPath, newDocumentVersionString); err != nil {
		return nil, nil, "", "", err
	}

	documentReference := filepath.Join(documentName, nextDocumentVersionString)

	documentVersionIpfsData, updatedDirectoryLinks, err := ipfs.UploadFileToIpfs(encryptedDocument, nextDocumentVersionString, documentReference, parentDirHash, parentDirObjectLinkHash)
	if err != nil {
		return nil, "", err
	}

	// create new document version
	newDocumentVersion := &documentVersion{
		Id:        bson.NewObjectId().Hex(),
		Name:      newVersion,
		IpfsData:  documentVersionIpfsData,
		CreatedAt: getTime(),
	}

	return newDocumentVersion, updatedDirectoryLinks, nil
}

func createFirstDocumentVersion(filename, documentName, accountPublicId, accountHash, accountObjectLinkHash string) (*ipfs.IpfsDirectoryData, *documentVersion, string, string, error) {

	// create new document directory in Ipfs network
	directoryName := documentName
	documentDirIpfsData, updatedAccountIpfsLinks, _, err := ipfs.CreateIpfsDocumentDirectory(directoryName, accountHash, accountObjectLinkHash)
	if err != nil {
		return nil, nil, "", "", err
	}

	// create document first version
	newDocumentVersionString := strconv.Itoa(1)

	// create Ipfs temporary document directory for saving rsa data on the server
	ipfsTempDocumentPath, err := ipfs.GetDocumentIpfsTempDirectory(personAccountsIpfsTempPath, accountPublicId, directoryName)
	if err != nil {
		return nil, nil, "", "", err
	}

	// RSA key is saved inside /ipfs/accountID/documentName/rsa
	// temporary solution for the current implementation - keys are supposed to be
	// sent to the client
	rsaPath := filepath.Join(ipfsTempDocumentPath, newDocumentVersionString, "rsa")
	rsaLink, err := crypto.GenerateRSAKeyPair(rsaPath)
	if err != nil {
		return nil, nil, "", "", err
	}

	// filename is the location of the scanned image before encryption
	// documentName is used as a passphrase for the document encryption
	// rsaPath - filename of rsa keys
	key := crypto.Key32byt(documentName)
	encryptedDocument, cipherKey, err := crypto.EncryptDocument(filename, key, rsaPath, rsaLink)
	if err != nil {
		return nil, nil, "", "", err
	}

	// save cipherKey
	// save encrypted cipher key to a file - temporary solution
	if err = crypto.SaveCipherKey(cipherKey, ipfsTempDocumentPath, newDocumentVersionString); err != nil {
		return nil, nil, "", "", err
	}

	documentReference := filepath.Join(documentName, newDocumentVersionString)
	documentVersionIpfsData, updatedDirectoryIpfsLinks, err := ipfs.UploadFileToIpfs(encryptedDocument, newDocumentVersionString, documentReference, documentDirIpfsData.ObjectHash, documentDirIpfsData.LinkObjectHash)

	if err != nil {
		return nil, nil, "", "", err
	}

	// create document version
	documentVersion := &documentVersion{
		Id:        bson.NewObjectId().Hex(),
		Name:      1,
		IpfsData:  documentVersionIpfsData,
		CreatedAt: getTime(),
	}

	documentDirIpfsData.LinkObjectHash = updatedDirectoryIpfsLinks

	return documentDirIpfsData, documentVersion, updatedAccountIpfsLinks, rsaLink, nil
}

func getNextDocumentVersion(documentVersions map[int]*documentVersion) int {

	keys := make([]int, 0, len(documentVersions))
	for k := range documentVersions {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	highestVersion := keys[len(keys)-1]
	nextVersion := highestVersion + 1

	return nextVersion
}

func getTime() string {

	currentDateTime := time.Now()
	CurrentDateTime := currentDateTime.Format("2006-01-02 15:04:05")

	return CurrentDateTime
}
