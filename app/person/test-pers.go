package person

import (
	"fmt"
	"os"
)

func TestPers() {

	//directory, err := ipfs.CreateGroupAccountsIpfsDirectory("personAccounts")
	//fmt.Println(directory)
	//fmt.Println(err)

	filename := os.Getenv("GOPATH") + "/src/cerberus/app/create-account-personal.png"
	fmt.Println(filename)

	id1 := "34a7570d84984f147f2cad26ff33d4fd" // document holder
	fmt.Println(id1)

	id2 := "649ba245a66626e80de07b35b7340b4f"
	fmt.Println(id2)

	id3 := "d5d4b155b92b611ce6f2704d06482844"
	fmt.Println(id3)

	//response, record, err := CreateAccount("anna", "AngeLOVA", "angeloWWA@gmail.COM", "123456")
	//data, response, err := UpdateAccountBySelector(id1, "Phone", "newName")
	//data, response, err := UpdateAccountFirstName(id2, "MyNENAMEEEE")
	//data, response, err := UpdateAccountLastName(id2, "myNELASTNAMEEE")
	//data, response, err := UpdateAccountPhone(id1, "123")

	//fmt.Println(response)
	//fmt.Println(data)
	//fmt.Println(record)
	//fmt.Println(err)

	//result, err := DeleteAccount(id3)

	//fmt.Println(result)
	//fmt.Println(err)

	//record, err := GetAccountById(id2)
	//record, err := GetAccountsByEmail("angelowwa@gmail.co")
	//record, err := GetAccountsByFirstName("newname")
	//record, err := GetAccountsByLastName("angelova")
	//record, err := GetAccountHistory(id1)
	//record, err := GetAccountsBySelector("email", "angelowwa@gmail.com")

	//fmt.Println(record)
	//fmt.Println(err)

	//rsaLink := "/hdd/server/go/src/cerberus/ipfs/personAccounts/2da281036a6febef53279395d5ecb59f/newdocument2/rsa/rsa_key.pem"
	//record, response, rsaLink, err := CreateNewDocument(id1, "newDocuMENT", "anna", "bulgaria", filename)
	//record, response, err := CreateDocumentVersion(id1, "newdocumenT2", filename, rsaLink)

	//fmt.Println(record)
	//fmt.Println(response)
	//fmt.Println(rsaLink)
	//fmt.Println(err)

	//record, response, err = UpdateDocumentHolderName(id1, "newDocument2", "newHoldeName")
	//record, response, err := UpdateDocumentCountryIssue(id1, "newdocument2", "newcountryIssue")
	//record, response, err := DeleteDocumentVersion(id1, "newdocument2", 3)
	//record, response, err := DeleteDocument(id1, "newdocument2")

	//fmt.Println(record)
	//fmt.Println(response)
	//fmt.Println(err)

	//result, err := GetAccountDocument(id1, "newdocument2")
	//result, err := GetAccountDocumentVersion(id1, "newdocument2", "5")
	//result, err := GetAccountDocumentVersions(id1, "newdocument2")

	//fmt.Println(result)
	//fmt.Println(err)

	// *********************************

	//publicId, result, err := CreateAccountDataRequest(id2, id1, []string{"firstName", "phone"})
	//response, record, err := RejectAccountDataRequest(id1, "e42ccff1bd28a50040e8ba531f6ddb78")
	//response, record, err := AcceptAccountDataRequest(id1, "ef4c06f241f4b9013cbd01eb09043bb4", []string{"firstName", "lastName"})

	//publicId, record, err := UpdateAccountDataRequest(id2, id1, "ef4c06f241f4b9013cbd01eb09043bb4", []string{"firstName"})
	//publicId, record, err := UpdateDocumentDataRequest(id2, id1, "4c53b2fd546bad2348d4c62d1b1c5dca", "newdocument", []string{"documentName"}, false)

	//fmt.Println(publicId)
	//fmt.Println(string(result))
	//fmt.Println(record)
	//fmt.Println(err)

	//fmt.Println(response)
	//fmt.Println(record)

	//id, record, err := CreateDocumentDataRequest(id2, id1, "newdocument", []string{"holder", "countryIssue", "documentName"}, false)
	response, record, _, err := AcceptDocumentDataRequest(id1, "413f6155ff15e1fbd30470dcdfb053b4", "1", []string{"holder", "countryIssue", "documentName"})

	//fmt.Println(id)
	//fmt.Println(string(record))

	fmt.Println(response)
	fmt.Println(record)
	fmt.Println(err)

	//response, record, err := RejectDocumentDataRequest(id1, "4c53b2fd546bad2348d4c62d1b1c5dca")

	//fmt.Println(response)
	//fmt.Println(record)
	//fmt.Println(err)

	// ***************************************
	//data, err := GetRequestsObjectsBySelector("any", "status", "rejected")
	//data, err := GetRequestsPublicIdsBySelector("documentData", "status", "rejected")
	//data, err := GetRequestsByRecipient("objects", "documentData", id2)
	//data, err := GetRequestsByRequester("publicIds", "any", id1)
	//data, err := GetRequestsByDocumentName("publicIds", "newdocument")
	//data, err := GetRequestsByStatus("objects", "any", "pending")
	//data, err := GetRequestObject("publicId", "4c53b2fd546bad2348d4c62d1b1c5dca")

	//fmt.Println(data)
	//fmt.Println(err)
}
