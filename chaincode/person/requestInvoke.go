package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// type, id, requesterId, recipientId, requested data
// type, id, requesterId, recipientId, documentName, requested data
func (t *CerberusPersonAccounts) createRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start createRequest initialization.")

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting at least 3")
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
	requestType := args[0]
	requestArgs := args[1:]

	switch requestType {
	case "accountData":
		return t.createAccountDataRequest(stub, requestArgs)

	case "documentData":
		return t.createDocumentDataRequest(stub, requestArgs)

	default:
		return shim.Error("Request type not found.")
	}
}

// id, requesterId, recipientId, requested data
func (t *CerberusPersonAccounts) createAccountDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestObject := args[0]
	requestData := args[1]

	newRequest := &accountDataRequest{}
	err := json.Unmarshal([]byte(requestObject), newRequest)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	_, err = t.checkRequestAttributes(stub, []string{newRequest.RequesterPublicID, newRequest.RecipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	// store requested data
	requestedFields := make(map[string]string)
	err = json.Unmarshal([]byte(requestData), &requestedFields)

	if err != nil {
		return shim.Error(err.Error())
	}

	requestedData := storeRequestedData(requestedFields)

	// check if same request already exists
	requestExistAsBytes, _, err := t.checkRequestImage(stub, []string{newRequest.RequesterPublicID, newRequest.RecipientPublicID, requestedData})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestExistAsBytes != nil {
		return shim.Error("Request with same image already exists")
	}

	// object -> finish creation
	newRequest.Status = "pending"
	newRequest.RequestedData = requestedData
	newRequest.AccountData = requestedFields
	newRequest.CreatedAt = getTime()

	requestAsBytes, err := json.Marshal(newRequest)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ledger invoke operation -> store with public id key
	err = stub.PutState(newRequest.PublicID, requestAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end createAccountDataRequest")
	return shim.Success(requestAsBytes)
}

func (t *CerberusPersonAccounts) createDocumentDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestObject := args[0]
	requestData := args[1]

	newRequest := &documentDataRequest{}
	err := json.Unmarshal([]byte(requestObject), newRequest)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{newRequest.RecipientPublicID, newRequest.RecipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check if requested document exists
	if _, ok := recipientAccount.Documents[newRequest.DocumentName]; !ok {
		return shim.Error("Document name " + newRequest.DocumentName + " does not exist in recipient account records and its data cannot be requested")
	}

	// store requested data
	requestedFields := make(map[string]string)
	err = json.Unmarshal([]byte(requestData), &requestedFields)

	if err != nil {
		return shim.Error(err.Error())
	}

	requestedData := storeRequestedData(requestedFields)

	// check if same request already exists
	requestExistAsBytes, _, err := t.checkRequestImage(stub, []string{newRequest.RequesterPublicID, newRequest.RecipientPublicID, requestedData})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestExistAsBytes != nil {
		return shim.Error("Request with same image already exists")
	}

	if _, ok := requestedFields["documentCopy"]; ok {
		newRequest.DocumentCopy = true
		delete(requestedFields, requestedFields["documentCopy"])
	}

	// object -> finish
	newRequest.Status = "pending"
	newRequest.RequestedData = requestedData
	newRequest.DocumentData = requestedFields
	newRequest.CreatedAt = getTime()

	requestDataAsBytes, err := json.Marshal(newRequest)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ledger invoke operation -> store with public id key
	err = stub.PutState(newRequest.PublicID, requestDataAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end createDocumentDataRequest")
	return shim.Success(requestDataAsBytes)
}

func (t *CerberusPersonAccounts) acceptRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start acceptRequest initialization.")

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting at least 4.")
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

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	// assign values
	requestType := args[0]
	requestArgs := args[1:]

	switch requestType {
	case "accountData":
		return t.acceptAccountDataRequest(stub, requestArgs)

	case "documentData":
		return t.acceptDocumentDataRequest(stub, requestArgs)

	default:
		return shim.Error("Request type not found.")
	}
}

func (t *CerberusPersonAccounts) acceptAccountDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestID := args[0]
	recipientPublicID := args[1]
	acceptedData := args[2]

	// obtain request object
	requestBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestBytes == nil {
		return shim.Error("Request with id: " + requestID + " does not exist")
	}

	request := &accountDataRequest{}
	err = json.Unmarshal(requestBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{request.RequesterPublicID, recipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	// get recipient account data
	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	if request.Status != "pending" {
		return shim.Error("Request status is already " + request.Status + " and cannot be accepted")
	}

	// obtain data
	fieldsData := recipientAccount.AccountData
	fieldsDataAsBytes, err := json.Marshal(fieldsData)

	if err != nil {
		return shim.Error(err.Error())
	}

	// umarshal account data as a map to create intersection
	var values map[string]string
	err = json.Unmarshal(fieldsDataAsBytes, &values)

	if err != nil {
		return shim.Error(err.Error())
	}

	acceptedFields := make(map[string]string)
	err = json.Unmarshal([]byte(acceptedData), &acceptedFields)

	if err != nil {
		return shim.Error(err.Error())
	}

	// match accepted fields values
	for field, value := range values { // data from retrieved account
		if _, ok := acceptedFields[field]; ok {
			acceptedFields[field] = value
		}
	}

	request.AccountData = acceptedFields
	request.Status = "accepted"
	request.UpdatedAt = getTime()

	acceptedRequestAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(requestID, acceptedRequestAsBytes) // store with public id again

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end acceptAccountDocumentRequest")
	return shim.Success(acceptedRequestAsBytes)
}

func (t *CerberusPersonAccounts) acceptDocumentDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestID := args[0]
	recipientPublicID := args[1]
	acceptedData := args[2]

	// obtain request object
	requestBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestBytes == nil {
		return shim.Error("Request with id: " + requestId + " does not exist")
	}

	request := &documentDataRequest{}
	err = json.Unmarshal(requestBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{request.RequesterPublicID, recipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	// get recipient account data
	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	if request.Status != "pending" {
		return shim.Error("Request status is already " + request.Status + " and cannot be accepted")
	}

	// obtain data
	fieldsData := recipientAccount.Documents[request.DocumentName].DocumentData
	fieldsDataAsBytes, err := json.Marshal(fieldsData)

	if err != nil {
		return shim.Error(err.Error())
	}

	// umarshal account data as a map to create intersection
	var values map[string]string
	err = json.Unmarshal(fieldsDataAsBytes, &values)

	if err != nil {
		return shim.Error(err.Error())
	}

	acceptedFields := make(map[string]string)
	err = json.Unmarshal([]byte(acceptedData), &acceptedFields)

	if err != nil {
		return shim.Error(err.Error())
	}

	// match accepted fields values
	for field, value := range values { // data from retrieved account
		if _, ok := acceptedFields[field]; ok {
			acceptedFields[field] = value
		}
	}

	request.DocumentData = acceptedFields
	request.Status = "accepted"
	request.UpdatedAt = getTime()

	acceptedRequestAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(requestID, acceptedRequestAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end acceptDocumentDataRequest")
	return shim.Success(acceptedRequestAsBytes)
}

func (t *CerberusPersonAccounts) rejectRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start rejectRequest initialization.")

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
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
	requestType := args[0]
	requestArgs := args[1:]

	switch requestType {
	case "accountData":
		return t.rejectAccountDataRequest(stub, requestArgs)

	case "documentData":
		return t.rejectDocumentDataRequest(stub, requestArgs)

	default:
		return shim.Error("Request type not recognized.")
	}
}

func (t *CerberusPersonAccounts) rejectAccountDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestD := args[0]
	recipientPublicID := args[1]

	// obtain request object
	requestBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestBytes == nil {
		return shim.Error("Request with ID: " + requestID + " does not exist")
	}

	request := &accountDataRequest{}
	err = json.Unmarshal(requestBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{request.RequesterPublicId, recipientPublicId})

	if err != nil {
		return shim.Error(err.Error())
	}

	// get recipient account data
	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	if request.Status != "pending" {
		return shim.Error("Request status is already " + request.Status + " and cannot be rejected")
	}

	request.Status = "rejected"
	request.UpdatedAt = getTime()

	requestUpdateAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(requestId, requestUpdateAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end rejectAccountDataRequest")
	return shim.Success(requestUpdateAsBytes)
}

func (t *CerberusPersonAccounts) rejectDocumentDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestID := args[0]
	recipientPublicID := args[1]

	// check if request exists
	queryResultBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes == nil {
		return shim.Error("No requests with ID: " + requestID + " exist.")
	}

	request := &documentDataRequest{}
	err = json.Unmarshal(queryResultBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{request.RequesterPublicID, recipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	// get recipient account data
	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	if request.Status != "pending" {
		return shim.Error("Request status is already " + request.Status + " and cannot be rejected")
	}

	request.Status = "rejected"
	request.UpdatedAt = getTime()

	requestUpdateAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(requestID, requestUpdateAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end rejectDocumentDataRequest")
	return shim.Success(requestUpdateAsBytes)
}

func (t *CerberusPersonAccounts) updateRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start updateRequest initialization.")

	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5.")
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

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	// assign values
	requestType := args[0]
	requestArgs := args[1:]

	switch requestType {
	case "accountData":
		return t.updateAccountDataRequest(stub, requestArgs)

	case "documentData":
		return t.updateDocumentDataRequest(stub, requestArgs)

	default:
		return shim.Error("Unknown request type")
	}
}

func (t *CerberusPersonAccounts) updateAccountDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start updateAccountDataRequest initialization.")

	requestID := args[0]
	requesterPublicID := args[1]
	recipientPublicID := args[2]
	data := args[3]

	// obtain request object
	requestBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestBytes == nil {
		return shim.Error("Request with id: " + requestID + " does not exist")
	}

	request := &accountDataRequest{}
	err = json.Unmarshal(requestBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{requesterPublicID, recipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	// check request status
	if request.Status != "pending" {
		fmt.Println("Request with ID " + request.PublicID + " has been " + request.Status + " and cannot be updated")
		return shim.Success(nil)
	}

	// store requested data
	requestedFields := make(map[string]string)
	err = json.Unmarshal([]byte(data), &requestedFields)

	if err != nil {
		return shim.Error(err.Error())
	}

	requestedData := storeRequestedData(requestedFields)

	// check if request already exists
	requestExistAsBytes, _, err := t.checkRequestImage(stub, []string{requesterPublicID, recipientPublicID, requestedData})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestExistAsBytes != nil {
		return shim.Error("Request with same image already exists.")
	}

	// object -> update
	request.RequestedData = requestedData
	request.AccountData = requestedFields
	request.UpdatedAt = getTime()

	requestAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ledger invoke operation -> store with public id key
	err = stub.PutState(requestID, requestAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateAccountDataRequest")
	return shim.Success(requestAsBytes)
}

func (t *CerberusPersonAccounts) updateDocumentDataRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start updateDocumentDataRequest initialization.")

	requestID := args[0]
	requesterPublicID := args[1]
	recipientPublicID := args[2]
	data := args[3]

	// obtain request object
	requestBytes, _, err := t.readRequest(stub, []string{"publicID", requestID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestBytes == nil {
		return shim.Error("Request with id: " + requestID + " does not exist")
	}

	request := &documentDataRequest{}
	err = json.Unmarshal(requestBytes, request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check attributes
	recipientAccountAsBytes, err := t.checkRequestAttributes(stub, []string{requesterPublicID, recipientPublicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	recipientAccount := &personAccount{}
	err = json.Unmarshal(recipientAccountAsBytes, recipientAccount)

	if err != nil {
		return shim.Error(err.Error())
	}

	if recipientAccount.PublicID != request.RecipientPublicID {
		return shim.Error("Request recipient ID does not match the provided ID")
	}

	// check request status
	if request.Status != "pending" {
		fmt.Println("Request with id " + request.PublicID + " has been " + request.Status + " and cannot be updated")
		return shim.Success(nil)
	}

	// data fields -> update
	requestedFields, err := updateFields(request.RequestedData, data)

	if err != nil {
		return shim.Error(err.Error())
	}

	// check if data has been already requested
	requestedData := storeRequestedData(requestedFields)

	requestExistAsBytes, _, err := t.checkRequestImage(stub, []string{requesterPublicID, recipientPublicID, requestedData})

	if err != nil {
		return shim.Error(err.Error())
	}

	if requestExistAsBytes != nil {
		return shim.Error("Request with same image already exists.")
	}

	// object -> update
	request.RequestedData = requestedData
	request.DocumentData = requestedFields
	request.UpdatedAt = getTime()

	requestAsBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ledger invoke operation -> store with public id key
	err = stub.PutState(requestID, requestAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateDocumentDataRequest")
	return shim.Success(requestAsBytes)
}

func (t *CerberusPersonAccounts) checkRequestAttributes(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	// assign values
	requesterPublicID := args[0]
	recipientPublicID := args[1]

	// check if requester exists
	requesterDataAsBytes, _, err := t.readAccount(stub, []string{requesterPublicID})

	if err != nil {
		return nil, err
	}

	if requesterDataAsBytes == nil {
		return nil, errors.New("Account with ID (for requester): " + requesterPublicID + " does not exist")
	}

	// check if recipient account exists
	recipientAccountAsBytes, _, err := t.readAccount(stub, []string{recipientPublicID})

	if err != nil {
		return nil, err
	}

	if recipientAccountAsBytes == nil {
		return nil, errors.New("Account with ID (for recipient): " + recipientPublicID + " does not exist")
	}

	return recipientAccountAsBytes, nil
}

func storeRequestedData(data map[string]string) string {

	var content string

	alphabeticKeys := make([]string, len(data))

	i := 0
	for key, _ := range data {
		alphabeticKeys[i] = key
		i++
	}

	sort.Strings(alphabeticKeys)

	for _, value := range alphabeticKeys {
		content += value + "+"
	}

	content = content[:len(content)-len("+")]

	return content
}

func updateFields(currentData, newData string) (map[string]string, error) {

	updateFields := make(map[string]string)

	err := json.Unmarshal([]byte(newData), &updateFields)

	if err != nil {
		return nil, err
	}

	currentDataSplit := strings.Split(currentData, "+")

	currentFields := make(map[string]string)
	for _, value := range currentDataSplit {
		currentFields[value] = ""
	}

	//newResult := make(map[string]string)
	for key, _ := range updateFields {

		if _, ok := currentFields[key]; !ok {
			currentFields[key] = ""
		}
	}

	return currentFields, nil
}
