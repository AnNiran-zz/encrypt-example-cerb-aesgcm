package ipfs

import (
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	ipfsUtils "github.com/ipfs/go-ipfs-util"
	"gopkg.in/mgo.v2/bson"
	"io"
	"math/rand"
	"time"
)

type IpfsDocumentVersionData struct {
	ContentIdentifier        string `json:"contentIdentifier"`
	ObjectHash               string `json:"objectHash"`
	Reference                string `json:"reference"`
	ParentDirectoryHash      string `json:"parentDirectoryHash"`
	ParentDirectoryReference string `json:"parentDirectoryReference"`
}

type documentVersion struct {
	Id        string                   `json:"id"`
	Name      int                      `json:"name"`
	IpfsData  *IpfsDocumentVersionData `json:"IpfsData"`
	CipherKey string                   `json:"cipherKey"`
	Hash      string                   `json:"hash"`
	CreatedAt string                   `json:"createdAt"`
	UpdateAt  string                   `json:"updatedAt"`
}

type IpfsDirectoryData struct {
	ContentIdentifier   string `json:"contentIdentifier"`
	ObjectHash          string `json:"objectHash"`
	Reference           string `json:"reference"`
	LinkObjectHash      string `json:"linkObjectHash"`
	ParentDirectoryHash string `json:"parentDirectoryHash"`
}

type documentDirectory struct {
	ObjectType                string                   `json:"docType"`
	DocumentName              string                   `json:"documentName"`
	PersonName                string                   `json:"personName"`
	CountryIssue              string                   `json:"countryIssue"`
	IpfsDocumentData          *IpfsDirectoryData       `json:"ipfsDocumentData"`
	IpfsDocumentVersionsData  map[int]*documentVersion `json:"ipfsDocumentVersionsData"`
	CreatedAt                 string                   `json:"createdAt"`
	UpdatedAt                 string                   `json:"updatedAt"`
}

type personAccount struct {
	Id                  bson.ObjectId                 `json:"id"`
	ObjectType          string                        `json:"docType"`
	FirstName           string                        `json:"firstName"`
	LastName            string                        `json:"lastName"`
	Email               string                        `json:"email"`
	Phone1              string                        `json:"phone1"`
	Phone2              string                        `json:"phone2"`
	IpfsAccountData     *IpfsDirectoryData            `json:"ipfsAccountData"`
	CreatedAt           string                        `json:"createdAt"`
	UpdatedAt           string                        `json:"updatedAt"`
	Documents           map[string]*documentDirectory `json:"documents"`
}

var sh     *shell.Shell
var ncalls int

var _ = time.ANSIC

func sleep() {
	ncalls++
}

func runShellInstance() {

	sh = shell.NewShell("localhost:5001")

	//for i := 0; i < 200; i++ {
	//	_, err := makeRandomObject()
	//	if err != nil {
	//		fmt.Println("err: ", err)
	//		return
	//	}
	//}
}

func makeRandomObject() (string, error) {

	x := rand.Intn(120) + 1
	y := rand.Intn(120) + 1
	z := rand.Intn(120) + 1
	size := x * y * z

	r := io.LimitReader(ipfsUtils.NewTimeSeededRand(), int64(size))
	sleep()

	return sh.Add(r)
}

func unpinFileFromIpfs(linkObject, documentName string) error {

	documentLocation, err := getFileFromIpfs(linkObject, documentName)

	if err != nil {
		return err
	}

	// here we use unpin mechanism -> garbage collector
	err = sh.Unpin(documentLocation)

	if err != nil {
		fmt.Println(err)
	}

	// TODO: unpin recursively for large files

	return nil
}

func createRandomDirectory(depth int, linkName string) (string, error) {

	if depth <= 0 {
		return makeRandomObject()
	}

	sleep()

	empty, err := sh.NewObject("unixfs-dir")

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	currentDir := empty

	for i := 0; i < rand.Intn(8) + 1; i++ {
		var object string
		if rand.Intn(2) == 1 {
			object, err = makeRandomObject()
			if err != nil {
				return "", err
			}
		} else {
			object, err = createRandomDirectory(depth - 1, linkName)
			if err != nil {
				return "", err
			}
		}

		name := linkName
		sleep()

		newObject, err := sh.PatchLink(currentDir, name, object, true)

		if err != nil {
			return "", err
		}

		currentDir = newObject
	}

	return currentDir, nil
}
