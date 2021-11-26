package helpers

import (
	"os"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
)

func NewCadenceClient(workflowClient workflowserviceclient.Interface) client.Client {
	domainName := os.Getenv("CADENCE_DOMAIN_NAME")
	return client.NewClient(workflowClient, domainName, &client.Options{})
}
