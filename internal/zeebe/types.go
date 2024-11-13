package zeebe

import (
	"context"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

// TODO: DEFINE IN INTERFACE HERE
type ZeebeCamunda interface {
	DeployProcessDefinition(resourceName string, formResources []string) ([]*pb.ProcessMetadata, []BPMNProcess, error)
	GenerateCRUDHandlers(processMetadata *pb.ProcessMetadata) error
	GenerateCRUDUserTaskServiceTaskHandler(bpmnProcess *[]BPMNProcess) error
	GenerateCRUDFromPayloadHandlers(processName, resourceName string, version int32, processDefinitionKey int64) error
	StartWorkflow(ctx context.Context, processDefinitionKey int64, variables map[string]interface{}) (*pb.CreateProcessInstanceResponse, error)
	CancelWorkflow(context.Context, int64) error
	StartWorker(jobType, nameWorker string, handler worker.JobHandler) (worker.JobWorker, error)
	Close() error
}

type Client struct {
	client zbc.Client
}

type StartWorkFlowResponse struct {
	ProcessDefinitionKey int64  `json:"processDefinitionKey"`
	BpmnProcessId        string `json:"bpmnProcessId"`
	Version              int32  `json:"version"`
	ProcessInstanceKey   int64  `json:"processInstanceKey"`
	TenantId             string `json:"tenantId"`
}

type FormDefinition struct {
	FormID string `xml:"formId,attr"`
}

type Propertie struct {
	Property string `xml:"property,attr"`
}

type TaskDefinition struct {
	Type string `xml:"type,attr"`
}

type AssignmentDefinition struct {
	Assignee        string `xml:"assignee,attr"`
	CandidateGroups string `xml:"candidateGroups,attr"`
	CandidateUsers  string `xml:"candidateUsers,attr"`
}

type TaskSchedule struct {
	DueDate string `xml:"dueDate,attr"`
}

type ExtensionElement struct {
	XMLName               xml.Name               `xml:"extensionElements"`
	FormDefinitions       []FormDefinition       `xml:"formDefinition"`
	AssignmentDefinitions []AssignmentDefinition `xml:"assignmentDefinition"`
	TaskSchedules         []TaskSchedule         `xml:"taskSchedule"`
	TaskDefinitions       []TaskDefinition       `xml:"taskDefinition"`
	Properties            []Propertie            `xml:"properties"`
}

type UserTask struct {
	XMLName           xml.Name           `xml:"userTask"`
	ID                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	ExtensionElements []ExtensionElement `xml:"extensionElements"`
}

type ServiceTask struct {
	XMLName           xml.Name           `xml:"serviceTask"`
	ID                string             `xml:"id,attr"`
	Name              string             `xml:"name,attr"`
	ExtensionElements []ExtensionElement `xml:"extensionElements"`
}

type BPMNProcess struct {
	XMLName     xml.Name      `xml:"process"`
	UserTasks   []UserTask    `xml:"userTask"`
	ServiceTask []ServiceTask `xml:"serviceTask"`
}

type BPMNDocument struct {
	XMLName   xml.Name      `xml:"definitions"`
	Processes []BPMNProcess `xml:"process"`
}

type TokenManager struct {
	clientID     string
	clientSecret string
	authURL      string
	authToken    string
	refreshToken string
	expiry       time.Time
}

type ZeebeClientRest struct {
	httpClient   *http.Client
	tokenManager *TokenManager
	zeebeAddr    string
}

// form id
type FormComponent struct {
	Label       string `json:"label"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Validate    struct {
		Required bool `json:"required"`
	} `json:"validate"`
}

type Form struct {
	Components []FormComponent `json:"components"`
}
