package persaccntschannel

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (persAccntsChannelClient *CerberusClient) QueryRecords(selectorKey, selectorValue string) (string, error) {

	// channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	if err != nil {
		return "", err
	}
	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)

	if err != nil {
		return "", err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "queryRecords",
		Args:        [][]byte{[]byte(selectorKey), []byte(selectorValue)},
	}

	//response, err := persAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Query(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return "", err
	}

	if len(response.Payload) < 5 { // small random number of bytes
		fmt.Println("No records with " + selectorKey + ":" + selectorValue + " exist.")
		return "", nil
	}

	return string(response.Payload), nil
}

func (persAccntsChannelClient *CerberusClient) QueryAccountData(queryType, publicId string) (string, error) {

	// channel instance -> create
	err := persAccntsChannelClient.setupPersonAccountsChannelClient()

	defer sdkInstance.Close()

	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// request -> prepare
	request := channel.Request{
		ChaincodeID: PersonAccountsChannelChainCode,
		Fcn:         "queryAccountData",
		Args:        [][]byte{[]byte(queryType), []byte(publicId)},
	}

	//response, err := persAccntsChannelClient.channelClient.Query(request)
	// or:
	response, err := persAccntsChannelClient.channelClient.Query(request, channel.WithTargetEndpoints(AnchorPrSipher))

	if err != nil {
		return "", err
	}

	if len(response.Payload) < 5 { // small random number of bytes
		fmt.Println("No records with id: " + publicId + " exist.")
		return "", nil
	}

	return string(response.Payload), nil
}
