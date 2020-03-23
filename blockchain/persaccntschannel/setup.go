package persaccntschannel

import (
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

const (
	PersonAccountsChannelID        = "persaccntschannel"
	PersonAccountsChannelChainCode = "persaccntschannelcc"

	SipherOrg   = "Sipher"
	WhiteBoxOrg = "WhiteBox"

	SipherAdmin   = "Admin"
	SipherUser    = "User1"
	WhiteBoxAdmin = "Admin"

	AnchorPrSipher   = "anchorpr.sipher.cerberus.dev"
	AnchorPrWhiteBox = "anchorpr.whitebox.cerberus.dev"
)

type CerberusClient struct {
	channelCtx    context.ChannelProvider
	channelClient *channel.Client
	registration  fab.Registration
	initialized   bool
	event         *event.Client
}

var sdkInstance *fabsdk.FabricSDK
var err error

func (persAccntsChannelClient *CerberusClient) setupPersonAccountsChannelClient() error {

	// check sdk instance
	if persAccntsChannelClient.initialized {
		fmt.Println("sdk is already initialized.")
		return errors.New("sdk is already initialized.") // fix check
	}

	// config file -> get
	configFile := os.Getenv("GOPATH") + "/src/cerberus/hl/config.yaml"

	// sdk instance -> create
	sdkInstance, err = fabsdk.New(config.FromFile(configFile))
	if err != nil {
		fmt.Println(err)
		return err
	}

	persAccntsChannelClient.channelCtx = sdkInstance.ChannelContext(PersonAccountsChannelID, fabsdk.WithUser(SipherUser), fabsdk.WithOrg(SipherOrg))

	// register event
	persAccntsChannelClient.event, err = event.New(persAccntsChannelClient.channelCtx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var notifier <-chan *fab.CCEvent
	persAccntsChannelClient.registration, notifier, err = persAccntsChannelClient.event.RegisterChaincodeEvent(PersonAccountsChannelChainCode, "event123")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer persAccntsChannelClient.event.Unregister(persAccntsChannelClient.registration) // optimize

	select {
	case chaincodeEvent := <-notifier:
		fmt.Printf("received chaincode event %v\n", chaincodeEvent)
	case <-time.After(time.Second * 5):
		fmt.Println("timeout while waiting for chaincode event")
	}

	// instantiate channel
	persAccntsChannelClient.channelClient, err = channel.New(persAccntsChannelClient.channelCtx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
