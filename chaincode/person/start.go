package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type CerberusPersonAccounts struct{}

func (t *CerberusPersonAccounts) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("Cerberus chaincode instantiation.")

	_, args := stub.GetFunctionAndParameters()

	t.createAccount(stub, args)
	return shim.Success(nil)
}

func (t *CerberusPersonAccounts) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running " + function)

	// Handle different functions
	switch function {

	case "createAccount":
		return t.createAccount(stub, args)

	// queryAccountData: getAccountHistory, getAccountRecords
	case "queryAccountData":
		return t.queryAccountData(stub, args)

	case "queryRecords":
		return t.queryRecords(stub, args)

	// updateRecords : updateAccount, updateDocumentRecords
	case "updateRecords":
		return t.updateRecords(stub, args)

	case "deleteAccount":
		return t.deleteAccount(stub, args)

	// requests
	case "createRequest":
		return t.createRequest(stub, args)

	case "queryRequestData":
		return t.queryRequestData(stub, args)

	case "acceptRequest":
		return t.acceptRequest(stub, args)

	case "rejectRequest":
		return t.rejectRequest(stub, args)

	case "updateRequest":
		return t.updateRequest(stub, args)

	case "queryRequests":
		return t.queryRequests(stub, args)

	default:
		return shim.Error("Function name not found.")
	}

	return shim.Success(nil)
}

func main() {

	err := shim.Start(new(CerberusPersonAccounts))

	if err != nil {
		fmt.Println("Error starting Person chaincode: %s", err.Error())
	}
}
