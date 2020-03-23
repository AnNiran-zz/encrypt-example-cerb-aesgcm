package person

import (
	"cerberus/blockchain/persaccntschannel"
	"encoding/json"
	"errors"
)

/*
selectorKeys:
- recipientId
- requesterId
- documentName
- status
- requestType
*/
func GetRequestsObjectsBySelector(requestType, selectorKey, selectorValue string) ([]string, error) {

	if requestType == "" {
		return nil, errors.New("Request type value cannot be an empty string")
	}

	if selectorKey == "" {
		return nil, errors.New("Selector key value cannot be an empty string")
	}

	if selectorValue == "" {
		return nil, errors.New("Selector value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests("objects", requestType, selectorKey, selectorValue)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

func GetRequestsPublicIdsBySelector(requestType, selectorKey, selectorValue string) ([]string, error) {

	if requestType == "" {
		return nil, errors.New("Request type value cannot be an empty string")
	}

	if selectorKey == "" {
		return nil, errors.New("Selector key value cannot be an empty string")
	}

	if selectorValue == "" {
		return nil, errors.New("Selector value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests("publicIds", requestType, selectorKey, selectorValue)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

// queryType:
// requestIds
// objects
func GetRequestsByRecipient(queryType, requestType, recipientPublicId string) ([]string, error) {

	if queryType == "" {
		return nil, errors.New("Query type value cannot be an empty string")
	}

	if requestType == "" {
		return nil, errors.New("Request type value cannot be an empty string")
	}

	if recipientPublicId == "" {
		return nil, errors.New("Recipient Id value cannot be an empty string")
	}

	selectorKey := "recipientPublicId"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests(queryType, requestType, selectorKey, recipientPublicId)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

func GetRequestsByRequester(queryType, requestType, requesterPublicId string) ([]string, error) {

	if queryType == "" {
		return nil, errors.New("Query type value cannot be an empty string")
	}

	if requestType == "" {
		return nil, errors.New("Request type value cannot be an empty string")
	}

	if requesterPublicId == "" {
		return nil, errors.New("Requester Id value cannot be an empty string")
	}

	selectorKey := "requesterPublicId"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests(queryType, requestType, selectorKey, requesterPublicId)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

func GetRequestsByDocumentName(queryType, documentName string) ([]string, error) {

	if queryType == "" {
		return nil, errors.New("Query type value cannot be an empty string")
	}

	if documentName == "" {
		return nil, errors.New("Document name value cannot be an empty string")
	}

	requestType := "documentData"
	selectorKey := "documentName"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests(queryType, requestType, selectorKey, documentName)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

func GetRequestsByStatus(queryType, requestType, status string) ([]string, error) {

	if queryType == "" {
		return nil, errors.New("Query type value cannot be an empty string")
	}

	if requestType == "" {
		return nil, errors.New("Request type value cannot be an empty string")
	}

	if status == "" {
		return nil, errors.New("Status value cannot be an empty string")
	}

	selectorKey := "status"

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestsData, err := persAccntsChannelClient.QueryRequests(queryType, requestType, selectorKey, status)

	if err != nil {
		return nil, err
	}

	return []string{string(requestsData)}, nil
}

func GetRequestsByRecipientAndType(recipientId string) ([]string, error) {

	// do we need this?

	return nil, nil
}

// idTypes:
// requestId
// publicId
func GetRequestObject(idType, id string) (string, error) {

	if idType == "" {
		return "", errors.New("Id type value cannot be an empty string")
	}

	if id == "" {
		return "", errors.New("Request Id value cannot be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestData, err := persAccntsChannelClient.QueryRequestData(idType, id)

	if err != nil {
		return "", err
	}

	return string(requestData), nil
}

// Types:
// accountData
// documentData
// general
func GetRequestPublicId(id, requestType string) (string, error) {

	if id == "" {
		return "", errors.New("Request Id value cannot be an empty string")
	}

	if requestType == "" {
		return "", errors.New("Request type value canont be an empty string")
	}

	persAccntsChannelClient := persaccntschannel.CerberusClient{}
	requestData, err := persAccntsChannelClient.QueryRequestData("requestId", id)

	if err != nil {
		return "", err
	}

	switch requestType {
	case "accountData":

		request := &accountDataRequest{}

		err = json.Unmarshal([]byte(requestData), request)

		if err != nil {
			return "", err
		}

		if request.RequestType == "accountData" {
			return request.PublicId, nil
		} else {
			return "", errors.New("Request is not of type \"accountData\"")
		}

	case "documentData":
		request := &documentDataRequest{}

		err = json.Unmarshal([]byte(requestData), request)

		if err != nil {
			return "", err
		}

		if request.RequestType == "documentData" {
			return request.PublicId, nil
		} else {
			return "", errors.New("Request is not of type \"documentData\"")
		}

	default:
		return "", errors.New("Unknown request type")
	}
}
