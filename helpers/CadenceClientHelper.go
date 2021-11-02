package helpers

import (
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
)

const (
	domainName = "poc"
)

func NewCadenceClient(workflowClient workflowserviceclient.Interface) client.Client {
	return client.NewClient(workflowClient, domainName, &client.Options{})
}
