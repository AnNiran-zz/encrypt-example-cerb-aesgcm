package persaccntschannel

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (persAccntsChannelClient *CerberusClient) CreateAccountDataRequest(newRequest, requestData []byte) ([]string, []byte, error) {

	//channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return nil, nil, err
	}
	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)

	if err != nil {
		return nil, nil, err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "createRequest",
		Args:        [][]byte{[]byte("accountData"), newRequest, requestData},
	}

	//response, err := persAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Execute(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return nil, nil, err
	}

	if response.ChaincodeStatus == 200 {
		fmt.Println("Request created successfully.")
		fmt.Println("Transaction ID is: " + response.TransactionID)
	}

	return []string{"200", string(response.TransactionID)}, response.Payload, nil
}

func (persAccntsChannelClient *CerberusClient) CreateDocumentDataRequest(newRequest, requestData []byte) ([]string, []byte, error) {

	//channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return nil, nil, err
	}
	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)

	if err != nil {
		return nil, nil, err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "createRequest",
		Args:        [][]byte{[]byte("documentData"), newRequest, requestData},
	}

	//response, err := persAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Execute(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return nil, nil, err
	}

	if response.ChaincodeStatus == 200 {
		fmt.Println("Request created successfully.")
		fmt.Println("Transaction ID is: " + response.TransactionID)
	}

	return []string{"200", string(response.TransactionID)}, response.Payload, nil
}

func (persAccntsChannelClient *CerberusClient) AcceptRequest(requestType, requestPublicId, recipientPublicId string, acceptedData []byte) ([]string, []byte, error) {

	// channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return nil, nil, err
	}
	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)

	if err != nil {
		return nil, nil, err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "acceptRequest",
		Args:        [][]byte{[]byte(requestType), []byte(requestPublicId), []byte(recipientPublicId), acceptedData},
	}

	//response, err := persAccntsChannelClient.channelClient.Execute(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Execute(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return nil, nil, err
	}

	if response.ChaincodeStatus == 200 {
		fmt.Println("Request accepted successfully.")
		fmt.Println("Transaction ID is: " + response.TransactionID)
	}

	return []string{"200", string(response.TransactionID)}, response.Payload, nil
}

func (persAccntsChannelClient *CerberusClient) RejectRequest(requestType, requestPublicId, recipientPublicId string) ([]string, []byte, error) {

	// channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return nil, nil, err
	}

	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "rejectRequest",
		Args:        [][]byte{[]byte(requestType), []byte(requestPublicId), []byte(recipientPublicId)},
	}

	//response, err := persAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Execute(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return nil, nil, err
	}

	if response.ChaincodeStatus == 200 {
		fmt.Println("Request accepted successfully.")
		fmt.Println("Transaction ID is: " + response.TransactionID)
	}

	return []string{"200", string(response.TransactionID)}, response.Payload, nil
}

func (persAccntsChannelClient *CerberusClient) UpdateRequest(requestType, requestPublicId, requesterPublicId, recipientId string, updatedData []byte) ([]string, []byte, error) {

	// channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return nil, nil, err
	}
	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)
	if err != nil {
		return nil, nil, err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "updateRequest",
		Args:        [][]byte{[]byte(requestType), []byte(requestPublicId), []byte(requesterPublicId), []byte(recipientId), []byte(updatedData)},
	}

	//response, err := instAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Execute(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return nil, nil, err
	}

	if response.ChaincodeStatus == 200 {
		fmt.Println("Request updated successfully.")
		fmt.Println("Transaction ID is: " + response.TransactionID)
	}

	return []string{"200", string(response.TransactionID)}, response.Payload, nil
}
