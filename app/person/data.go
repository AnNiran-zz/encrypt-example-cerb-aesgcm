package person

import (
	"cerberus/services/ipfs"
	"os"
)

type documentShareableData struct {
	documentId   string
	documentName string
	holder       string
	countryIssue string
	// some other data
}

type accountShareableData struct { // better name?
	firstName string
	lastName  string
	email     string
	phone     string
	// some other data
}

type accountDataRequest struct {
	Id                string            `json:"id"`
	PublicId          string            `json:"publicId"`
	RequestType       string            `json:"requestType"`
	ObjectType        string            `json:"docType"`
	RequesterPublicId string            `json:"requesterPublicId"`
	RecipientPublicId string            `json:"recipientPublicId"`
	RequestedData     string            `json:"requestedData"`
	AccountData       map[string]string `json:"accountData"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Status            string            `json:"status"`
}

type documentDataRequest struct {
	Id                string            `json:"id"`
	PublicId          string            `json:"publicId"`
	RequestType       string            `json:"requestType"`
	ObjectType        string            `json:"docType"`
	RequesterPublicId string            `json:"requesterPublicId"`
	RecipientPublicId string            `json:"recipientPublicId"`
	DocumentName      string            `json:"documentName"`
	DocumentData      map[string]string `json:"documentData"`
	DocumentCopy      bool              `json:"documentCopy"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Status            string            `json:"status"`
}

type documentVersion struct {
	Id        string                        `json:"id"`
	Name      int                           `json:"name"`
	IpfsData  *ipfs.IpfsDocumentVersionData `json:"ipfsData"`
	CreatedAt string                        `json:"createdAt"`
	UpdateAt  string                        `json:"updatedAt"`
}

type documentData struct {
	DocumentId   string `json:"documentId"`
	DocumentName string `json:"documentName"`
	Holder       string `json:"holder"`
	CountryIssue string `json:"countryIssue"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type documentDirectory struct {
	Id                        string                   `json:"id"`
	ObjectType                string                   `json:"docType"`
	DocumentData              *documentData            `json:"documentData"`
	IpfsDocumentDirectoryData *ipfs.IpfsDirectoryData  `json:"ipfsDocumentDirectoryData"`
	IpfsDocumentVersionsData  map[int]*documentVersion `json:"ipfsDocumentVersionsData"`
	UpdatedAt                 string                   `json:"updatedAt"`
}

type accountData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type personAccount struct {
	Id              string                        `json:"id"`
	PublicId        string                        `json:"publicId"`
	ObjectType      string                        `json:"docType"`
	AccountData     *accountData                  `json:"accountData"`
	IpfsAccountData *ipfs.IpfsDirectoryData       `json:"ipfsAccountData"`
	Documents       map[string]*documentDirectory `json:"documents"`
}

var personAccountsIpfsHash = "QmNMEveW869ERwNWCk4YSmuU5bR3j4AonmYoXaruTD51rA"
var personAccountsIpfsTempPath = os.Getenv("GOPATH") + "/src/cerberus/ipfs/personAccounts"
