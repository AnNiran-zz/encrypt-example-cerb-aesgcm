package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *CerberusPersonAccounts) queryAccountData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Start queryAccountData initialization.")

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	// get query function
	queryFunction := args[0]
	accountPublicID := args[1]

	switch queryFunction {

	case "getAccountHistory":
		return t.getAccountHistory(stub, []string{accountPublicID})

	case "getAccountRecords":
		return t.getAccountRecords(stub, []string{accountPublicID})

	default:
		return shim.Error("Function name not found.")
	}

	return shim.Success(nil)
}

func (t *CerberusPersonAccounts) queryRecords(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 1 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	// assign values
	selectorKey := args[0]
	selectorValue := strings.ToLower(args[1])

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"person\", \"accountData\":{\"%s\":\"%s\"}}}", selectorKey, selectorValue)

	// obtain records
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end queryRecords by: " + selectorKey + ": " + string(queryResults))
	return shim.Success(queryResults)
}

func (t *CerberusPersonAccounts) getAccountHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// assign values
	publicID := args[0]

	// obtain current records
	resultsIterator, err := stub.GetHistoryForKey(publicID)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the account
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON account)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAccountHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *CerberusPersonAccounts) getAccountRecords(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

	fmt.Println("- end getAccountRecords: " + string(queryResultBytes))
	return shim.Success(queryResultBytes)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)

	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

func (t *CerberusPersonAccounts) readAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, string, error) {

	var resultBytes []byte

	if len(args) < 1 {
		return nil, "", fmt.Errorf("Not enough arguments provided, Expecting 1")
	}

	if len(args[0]) <= 0 {
		return nil, "", fmt.Errorf("1st argument must be a non-empty string")
	}

	// assign values
	publicID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"person\",\"publicID\":\"%s\"}}", publicID)

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

	fmt.Println("- end readAccount")
	return resultBytes, string(buffer.Bytes()), nil
}

func getTime() string {

	currentDateTime := time.Now()
	CurrentDateTime := currentDateTime.Format("2006-01-02 15:04:05")

	return CurrentDateTime
}
