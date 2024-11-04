package zeebe

import (
	"context"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

// TODO: DEFINE IN INTERFACE HERE
type ZeebeCamunda interface {
	DeployProcessDefinition(resourceName string) (*pb.ProcessMetadata, error)
	GenerateCRUDHandlers(processMetadata *pb.ProcessMetadata) error
	GenerateCRUDFromPayloadHandlers(processName, resourceName string, version int32, processDefinitionKey int64) error
	StartWorkflow(ctx context.Context, workflowName string, variables map[string]interface{}) (string, error)
	StartWorker(jobType, nameWorker string, handler worker.JobHandler) (worker.JobWorker, error)
	Close() error
}

type Client struct {
	client zbc.Client
}
