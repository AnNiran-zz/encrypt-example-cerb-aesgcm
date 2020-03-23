package ipfs

import (
	"bytes"
	"cerberus/services/crypto"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// creates additional object - function according to documentations
// parentDirectory = "document1"
// fileName = linkName = x
func UploadFileToIpfs(documentData []byte, documentName, parentDirectoryReference, parentDirectoryHash, parentLinkObject string) (*IpfsDocumentVersionData, string, error) {

	var updatedParentLinkObject string

	// upload data to ipfs
	runShellInstance()

	reader := bytes.NewReader(documentData)

	contentIdentifier, err := sh.Add(reader)

	if err != nil {
		return nil, updatedParentLinkObject, err
	}

	documentObjectHash, err := sh.PatchLink(parentDirectoryHash, documentName, contentIdentifier, true)

	if err != nil {
		return nil, updatedParentLinkObject, err
	}

	newDocumentVersion := &IpfsDocumentVersionData{
		ContentIdentifier:        contentIdentifier,
		ObjectHash:               documentObjectHash,
		Reference:                documentName,
		ParentDirectoryHash:      parentDirectoryHash,
		ParentDirectoryReference: parentDirectoryReference,
	}

	// update parent directory linkObject
	updatedParentLinkObject, err = sh.Patch(parentLinkObject, "add-link", documentName, documentObjectHash)

	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	// delete previous linkObject
	sh.Request("rm", parentLinkObject)

	return newDocumentVersion, updatedParentLinkObject, nil
}

// documentHash - CID
// parentHas
// parentObjectLink
// document reference
func UpdateIpfsFileReference(documentVersion *IpfsDocumentVersionData, parentDirectoryHash, parentLinkObject string) (*IpfsDocumentVersionData, string, error) {

	var updatedParentLinkObject string
	var documentObjectHash string

	runShellInstance()

	var err error

	documentObjectHash, err = sh.PatchLink(parentDirectoryHash, documentVersion.ParentDirectoryReference, documentVersion.ContentIdentifier, true)

	if err != nil {
		return nil, "", err
	}

	updatedDocumentVersion := &IpfsDocumentVersionData{
		ContentIdentifier:        documentVersion.ContentIdentifier,
		ObjectHash:               documentObjectHash,
		Reference:                documentVersion.Reference,
		ParentDirectoryHash:      parentDirectoryHash,
		ParentDirectoryReference: documentVersion.ParentDirectoryReference,
	}

	updatedParentLinkObject, err = sh.Patch(parentLinkObject, "add-link", documentVersion.Reference, documentObjectHash)

	if err != nil {
		return nil, "", err
	}

	return updatedDocumentVersion, updatedParentLinkObject, nil
}

func DeleteDocumentObjectFromIpfs(document *IpfsDocumentVersionData) {

	runShellInstance()

	sh.Request("rm", document.ContentIdentifier)
	sh.Request("rm", document.ObjectHash)
}

func ExportFileFromIpfs(objectHash, documentVersionName, destinationPath string, cipherKey []byte) (string, error) {

	runShellInstance()

	fileData, err := readFileFromIpfs(objectHash, documentVersionName, destinationPath)

	if err != nil {
		return "", err
	}

	// read file data form file
	fileDataEncrypted, err := ioutil.ReadFile(fileData)

	// decrypt process
	rsaPath := filepath.Join(destinationPath, documentVersionName, "rsa")
	rsa := rsaPath + "/rsa_key.pem"
	fileDataAsBytes, err := crypto.DecryptDocument(fileDataEncrypted, cipherKey, rsa)

	if err != nil {
		return "", err
	}

	// save file data as png
	filePath, err := convertToPng(fileDataAsBytes, destinationPath+"/"+documentVersionName+".png")

	if err != nil {
		return "", err
	}

	if err = os.Remove(fileData); err != nil {
		return "", err
	}

	return filePath, nil
}

func getFileFromIpfs(objectHash, documentVersionName string) (string, error) {

	runShellInstance()

	var documentLocation string

	object, err := sh.ObjectGet("/ipfs/" + objectHash)

	if err != nil {
		return documentLocation, err
	}

	documentLinks := object.Links

	for _, value := range documentLinks {

		if value.Name == documentVersionName {
			documentLocation = value.Hash
		}
	}

	if documentLocation == "" {
		return documentLocation, errors.New("Provided document name: " + documentVersionName + " does not match any links for object: " + objectHash)
	}

	return documentLocation, nil
}

func readFileFromIpfs(cid, documentVersionName, destinationDirectory string) (string, error) {

	documentLocation, err := getFileFromIpfs(cid, documentVersionName)

	// obtain document data
	reader, err := sh.Cat(documentLocation)

	defer reader.Close()

	if err != nil {
		return "", err
	}

	fileLocation := filepath.Join(destinationDirectory, documentVersionName)

	outputFile, err := os.Create(fileLocation)

	_, err = io.Copy(outputFile, reader)

	return fileLocation, nil
}

func convertToPng(fileDataAdBytes []byte, filePath string) (string, error) {

	img, _, _ := image.Decode(bytes.NewReader(fileDataAdBytes))

	out, err := os.Create(filePath) // create png extension file

	if err != nil {
		return "", err
	}

	err = png.Encode(out, img)

	if err != nil {
		return "", err
	}

	return filePath, nil
}
