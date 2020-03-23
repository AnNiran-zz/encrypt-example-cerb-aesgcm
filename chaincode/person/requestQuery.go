package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *CerberusPersonAccounts) queryRequestData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start queryRequestData a initialization.")

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	idType := args[0]
	publicID := args[1]

	// check if request exists
	queryResultBytes, _, err := t.readRequest(stub, []string{idType, publicID})

	if err != nil {
		return shim.Error(err.Error())
	}

	if queryResultBytes == nil {
		return shim.Error("No requests with " + idType + " : " + publicID + " exist.")
	}

	fmt.Println("- end queryRequestData: " + string(queryResultBytes))
	return shim.Success(queryResultBytes)
}

func (t *CerberusPersonAccounts) readRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, string, error) {

	var resultBytes []byte

	if len(args) < 2 {
		return nil, "", fmt.Errorf("Not enough arguments provided")
	}

	if len(args[0]) <= 0 {
		return nil, "", fmt.Errorf("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return nil, "", fmt.Errorf("2nd argument must be a non-empty string")
	}

	// assign values
	idType := args[0]
	ID := args[1]

	var queryString string

	switch idType {
	case "publicID":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"publicID\":\"%s\"}}", ID)

	case "requestID":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"ID\":\"%s\"}}", ID)

	default:
		return nil, "", fmt.Errorf("Unknown ID type")
	}

	// obtain records
	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return nil, "", err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, "", err
		}

		resultBytes = response.Value
	}

	buffer, err := constructQueryResponseFromIterator(resultsIterator)

	if err != nil {
		return nil, "", err
	}

	fmt.Println("- end readRequest")
	return resultBytes, string(buffer.Bytes()), nil
}

func (t *CerberusPersonAccounts) queryRequests(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start queryRequests initialization")

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 1 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	// assign values
	queryType := args[0]
	queryArgs := args[1:]

	switch queryType {
	case "objects":
		return t.queryRequestsObjects(stub, queryArgs)

	case "publicIDs":
		return t.queryRequestsPublicIDs(stub, queryArgs)

	default:
		return shim.Error("Query type not found")
	}
}

func (t *CerberusPersonAccounts) queryRequestsObjects(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// assign values
	requestType := args[0]
	selectorKey := args[1]
	selectorValue := strings.ToLower(args[2])

	var queryString string

	switch requestType {
	case "accountData":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"requestType\":\"accountData\",\"%s\":\"%s\"}}", selectorKey, selectorValue)

	case "documentData":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"requestType\":\"documentData\",\"%s\":\"%s\"}}", selectorKey, selectorValue)

	case "any":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"%s\":\"%s\"}}", selectorKey, selectorValue)

	default:
		return shim.Error("Unknown request type")
	}

	// obtain records
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end queryRequestsObjects by: " + selectorKey + ": " + string(queryResults))
	return shim.Success(queryResults)
}

func (t *CerberusPersonAccounts) queryRequestsPublicIDs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// assign values
	requestType := args[0]
	selectorKey := args[1]
	selectorValue := strings.ToLower(args[2])

	var queryString string

	switch requestType {
	case "accountData":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"requestType\":\"accountData\",\"%s\":\"%s\"},\"fields\":[\"publicId\"]}", selectorKey, selectorValue)

	case "documentData":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"requestType\":\"documentData\",\"%s\":\"%s\"},\"fields\":[\"publicId\"]}", selectorKey, selectorValue)

	case "any":
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"%s\":\"%s\"},\"fields\":[\"publicId\"]}", selectorKey, selectorValue)

	default:
		return shim.Error("Unknown request type")

	}

	// obtain records
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end queryRequestsPublicIDs by: " + selectorKey + ": " + string(queryResults))
	return shim.Success(queryResults)
}

func (t *CerberusPersonAccounts) checkRequestImage(stub shim.ChaincodeStubInterface, args []string) ([]byte, string, error) {

	if len(args) < 3 {
		return nil, "", fmt.Errorf("Not enough arguments provided")
	}

	if len(args[0]) <= 0 {
		return nil, "", fmt.Errorf("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return nil, "", fmt.Errorf("2nd argument must be a non-empty string")
	}

	if len(args[2]) <= 0 {
		return nil, "", fmt.Errorf("3rd argument must be a non-empty string")
	}

	requesterPublicID := args[0]
	recipientPublicID := args[1]
	requestedData := args[2]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"persAccntsRequest\",\"requesterPublicID\":\"%s\",\"recipientPublicID\":\"%s\",\"requestedData\":\"%s\",\"status\":\"pending\"}}", requesterPublicID, recipientPublicID, requestedData)

	// obtain records
	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return nil, "", err
	}
	defer resultsIterator.Close()

	var resultBytes []byte

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, "", err
		}

		resultBytes = response.Value
	}

	buffer, err := constructQueryResponseFromIterator(resultsIterator)

	if err != nil {
		return nil, "", err
	}

	fmt.Println("- end readRequest")
	return resultBytes, string(buffer.Bytes()), nil
}
