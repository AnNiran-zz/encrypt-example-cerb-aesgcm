package ipfs

import (
	"fmt"
	"github.com/ipfs/go-ipfs-api"
)

func createNewEmptyDirectory() (string, error) {

	runShellInstance()

	var newEmptyDirectory string

	newEmptyDirectory, err := sh.NewObject("unixfs-dir")

	if err != nil {
		return newEmptyDirectory, err
	}

	return newEmptyDirectory, nil
}

// directoryName = linkName = "personAccounts"
func CreateGroupAccountsIpfsDirectory(directoryName string) (string, error) {

	runShellInstance()

	newEmptyDirectory, err := sh.NewObject("unixfs-dir")

	if err != nil {
		return "", err
	}

	currentDirectory := newEmptyDirectory

	object, err := makeRandomObject()

	if err != nil {
		return "", nil
	}

	newObject, err := sh.PatchLink(currentDirectory, directoryName, object, true)

	if err != nil {
		return "", err
	}

	currentDirectory = newObject

	return currentDirectory, nil
}

// creates additional object - function according to documentations
// parentDirectory = "personAccounts"
func CreateIpfsAccountDirectory(directoryName, parentDirectory string) (*IpfsDirectoryData, string, error) {

	runShellInstance()

	newEmptyDirectory, err := sh.NewObject("unixfs-dir")

	if err != nil {
		return nil, "", err
	}

	currentDirectory := newEmptyDirectory

	object, err := makeRandomObject()

	if err != nil {
		return nil, "", err
	}

	// arguments: root, path, childhash, create
	newObject, err := sh.PatchLink(parentDirectory, directoryName, object, true)

	currentDirectory = newObject

	// create link object
	linkObject, err := sh.Patch(parentDirectory, "add-link", "links", currentDirectory)

	if err != nil {
		return nil, "", err
	}

	newAccountDirectory := &IpfsDirectoryData{
		ContentIdentifier:   newEmptyDirectory,
		ObjectHash:          currentDirectory,
		Reference:           directoryName,
		LinkObjectHash:      linkObject,
		ParentDirectoryHash: parentDirectory,
	}

	return newAccountDirectory, "", nil
}

// creates additional object - function according to documentations
// parentDirectory = "email"
// directoryName = linkName = "documentName"
func CreateIpfsDocumentDirectory(directoryName, parentDirectoryHash, parentLinkObject string) (*IpfsDirectoryData, string, string, error) {

	runShellInstance()

	newEmptyDirectory, err := sh.NewObject("unixfs-dir")

	if err != nil {
		return nil, "", newEmptyDirectory, err
	}

	currentDirectory := newEmptyDirectory

	object, err := makeRandomObject()

	if err != nil {
		return nil, "", newEmptyDirectory, err
	}

	newObject, err := sh.PatchLink(parentDirectoryHash, directoryName, object, true)

	if err != nil {
		return nil, "", "", err
	}

	currentDirectory = newObject

	// create link object
	// arguments: root, action, args ...string
	linkObject, err := sh.Patch(parentDirectoryHash, "add-link", "links", currentDirectory)

	if err != nil {
		return nil, "", "", err
	}

	newDocumentDirectory := &IpfsDirectoryData{
		ContentIdentifier:   newEmptyDirectory,
		ObjectHash:          currentDirectory,
		Reference:           directoryName,
		LinkObjectHash:      linkObject,
		ParentDirectoryHash: parentDirectoryHash,
	}

	// update parent directory linkObject
	newParentLinkObject, err := sh.Patch(parentLinkObject, "add-link", directoryName, currentDirectory)

	if err != nil {
		fmt.Println(err)
		return nil, "", "", err
	}

	// delete previous linkObject
	sh.Request("rm", parentLinkObject)

	return newDocumentDirectory, newParentLinkObject, "", nil
}

func getIpfsAccountDirectory(accountDirectoryHash string) (*shell.IpfsObject, error) {

	runShellInstance()

	accountIpfsObject, err := sh.ObjectGet(accountDirectoryHash)

	if err != nil {
		return nil, err
	}

	return accountIpfsObject, nil
}

func getIpfsDocumentDirectory(directoryHash string) (*shell.IpfsObject, error) {

	runShellInstance()

	directoryIpfsObject, err := sh.ObjectGet(directoryHash)

	if err != nil {
		return nil, err
	}

	return directoryIpfsObject, nil
}

func DeleteDirectoryFromIpfs(directoryHash, linkObjectHash string) {

	runShellInstance()

	sh.Request("rm", directoryHash)
	sh.Request("rm", linkObjectHash)
}
