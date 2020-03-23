package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *CerberusPersonAccounts) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start Person account initialization.")

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	publicID := args[0]
	accountObject := args[1]

	// check if account exists
	queryResultBytes, _, err := t.readAccount(stub, []string{publicID})
	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes != nil {
		return shim.Error("Record with public id: " + publicID + " already exists.")
	}

	// ledger invoke operation
	err = stub.PutState(publicID, accountObject)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end createAccount")
	return shim.Success(publicID)
}

func (t *CerberusPersonAccounts) updateRecords(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Initialize updateRecords")

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}

	// assign values
	updateFunction := args[0]
	updateArgs := args[1:]

	switch updateFunction {

	case "updateAccount":
		return t.updateAccount(stub, updateArgs)

	case "updateDocumentRecords":
		return t.updateDocumentRecords(stub, updateArgs)

	default:
		return shim.Error("Function name not found.")
	}
}

func (t *CerberusPersonAccounts) updateAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// assign values
	publicID := args[0]
	passphrase := args[1]
	dataField := args[2]
	updateValue := args[3]

	// check if account exists
	queryResultBytes, _, err := t.readAccount(stub, []string{publicID})
	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes == nil {
		return shim.Error("No records with provided id exist.")
	}

	// object -> get
	currentRecord, err := decrAESGCM(queryResultBytes, []byte(passphrase)
	if err != nil {
		return shim.Error(err.Error())
	}

	recordUpdate := &personAccount{}
	err = json.Unmarshal(currentRecord, recordUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// object -> update
	value := reflect.ValueOf(recordUpdate.AccountData).Elem().FieldByName(dataField)
	if value.IsValid() {
		value.SetString(updateValue)
	}

	recordUpdate.AccountData.UpdatedAt = getTime()
	recordUpdateAsBytes, err := json.Marshal(recordUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// encrypt again
	encrRecord, err := encrAESGCM(recordUpdateAsBytes, passphrase)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ledger invoke operation
	err = stub.PutState(publicID, encrRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateAccount: ")
	return shim.Success()
}

func (t *CerberusPersonAccounts) deleteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// input sanitation
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	// assign values
	publicID := args[0]

	// check if account exists
	queryResultBytes, _, err := t.readAccount(stub, []string{publicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes == nil {
		return shim.Error("No records with provided id exist.")
	}

	// ledger invoke operation
	err = stub.DelState(publicID)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end deleteAccount")
	return shim.Success(queryResultBytes)
}

func (t *CerberusPersonAccounts) updateDocumentRecords(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// assign values
	publicID := args[0]
	data := args[2]

	// check if account exists
	queryResultBytes, _, err := t.readAccount(stub, []string{publicID})
	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes == nil {
		return shim.Error("No records with provided id exist.")
	}

	// ledger invoke operation
	err = stub.PutState(publicID, []byte(data))

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateDocumentRecords: " + string(data))
	return shim.Success([]byte(data))
}
