package main

import (
	"errors"
	"time"

	"github.com/damarteplok/social/internal/auth"
	"github.com/damarteplok/social/internal/mailer"
	"github.com/damarteplok/social/internal/minioupload"
	"github.com/damarteplok/social/internal/ratelimiter"
	"github.com/damarteplok/social/internal/store"
	"github.com/damarteplok/social/internal/store/cache"
	"github.com/damarteplok/social/internal/zeebe"
	"go.uber.org/zap"
)

// api types
type application struct {
	config          config
	store           store.Storage
	cacheStorage    cache.Storage
	logger          *zap.SugaredLogger
	mailer          mailer.Client
	authenticator   auth.Authenticator
	rateLimiter     ratelimiter.Limiter
	zeebeClient     zeebe.ZeebeCamunda
	zeebeClientRest zeebe.ZeebeClientRest
	minioClient     minioupload.MinioApi
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	minio       minioConfig
	camunda     camundaConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
	camundaRest camundaRestConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	sendgrid  sendGridConfig
	exp       time.Duration
	fromEmail string
}

type minioConfig struct {
	addr      string
	port      int
	ssl       bool
	accessKey string
	secretKey string
	bucket    string
	expires   time.Duration
	enabled   bool
}

type camundaConfig struct {
	zeebeAddr          string
	zeebeClientId      string
	zeebeClientSecret  string
	zeebeAuthServerUrl string
}

type camundaRestConfig struct {
	zeebeRestAddress       string
	zeebeGrpcAddress       string
	camundaTasklistBaseUrl string
	camundaOperateBaseUrl  string
	camundaOptimizeBaseUrl string
}

type sendGridConfig struct {
	apiKey string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

// health types
type HealthResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// posts types
type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

type DataStorePostWrapper struct {
	Data store.Post `json:"data"`
}

// users types
type DataStoreUserWrapper struct {
	Data store.User `json:"data"`
}

// camunda types
type DeployBpmnPayload struct {
	ResourceName  string   `json:"resource_name" validate:"required"`
	FormResources []string `json:"form_resources" validate:"omitempty,min=0,dive"`
}
type CrudPayload struct {
	ProcessName          string `json:"process_name" validate:"required,max=255"`
	ResourceName         string `json:"resource_name" validate:"required,max=255"`
	Version              int32  `json:"version" validate:"required"`
	ProcessDefinitionKey int64  `json:"process_definition_key" validate:"required"`
}

type StartInstruction struct {
	ElementID *string `json:"elementId,omitempty"`
}

type CreateProcessInstancePayload struct {
	ProcessDefinitionKey int64               `json:"processDefinitionKey" validate:"required"`
	Variables            *map[string]string  `json:"variables,omitempty"`
	TenantID             *string             `json:"tenantId,omitempty"`
	OperationReference   *int64              `json:"operationReference,omitempty"`
	StartInstructions    *[]StartInstruction `json:"startInstructions,omitempty"`
	AwaitCompletion      *bool               `json:"awaitCompletion,omitempty"`
	FetchVariables       *[]string           `json:"fetchVariables,omitempty"`
	RequestTimeout       *int64              `json:"requestTimeout,omitempty"`
}

type CreateProcessInstancesResponse struct {
	ProcessDefinitionKey     int64             `json:"processDefinitionKey"`
	ProcessDefinitionId      string            `json:"processDefinitionId"`
	ProcessDefinitionVersion int32             `json:"processDefinitionVersion"`
	ProcessInstanceKey       int64             `json:"processInstanceKey"`
	TenantId                 string            `json:"tenantId"`
	Variables                map[string]string `json:"variables"`
}

type SortSearchTasklist struct {
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}

type SearchTaskListPayload struct {
	State                string               `json:"state,omitempty"`
	Assigned             bool                 `json:"assgined,omitempty"`
	Assignee             string               `json:"assignee,omitempty"`
	Assignees            []string             `json:"assignees,omitempty"`
	TaskDefinitionId     string               `json:"taskDefinitionId,omitempty"`
	CandidateGroup       string               `json:"candidateGroup,omitempty"`
	CandidateGroups      []string             `json:"candidateGroups,omitempty"`
	CandidateUser        string               `json:"candidateUser,omitempty"`
	CandidateUsers       []string             `json:"candidateUsers,omitempty"`
	ProcessDefinitionKey string               `json:"processDefinitionKey,omitempty"`
	ProcessInstanceKey   string               `json:"processInstanceKey,omitempty"`
	PageSize             int32                `json:"pageSize,omitempty"`
	Sort                 []SortSearchTasklist `json:"sort,omitempty"`
	SearchAfter          []string             `json:"searchAfter,omitempty"`
	SearchAfterOrEqual   []string             `json:"searchAfterOrEqual,omitempty"`
	SearchBefore         []string             `json:"searchBefore,omitempty"`
	SearchBeforeOrEqual  []string             `json:"searchBeforeOrEqual,omitempty"`
}

type QueryUserTaskPayload struct {
	Sort   []Sort `json:"sort"`
	Filter Filter `json:"filter,omitempty"`
	Page   Page   `json:"page,omitempty"`
}

type Page struct {
	From         int64                `json:"from"`
	Limit        int64                `json:"limit"`
	SearchAfter  []SearchAfterPayload `json:"searchAfter"`
	SearchBefore []SearchAfterPayload `json:"searchBefore"`
}

type SearchAfterPayload struct {
	Object []interface{} `json:"object,omitempty"`
}

type Variable struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Filter struct {
	Key                  int64      `json:"key,omitempty"`
	State                string     `json:"state,omitempty"`
	Assignee             string     `json:"assignee,omitempty"`
	ElementId            string     `json:"elementId,omitempty"`
	CandidateGroup       string     `json:"candidateGroup,omitempty"`
	CandidateUser        string     `json:"candidateUser,omitempty"`
	ProcessDefinitionKey int64      `json:"processDefinitionKey,omitempty"`
	ProcessInstanceKey   int64      `json:"processInstanceKey,omitempty"`
	TenantIds            string     `json:"tenantIds,omitempty"`
	ProcessDefinitionId  string     `json:"processDefinitionId,omitempty"`
	Variables            []Variable `json:"variables,omitempty"`
}

type Sort struct {
	Field string `json:"field" validate:"required"`
	Order string `json:"order,omitempty"`
}

func (p *SearchTaskListPayload) IsValidState() error {
	switch p.State {
	case StateCreated, StateCompleted, StateCanceled, StateFailed:
		return nil
	default:
		return errors.New("invalid state: must be one of CREATED, COMPLETED, CANCELED, FAILED")
	}
}

// authenticated types
type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

type FlowNodeQueryParams struct {
	Size                     string
	Order                    string
	Sort                     string
	SearchAfter              string
	SearchBefore             string
	Type                     string
	State                    string
	BpmnProcessId            string
	ProcessDefinitionKey     string
	ParentProcessInstanceKey string
	StartDate                string
	EndDate                  string
}

type TaskListQueryParams struct {
	Size               int32
	State              string
	TaskDefinitionId   string
	ProcessInstanceKey string
	Sort               string
	Order              string
	SearchAfter        string
	SearchBefore       string
}

type OperateCoreStats struct {
	Running       int `json:"running"`
	Active        int `json:"active"`
	WithIncidents int `json:"withIncidents"`
}

type Process struct {
	ProcessID                         string  `json:"processId"`
	Version                           int     `json:"version"`
	Name                              string  `json:"name"`
	BpmnProcessID                     string  `json:"bpmnProcessId"`
	TenantID                          string  `json:"tenantId"`
	ErrorMessage                      *string `json:"errorMessage"`
	InstancesWithActiveIncidentsCount int     `json:"instancesWithActiveIncidentsCount"`
	ActiveInstancesCount              int     `json:"activeInstancesCount"`
}

type OperateProcessStats struct {
	BpmnProcessID                     string    `json:"bpmnProcessId"`
	TenantID                          string    `json:"tenantId"`
	ProcessName                       *string   `json:"processName"`
	InstancesWithActiveIncidentsCount int       `json:"instancesWithActiveIncidentsCount"`
	ActiveInstancesCount              int       `json:"activeInstancesCount"`
	Processes                         []Process `json:"processes"`
}
