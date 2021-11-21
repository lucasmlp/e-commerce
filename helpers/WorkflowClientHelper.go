package helpers

import (
	"errors"
	"os"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
)

func NewWorkflowClient(serviceNameCadenceClient string, serviceNameCadenceFrontend string) (workflowserviceclient.Interface, error) {

	cadenceHostAddressAndPort := os.Getenv("CADENCE_HOST_ADDRESS_AND_PORT")

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(serviceNameCadenceClient))
	if err != nil {
		return nil, err
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: serviceNameCadenceClient,
		Outbounds: yarpc.Outbounds{
			serviceNameCadenceFrontend: {Unary: ch.NewSingleOutbound(cadenceHostAddressAndPort)},
		},
	})

	if dispatcher == nil {
		return nil, errors.New("failed to create dispatcher")
	}

	if err := dispatcher.Start(); err != nil {
		panic(err)
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(serviceNameCadenceFrontend)), nil
}
