package main

type accountDataRequest struct {
	ID                string            `json:"id"`
	PublicID          string            `json:"publicID"`
	RequestType       string            `json:"requestType"`
	ObjectType        string            `json:"docType"`
	RequesterPublicID string            `json:"requesterPublicID"`
	RecipientPublicID string            `json:"recipientPublicID"`
	RequestedData     string            `json:"requestedData"`
	AccountData       map[string]string `json:"accountData"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Status            string            `json:"status"`
}

type documentDataRequest struct {
	ID                string            `json:"id"`
	PublicID          string            `json:"publicID"`
	RequestType       string            `json:"requestType"`
	ObjectType        string            `json:"docType"`
	RequesterPublicID string            `json:"requesterPublicID"`
	RecipientPublicID string            `json:"recipientPublicID"`
	RequestedData     string            `json:"requestedData"`
	DocumentName      string            `json:"documentName"`
	DocumentData      map[string]string `json:"documentData"`
	DocumentCopy      bool              `json:"documentCopy"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Status            string            `json:"status"`
}

type ipfsDocumentVersionData struct {
	ContentIdentifier        string `json:"contentIdentifier"`
	ObjectHash               string `json:"objectHash"`
	Reference                string `json:"reference"`
	ParentDirectoryHash      string `json:"parentDirectoryHash"`
	ParentDirectoryReference string `json:"parentDirectoryReference"`
}

type documentVersion struct {
	ID        string                   `json:"id"`
	Name      int                      `json:"name"`
	IpfsData  *ipfsDocumentVersionData `json:"ipfsData"`
	CipherKey string                   `json:"cipherKey"`
	CreatedAt string                   `json:"createdAt"`
	UpdateAt  string                   `json:"updatedAt"`
}

type IPFSDirectoryData struct {
	ContentIdentifier   string `json:"contentIdentifier"`
	ObjectHash          string `json:"objectHash"`
	Reference           string `json:"reference"`
	LinkObjectHash      string `json:"linkObjectHash"`
	ParentDirectoryHash string `json:"parentDirectoryHash"`
}

type documentData struct {
	DocumentID   string `json:"documentID"`
	DocumentName string `json:"documentName"`
	Holder       string `json:"holder"`
	CountryIssue string `json:"countryIssue"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updateAd"`
}

type documentDirectory struct {
	ID                      string                   `json:"id"`
	ObjectType              string                   `json:"objectType"`
	DocumentData            string                   `json:"documentData"`
	IPFSDocumentData        string                   `json:"ipfsDocumentData"`
	IPFSDocumentVersionData map[int]*documentVersion `json:"ipfsDocumentVersionData"`
	UpdatedAt               string                   `json:"updatedAt"`
}

type accountData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type personalAccount struct {
	ID              string                        `json:"id"`
	PublicID        string                        `json:"publicID"`
	ObjectType      string                        `json:"objectType"`
	AccountData     *accountData                  `json:"accountData"`
	IPFSAccountData *IPFSDirectoryData            `json:"ipfsAccountData"`
	Documents       map[string]*documentDirectory `json:"documents"`
}
