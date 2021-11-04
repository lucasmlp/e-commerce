package workflows

import (
	"fmt"

	"go.uber.org/cadence/workflow"
)

func RunOrder(ctx workflow.Context) error {
	fmt.Println("executing order workflow")
	return nil
}
