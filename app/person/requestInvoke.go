package person

import (
	"cerberus/blockchain/instaccntschannel"
	"cerberus/blockchain/persaccntschannel"
	"cerberus/services/crypto"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func CreateAccountDataRequest(requesterPublicId, recipientPublicId string, args []string) (string, []byte, error) {

	if requesterPublicId == "" {
		return "", nil, errors.New("Requester ID value cannot be an empty string")
	}

	if recipientPublicId == "" {
		return "", nil, errors.New("Recipient ID value cannot be an empty string")
	}

	// x arguments
	if len(args) < 1 {
		return "", nil, errors.New("Insufficient number of arguments Expecting at least 1")
	}

	// check recipient and requester ids
	if requesterPublicId == recipientPublicId {
		return "", nil, errors.New("Recipient and requester ids cannot be identical")
	}

	// request -> create
	requestDataSample := &accountShareableData{}
	requestedFields := make(map[string]string)

	for key, argument := range args {

		if len(argument) <= 0 {
			return "", nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
		}

		// check if argument match any fields from data share list
		value := reflect.ValueOf(requestDataSample).Elem().FieldByName(argument)

		if !value.IsValid() {
			return "", nil, errors.New("Field value: " + argument + " does not match acceptable data fields.")
		}

		//value.SetString(argument)
		requestedFields[argument] = ""
	}

	requestDataAsBytes, err := json.Marshal(requestedFields)
	if err != nil {
		return "", nil, err
	}

	// create object
	id := bson.NewObjectId().Hex()
	publicId := crypto.GetMD5Hash(id)

	newRequest := &accountDataRequest{
		Id:                id,
		PublicId:          publicId,
		RequestType:       "accountData",
		ObjectType:        "persAccntsRequest",
		RequesterPublicId: requesterPublicId,
		RecipientPublicId: recipientPublicId,
		AccountData:       make(map[string]string),
	}

	newRequestAsBytes, err := json.Marshal(newRequest)
	if err != nil {
		return "", nil, err
	}

	// call chanicode function that sends the request to the data holder
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	_, requestRecord, err := persAccntsChannelClient.CreateAccountDataRequest(newRequestAsBytes, requestDataAsBytes)
	if err != nil {
		return "", nil, err
	}

	request := &accountDataRequest{}
	if err = json.Unmarshal([]byte(requestRecord), request); err != nil {
		return "", nil, err
	}

	return request.PublicId, requestRecord, nil
}

func CreateDocumentDataRequest(requesterPublicId, recipientPublicId, documentName string, args []string, documentCopy bool) (string, []byte, error) {

	if requesterPublicId == "" {
		return "", nil, errors.New("Requester ID value cannot be an empty string")
	}

	if recipientPublicId == "" {
		return "", nil, errors.New("Recipient ID value cannot be an empty string")
	}

	if documentName == "" {
		return "", nil, errors.New("Document name value cannot be an empty string")
	}

	// check recipient and requester ids
	if requesterPublicId == recipientPublicId {
		return "", nil, errors.New("Recipient and requester ids cannot be identical")
	}

	// request -> create
	requestDataSample := &documentShareableData{}
	requestedFields := make(map[string]string)

	// x arguments
	if len(args) >= 1 {
		for key, argument := range args {

			if len(argument) <= 0 {
				return "", nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
			}

			// check if argument match any fields from data share list
			value := reflect.ValueOf(requestDataSample).Elem().FieldByName(argument)

			if !value.IsValid() {
				return "", nil, errors.New("Field value: " + argument + " does not match acceptable data fields.")
			}

			requestedFields[argument] = ""
		}
	}

	if documentCopy == true {
		requestedFields["documentCopy"] = ""
	}

	requestDataAsBytes, err := json.Marshal(requestedFields)
	if err != nil {
		return "", nil, err
	}

	// create object
	id := bson.NewObjectId().Hex()
	publicId := crypto.GetMD5Hash(id)

	newRequest := &documentDataRequest{
		Id:                id,
		PublicId:          publicId,
		RequestType:       "documentData",
		ObjectType:        "persAccntsRequest",
		RequesterPublicId: requesterPublicId,
		RecipientPublicId: recipientPublicId,
		DocumentName:      strings.ToLower(documentName),
		DocumentData:      make(map[string]string),
	}

	newRequestAsBytes, err := json.Marshal(newRequest)
	if err != nil {
		return "", nil, err
	}

	// call chanicode function that sends the request to the data holder
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	_, requestRecord, err := persAccntsChannelClient.CreateDocumentDataRequest(newRequestAsBytes, requestDataAsBytes)

	if err != nil {
		return "", nil, err
	}

	request := &documentDataRequest{}
	if err = json.Unmarshal(requestRecord, request); err != nil {
		return "", nil, err
	}

	return request.PublicId, requestRecord, nil
}

func AcceptAccountDataRequest(recipientPublicId, requestPublicId string, args []string) ([]string, []string, error) {

	if recipientPublicId == "" {
		return nil, nil, errors.New("Account Id value cannot be an empty string")
	}

	if requestPublicId == "" {
		return nil, nil, errors.New("Request Id value cannot be an empty string")
	}

	// x arguments
	if len(args) < 1 {
		return nil, nil, errors.New("Insufficient number of arguments Expecting at least 1")
	}

	for key, argument := range args {

		if len(argument) <= 0 {
			return nil, nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
		}
	}

	// get request from blockchain
	requestData, err := GetRequestObject("publicId", requestPublicId)
	if err != nil {
		return nil, nil, err
	}

	request := &accountDataRequest{}
	if err = json.Unmarshal([]byte(requestData), request); err != nil {
		return nil, nil, err
	}

	// match accepted fields
	acceptedFields := make(map[string]string)
	acceptedFields = GetIntersection(request.AccountData, args)

	acceptedFieldsAsBytes, err := json.Marshal(acceptedFields)
	if err != nil {
		return nil, nil, err
	}

	// send request
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, record, err := persAccntsChannelClient.AcceptRequest("accountData", requestPublicId, recipientPublicId, acceptedFieldsAsBytes)
	if err != nil {
		return nil, nil, err
	}

	return response, []string{string(record)}, nil
}

func AcceptDocumentDataRequest(recipientPublicId, requestPublicId string, documentCopy string, args []string) ([]string, []string, []string, error) {

	if recipientPublicId == "" {
		return nil, nil, nil, errors.New("Account Id value cannot be an empty string")
	}

	if requestPublicId == "" {
		return nil, nil, nil, errors.New("Request Id value cannot be an empty string")
	}

	// x arguments
	if len(args) < 1 {
		return nil, nil, nil, errors.New("Insufficient number of arguments Expecting at least 1")
	}

	for key, argument := range args {

		if len(argument) <= 0 {
			return nil, nil, nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
		}
	}

	// get request from blockchain
	requestData, err := GetRequestObject("publicId", requestPublicId)
	if err != nil {
		return nil, nil, nil, err
	}

	request := &documentDataRequest{}
	if err = json.Unmarshal([]byte(requestData), request); err != nil {
		return nil, nil, nil, err
	}

	// match accepted fields
	acceptedFields := make(map[string]string)
	acceptedFields = GetIntersection(request.DocumentData, args)

	acceptedFieldsAsBytes, err := json.Marshal(acceptedFields)

	if err != nil {
		return nil, nil, nil, err
	}

	var documentCopyData []string

	if request.DocumentCopy == true {
		if documentCopy != "" {
			documentCopyData, err = GetAccountDocumentVersion(recipientPublicId, request.DocumentName, documentCopy)

			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	// send request
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, record, err := persAccntsChannelClient.AcceptRequest("documentData", requestPublicId, recipientPublicId, acceptedFieldsAsBytes)
	if err != nil {
		return nil, nil, nil, err
	}

	return response, []string{string(record)}, documentCopyData, nil
}

func RejectAccountDataRequest(recipientPublicId, requestPublicId string) ([]string, []string, error) {

	if recipientPublicId == "" {
		return nil, nil, errors.New("Account Id value cannot be an empty string")
	}

	if requestPublicId == "" {
		return nil, nil, errors.New("Request Id value cannot be an empty string")
	}

	// send request
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, record, err := persAccntsChannelClient.RejectRequest("accountData", requestPublicId, recipientPublicId)
	if err != nil {
		return nil, nil, err
	}

	return response, []string{string(record)}, nil
}

func RejectDocumentDataRequest(recipientPublicId, requestPublicId string) ([]string, []string, error) {

	if recipientPublicId == "" {
		return nil, nil, errors.New("Account Id value cannot be an empty string")
	}

	if requestPublicId == "" {
		return nil, nil, errors.New("Request Id value cannot be an empty string")
	}

	// send request
	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	response, record, err := persAccntsChannelClient.RejectRequest("documentData", requestPublicId, recipientPublicId)
	if err != nil {
		return nil, nil, err
	}

	return response, []string{string(record)}, nil
}

func UpdateAccountDataRequest(requesterPublicId, recipientPublicId, requestPublicId string, args []string) (string, []byte, error) {

	if requesterPublicId == "" {
		return "", nil, errors.New("Requester ID value cannot be an empty string")
	}

	if recipientPublicId == "" {
		return "", nil, errors.New("Recipient ID value cannot be an empty string")
	}

	// x arguments
	if len(args) < 1 {
		return "", nil, errors.New("Insufficient number of arguments Expecting at least 1")
	}

	// check recipient and requester ids
	if requesterPublicId == recipientPublicId {
		return "", nil, errors.New("Recipient and requester ids cannot be identical")
	}

	requestData, err := GetRequestObject("publicId", requestPublicId)
	if err != nil {
		return "", nil, err
	}

	accountRequest := &accountDataRequest{}
	if err = json.Unmarshal([]byte(requestData), accountRequest); err != nil {
		return "", nil, err
	}

	if accountRequest.Status != "pending" {

		fmt.Println("Request status is " + accountRequest.Status)
		fmt.Println("Creating new request")

		publicId, requestRecord, err := CreateAccountDataRequest(requesterPublicId, recipientPublicId, args)

		if err != nil {
			return "", nil, err
		}

		return publicId, requestRecord, nil
	}

	// validate fields
	requestDataSample := &accountShareableData{}
	requestedFields := make(map[string]string)

	for key, argument := range args {

		if len(argument) <= 0 {
			return "", nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
		}

		// check if argument match any fields from data share list
		value := reflect.ValueOf(requestDataSample).Elem().FieldByName(argument)

		if !value.IsValid() {
			return "", nil, errors.New("Field value: " + argument + " does not match acceptable data fields.")
		}

		//value.SetString(argument)
		requestedFields[argument] = ""
	}

	updateDataAsBytes, err := json.Marshal(requestedFields)
	if err != nil {
		return "", nil, err
	}

	// send updates to ledger
	instAccntsChannelClient := instaccntschannel.CerberusClient{}
	_, updateRecord, err := instAccntsChannelClient.UpdateRequest("accountData", requestPublicId, requesterPublicId, recipientPublicId, updateDataAsBytes)
	if err != nil {
		return "", nil, err
	}

	request := &accountDataRequest{}
	if err = json.Unmarshal([]byte(updateRecord), request); err != nil {
		return "", nil, err
	}

	// obtain current
	return request.PublicId, updateRecord, nil
}

func UpdateDocumentDataRequest(requesterPublicId, recipientPublicId, requestPublicId, documentName string, args []string, documentCopy bool) (string, []byte, error) {

	if requesterPublicId == "" {
		return "", nil, errors.New("Requester ID value cannot be an empty string")
	}

	if recipientPublicId == "" {
		return "", nil, errors.New("Recipient ID value cannot be an empty string")
	}

	if documentName == "" {
		return "", nil, errors.New("Document name value cannot be an empty string")
	}

	// check recipient and requester ids
	if requesterPublicId == recipientPublicId {
		return "", nil, errors.New("Recipient and requester ids cannot be identical")
	}

	requestData, err := GetRequestObject("publicId", requestPublicId)
	if err != nil {
		return "", nil, err
	}

	documentRequest := &documentDataRequest{}
	if err = json.Unmarshal([]byte(requestData), documentRequest); err != nil {
		return "", nil, err
	}

	if documentRequest.Status != "pending" {

		fmt.Println("Request status is " + documentRequest.Status)
		fmt.Println("Creating new request")

		publicId, requestRecord, err := CreateDocumentDataRequest(requesterPublicId, recipientPublicId, documentName, args, documentCopy)

		if err != nil {
			return "", nil, err
		}

		return publicId, requestRecord, nil
	}

	// request -> create
	requestDataSample := &documentShareableData{}
	requestedFields := make(map[string]string)

	// x arguments
	if len(args) >= 1 {
		for key, argument := range args {

			if len(argument) <= 0 {
				return "", nil, errors.New("Argument: " + strconv.Itoa(key) + " cannot be empty.")
			}

			// check if argument match any fields from data share list
			value := reflect.ValueOf(requestDataSample).Elem().FieldByName(argument)

			if !value.IsValid() {
				return "", nil, errors.New("Field value: " + argument + " does not match acceptable data fields.")
			}

			requestedFields[argument] = ""
		}
	}

	if documentCopy == true {
		requestedFields["documentCopy"] = ""
	}

	updateDataAsBytes, err := json.Marshal(requestedFields)
	if err != nil {
		return "", nil, err
	}

	// send updates to ledger
	instAccntsChannelClient := instaccntschannel.CerberusClient{}
	_, updateRecord, err := instAccntsChannelClient.UpdateRequest("accountData", requestPublicId, requesterPublicId, recipientPublicId, updateDataAsBytes)

	if err != nil {
		return "", nil, err
	}

	request := &accountDataRequest{}
	if err = json.Unmarshal([]byte(updateRecord), request); err != nil {
		return "", nil, err
	}

	// obtain current
	return request.PublicId, updateRecord, nil
}

func GetIntersection(a map[string]string, b []string) map[string]string {

	valuesMap := make(map[string]bool)
	//var result []string
	intersectionMap := make(map[string]string)

	for item, _ := range a {
		valuesMap[item] = true
	}

	for _, item := range b {
		if _, ok := valuesMap[item]; ok {
			//result = append(result, item)
			intersectionMap[item] = ""
		}
	}

	return intersectionMap
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
