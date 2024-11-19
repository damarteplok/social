package zeebe

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

//go:embed "resources"
var res embed.FS

func NewZeebeClient(clientID, clientSecret, authURL, zeebeAddr string) (ZeebeCamunda, error) {
	credentials, err := zbc.NewOAuthCredentialsProvider(&zbc.OAuthProviderConfig{
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		AuthorizationServerURL: authURL,
		Audience:               "zeebe-api",
	})
	if err != nil {
		return nil, err
	}
	config := zbc.ClientConfig{
		UsePlaintextConnection: true,
		GatewayAddress:         "localhost:26500",
		CredentialsProvider:    credentials,
	}
	client, err := zbc.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zeebe client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// UpdateProcessInstance updates the variables of a process instance.
func (c *Client) UpdateProcessInstance(ctx context.Context, processInstanceKey int64, variables map[string]interface{}) error {
	request, err := c.client.NewSetVariablesCommand().ElementInstanceKey(processInstanceKey).VariablesFromMap(variables)
	if err != nil {
		return fmt.Errorf("failed to update process instance: %w", err)
	}

	_, err = request.Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to update variables from process instance: %w", err)
	}

	return nil
}

// StartWorkflow starts a new workflow instance.
func (c *Client) StartWorkflow(ctx context.Context, processDefinitionKey int64, variables map[string]interface{}) (*pb.CreateProcessInstanceResponse, error) {
	request, err := c.client.NewCreateInstanceCommand().ProcessDefinitionKey(processDefinitionKey).VariablesFromMap(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow instance request: %w", err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start workflow: %w", err)
	}

	return result, nil
}

// CancelWorkflow cancels a workflow instance.
func (c *Client) CancelWorkflow(ctx context.Context, processInstanceKey int64) error {
	_, err := c.client.NewCancelInstanceCommand().ProcessInstanceKey(processInstanceKey).Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to cancel workflow: %w", err)
	}
	return nil
}

// DeployProcessDefinition deploys a process definition.
func (z *Client) DeployProcessDefinition(resourceName string, formResources []string) ([]*pb.ProcessMetadata, []BPMNProcess, error) {
	definition, err := MustReadFile(resourceName)
	if err != nil {
		return nil, nil, errors.New("failed to read file")
	}

	bpmnProcess, err := unMarshalBpmn(definition)
	if err != nil {
		return nil, nil, errors.New("failed to unmarshal xml")
	}

	command := z.client.NewDeployResourceCommand().AddResource(definition, resourceName)

	for _, formResource := range formResources {
		formDefinition, err := MustReadFile(formResource)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read form file %s: %w", formResource, err)
		}
		command = command.AddResource(formDefinition, formResource)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	resource, err := command.Send(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(resource.GetDeployments()) < 1 {
		return nil, nil, errors.New("failed to deploy model; nothing was deployed")
	}

	var processes []*pb.ProcessMetadata
	for _, deployment := range resource.GetDeployments() {
		process := deployment.GetProcess()
		if process != nil {
			processes = append(processes, process)
		}
	}

	return processes, bpmnProcess, nil
}

// StartWorker starts a worker for a given job type.
func (z *Client) StartWorker(jobType, nameWorker string, handler worker.JobHandler) (worker.JobWorker, error) {
	w := z.client.NewJobWorker().
		JobType(jobType).
		Handler(handler).
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name(nameWorker).
		Open()

	return w, nil
}

// GenerateCRUDFromPayloadHandlers generates CRUD handlers from payload.
func (c *Client) GenerateCRUDFromPayloadHandlers(processName, resourceName string, version int32, processDefinitionKey int64) error {
	processNameTitle := toCamelCase(processName)
	tableName := processName
	if err := generateCrudProcess(processNameTitle, resourceName, tableName, version, processDefinitionKey); err != nil {
		return err
	}

	return nil
}

// GenerateCRUDHandlers generates CRUD handlers from process metadata.
func (c *Client) GenerateCRUDHandlers(processMetadata *pb.ProcessMetadata) error {
	processName := toCamelCase(processMetadata.GetBpmnProcessId())
	tableName := processMetadata.GetBpmnProcessId()
	version := processMetadata.GetVersion()
	resourceName := processMetadata.GetResourceName()
	processDefinitionKey := processMetadata.GetProcessDefinitionKey()

	if err := generateCrudProcess(processName, resourceName, tableName, version, processDefinitionKey); err != nil {
		return err
	}

	return nil
}

// GenerateCRUDUserTaskServiceTaskHandler generates CRUD handlers for user task and service task.
func (c *Client) GenerateCRUDUserTaskServiceTaskHandler(bpmnProcess *[]BPMNProcess) error {
	for _, process := range *bpmnProcess {
		for _, serviceTask := range process.ServiceTask {
			if err := generateCrudServiceTask(serviceTask); err != nil {
				return err
			}
		}
		for _, userTask := range process.UserTasks {
			if err := generateCrudUserTask(userTask); err != nil {
				return err
			}
		}
	}
	return nil
}

// generate crud code for service task
func generateCrudServiceTask(serviceTask ServiceTask) error {
	idServiceTask := serviceTask.ID
	nameServiceTask := serviceTask.Name
	var serviceTaskName string
	var nameFile string
	var taskDefinitionType string

	for _, extensionElement := range serviceTask.ExtensionElements {
		for _, taskDefinition := range extensionElement.TaskDefinitions {
			taskDefinitionType = taskDefinition.Type
			serviceTaskName = toCamelCase(taskDefinition.Type)
			nameFile = strings.ToLower(taskDefinition.Type)
		}
	}
	filePathStore := fmt.Sprintf("./internal/service/%s_service_task.go", nameFile)

	serviceTaskCode := fmt.Sprintf(`package store
const (
	%sID = "%s"
	%sName = "%s"
	%sType = "%s"
)
// TODO: DO SOMETHING IN SERVICE TASK
`, serviceTaskName, idServiceTask, serviceTaskName, nameServiceTask, serviceTaskName, taskDefinitionType)
	err := os.WriteFile(filePathStore, []byte(serviceTaskCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write service task file: %w", err)
	}

	return nil
}

// generate code crud for user task
func generateCrudUserTask(userTask UserTask) error {
	idServiceTask := userTask.ID
	nameServiceTask := userTask.Name
	userTaskName := toCamelCase(userTask.Name)
	nameFile := strings.ToLower(userTaskName)
	var formID string
	var assignee string
	var candidateGroup string
	var candidateUser string
	var dueDate string
	for _, extensionElement := range userTask.ExtensionElements {
		for _, formDefinition := range extensionElement.FormDefinitions {
			formID = formDefinition.FormID
		}
		for _, assigneeDefinition := range extensionElement.AssignmentDefinitions {
			assignee = assigneeDefinition.Assignee
			candidateGroup = assigneeDefinition.CandidateGroups
			candidateUser = assigneeDefinition.CandidateUsers
		}
		for _, taskSchedule := range extensionElement.TaskSchedules {
			dueDate = taskSchedule.DueDate
		}
	}

	moduleName, errModule := getModuleName()
	if errModule != nil {
		return errModule
	}

	filePathHandler := fmt.Sprintf("./cmd/api/%s_user_task.go", nameFile)
	filePathScripts := fmt.Sprintf("./scripts/%s_user_task.sql", nameFile)
	scriptCode := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	task_id VARCHAR(256),
	form_id VARCHAR(256),
	properties JSONB,
	created_by BIGINT NOT NULL,
	updated_by BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE
);

DROP TABLE IF EXISTS %s;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_%s_properties ON %s USING gin (properties);

	`, nameFile, nameFile, nameFile, nameFile)

	err := os.WriteFile(filePathScripts, []byte(scriptCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	filePathStore := fmt.Sprintf("./internal/store/%s_user_task.go", nameFile)
	modelCode := fmt.Sprintf(`package store
import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const (
	%sID = "%s"
	%sName = "%s"
	%sFormID = "%s"
	%sAssignee = "%s"
	%sCandidateGroup = "%s"
	%sCandidateUser = "%s"
	%sSchedule = %s%s%s

)

type %s struct {
    ID         int64    `+"`json:\"id\"`"+`
	Name       string   `+"`json:\"name\"`"+`
	TaskId     string   `+"`json:\"task_id\"`"+`
	FormId     string   `+"`json:\"form_id\"`"+`
	Properties []string `+"`json:\"properties\"`"+`
	CreatedBy  int64    `+"`json:\"created_by\"`"+`
	UpdatedBy  *int64   `+"`json:\"updated_by\"`"+`
	CreatedAt  string   `+"`json:\"created_at\"`"+`
	UpdatedAt  string   `+"`json:\"updated_at\"`"+`
	DeletedAt  *string  `+"`json:\"deleted_at\"`"+`
}

type %sStore struct {
	db *sql.DB
}

func (s *%sStore) Create(ctx context.Context, model *%s) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *%sStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *%sStore) Update(ctx context.Context, model *%s) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *%sStore) create(ctx context.Context, tx *sql.Tx, model *%s) error {
	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}

	query := %s
		INSERT INTO %s (name, form_id, properties, created_by, task_id)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING 
		 	id, name, form_id, task_id, properties, created_by, updated_by,
			created_at, updated_at
		%s
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var propertiesData []byte
	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.FormId,
		propertiesJSON,
		model.CreatedBy,
		model.TaskId,
	).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
		&model.TaskId,
		&propertiesData,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *%sStore) GetByID(ctx context.Context, id int64) (*%s, error) {
	query := %s
		SELECT id, name, form_id, task_id, properties, created_by, 
		updated_by, created_at, updated_at
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model %s
	var propertiesData []byte
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
		&model.TaskId,
		&propertiesData,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
		
	if len(propertiesData) > 0 {
		if err := json.Unmarshal(propertiesData, &model.Properties); err != nil {
			return nil, err
		}
	} else {
		model.Properties = []string{}
	}

	return &model, nil
}

func (s *%sStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := %sUPDATE %s SET deleted_at = NOW() WHERE id = $1;%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *%sStore) update(ctx context.Context, tx *sql.Tx, model *%s) error {

	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}
	query := %s
		UPDATE %s
		SET name = $1, form_id = $2, properties = $3, updated_by = $4, task_id = $5,  updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, form_id, task_id, properties, created_by, updated_by, created_at, updated_at;
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var propertiesData []byte
	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.FormId,
		propertiesJSON,
		model.UpdatedBy,
		model.TaskId,
		model.ID,
	).Scan(&model.ID, 
		&model.Name, 
		&model.FormId, 
		&model.TaskId,
		propertiesData, 
		&model.CreatedBy, 
		&model.UpdatedBy, 
		&model.CreatedAt, 
		&model.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

`,
		userTaskName,
		idServiceTask,
		userTaskName,
		nameServiceTask,
		userTaskName,
		formID,
		userTaskName,
		assignee,
		userTaskName,
		candidateGroup,
		userTaskName,
		candidateUser,
		userTaskName,
		"`",
		dueDate,
		"`",
		userTaskName,
		userTaskName,

		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,

		userTaskName,
		userTaskName,
		"`",
		nameFile,
		"`",
		userTaskName,
		userTaskName,
		"`",
		nameFile,
		"`",
		userTaskName,

		userTaskName,
		"`",
		nameFile,
		"`",
		userTaskName,
		userTaskName,
		"`",
		nameFile,
		"`",
	)

	err = os.WriteFile(filePathStore, []byte(modelCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	// if formID is not empty then generate code form
	structCode := ""
	if formID != "" {
		formFile := formID + ".form"
		filePathForm := fmt.Sprintf("./internal/zeebe/resources/%s", formFile)

		form, errForm := readFormFile(filePathForm)
		if errForm != nil {
			return errForm
		}

		structCode = generateStructCode(form, userTaskName)
	}

	// create handler usertask
	handlerUserTaskCode := fmt.Sprintf(`package main
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"%s/internal/store"
)

%s

// GetUserTaskActive %s godoc
//
//	@Summary		GetUserTaskActive %s
//	@Description	GetUserTaskActive %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			size			query		string	false	"Size 50"
//	@Param			order			query		string	false	"Order DESC ASC"
//	@Param			sort			query		string	false	"Sort creationTime"
//
// @Param			state			query		string	false	"State CREATED"
//
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200				{string}	string	"%s GetUserTaskActive"
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/search  [get]
func (app *application) getUserTaskActive%sHandler(w http.ResponseWriter, r *http.Request) {
	taskListQueryParams, err := getTaskListQueryParams(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if taskListQueryParams.Sort == "" {
		taskListQueryParams.Sort = "creationTime"
	}

	if taskListQueryParams.State == "" {
		taskListQueryParams.State = "CREATED"
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	jsonTemplate := %s
{
	"taskDefinitionId": "%s",
	"state:": "%s",
    "pageSize": %s, 
    "sort": [{"field": "%s", "order": "%s"}]
    %s
}%s

	var searchAfterStr, searchBeforeStr string

	if taskListQueryParams.SearchAfter != "" {
		searchAfterStr = fmt.Sprintf(%s, "searchAfter": %s%s, taskListQueryParams.SearchAfter)
	}

	if taskListQueryParams.SearchBefore != "" {
		searchBeforeStr = fmt.Sprintf(%s, "searchBefore": %s%s, taskListQueryParams.SearchBefore)
	}

	searchParams := searchAfterStr + searchBeforeStr

	body := []byte(fmt.Sprintf(
		jsonTemplate,
		store.%sID,
		taskListQueryParams.State,
		taskListQueryParams.Size,
		taskListQueryParams.Sort,
		taskListQueryParams.Order,
		searchParams))

	// get history from operate api rest api
	url := fmt.Sprintf("%s%s/search", app.config.camundaRest.camundaTasklistBaseUrl, V1TasklistUrl)
	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	var jsonData interface{}
	if err := json.Unmarshal(resp, &jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	
`,
		moduleName,
		structCode,
		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,
		nameFile,
		userTaskName,

		"`",
		"%s",
		"%s",
		"%d",
		"%s",
		"%s",
		"%s",
		"`",

		"`",
		"%s",
		"`",

		"`",
		"%s",
		"`",

		userTaskName,

		"%s",
		"%s",
	)

	err = os.WriteFile(filePathHandler, []byte(handlerUserTaskCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	filePathEditRoutes := "./cmd/api/api.go"
	generateCodeRoutes := fmt.Sprintf(`
			r.Route("/%s", func(r chi.Router) {
				r.Get("/search", app.getUserTaskActive%sHandler)
			})	
`, nameFile, userTaskName)

	err = insertGeneratedCode(filePathEditRoutes, generateCodeRoutes, "// GENERATE USER TASK ROUTES API")
	if err != nil {
		return err
	}

	filePathEditStorage := "./internal/store/storage.go"

	// edit file storage
	generateCodeStorage := fmt.Sprintf(`
	%s interface {
		Create(context.Context, *%s) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*%s, error)
	}
`,
		userTaskName,
		userTaskName,
		userTaskName,
	)
	generateCodeConstructor := fmt.Sprintf(`
		%s:   &%sStore{db},
`, userTaskName, userTaskName)

	err = insertGeneratedCode(filePathEditStorage, generateCodeStorage, "// GENERATED CODE INTERFACE")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditStorage, generateCodeConstructor, "// GENERATED CODE CONSTRUCTOR")
	if err != nil {
		return err
	}

	// cache
	filePathEditCacheStorage := "./internal/store/cache/storage.go"
	filePathStoreCache := fmt.Sprintf("./internal/store/cache/%s_user_task.go", nameFile)
	modelCacheCode := fmt.Sprintf(`package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"%s/internal/store"
	"github.com/go-redis/redis/v8"
)

type %sStore struct {
	rdb *redis.Client
}

const %sExpTime = time.Hour * 24 * 7
	
func (s *%sStore) Get(ctx context.Context, modelID int64) (*store.%s, error) {
	cacheKey := fmt.Sprintf("%s-%s", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.%s
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *%sStore) Set(ctx context.Context, model *store.%s) error {
	cacheKey := fmt.Sprintf("%s-%s", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, %sExpTime).Err()
}

func (s *%sStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("%s-%s", modelID)
	s.rdb.Del(ctx, cacheKey)
}
`,
		moduleName,
		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,
		userTaskName,
		"%v",
		userTaskName,

		userTaskName,
		userTaskName,
		userTaskName,
		"%v",
		userTaskName,
		userTaskName,
		userTaskName,
		"%v",
	)

	err = os.WriteFile(filePathStoreCache, []byte(modelCacheCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// edit file cache storage
	generateCodeCacheStorage := fmt.Sprintf(`
	%s interface {
		Get(context.Context, int64) (*store.%s, error)
		Set(context.Context, *store.%s) error
		Delete(context.Context, int64)
	}
`,
		userTaskName,
		userTaskName,
		userTaskName,
	)

	generateCodeCacheInterface := fmt.Sprintf(`
		%s: &%sStore{
			rdb: rbd,
		},
`,
		userTaskName,
		userTaskName,
	)

	err = insertGeneratedCode(filePathEditCacheStorage, generateCodeCacheStorage, "// GENERATED CACHE CODE INTERFACE")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditCacheStorage, generateCodeCacheInterface, "// GENERATED CACHE CODE CONSTRUCTOR")
	if err != nil {
		return err
	}

	return nil
}

func generateCrudProcess(processName, resourceName, tableName string, version int32, processDefinitionKey int64) error {
	filePathHandler := fmt.Sprintf("./cmd/api/%s_process.go", tableName)
	filePathStore := fmt.Sprintf("./internal/store/%s_process.go", tableName)
	filePathStoreCache := fmt.Sprintf("./internal/store/cache/%s_process.go", tableName)
	filePathScripts := fmt.Sprintf("./scripts/%s_process.sql", tableName)
	filePathEditCacheStorage := "./internal/store/cache/storage.go"
	filePathEditStorage := "./internal/store/storage.go"
	filePathEditRoutes := "./cmd/api/api.go"

	moduleName, errModule := getModuleName()
	if errModule != nil {
		return errModule
	}

	scriptCode := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id BIGSERIAL PRIMARY KEY,
	process_definition_key BIGINT NOT NULL,
	version INT NOT NULL,
	resource_name VARCHAR(256) NOT NULL,
	process_instance_key BIGINT,
	task_definition_id VARCHAR(256),
	task_state VARCHAR(20) NOT NULL DEFAULT 'CREATED', 
	created_by BIGINT NOT NULL,
	updated_by BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE,
	CONSTRAINT task_state_check CHECK (task_state IN ('CREATED', 'COMPLETED', 'CANCELED', 'FAILED'))
)

DROP TABLE IF EXISTS %s;
	`, tableName, tableName)

	err := os.WriteFile(filePathScripts, []byte(scriptCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	modelCode := fmt.Sprintf(`package store
import (
	"context"
	"database/sql"
	"errors"
)

const (
	%sVersion = %d
	%sProcessDefinitionKey = %d
	%sResourceName = "%s"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type %s struct {
    ID                   int64   `+"`json:\"id\"`"+`
	ProcessDefinitionKey int64   `+"`json:\"process_definition_key\"`"+`
	Version              int32   `+"`json:\"version\"`"+`
	ResourceName         string  `+"`json:\"resource_name\"`"+`
	ProcessInstanceKey   int64   `+"`json:\"process_instance_key\"`"+`
	TaskDefinitionId     *string `+"`json:\"task_definition_id\"`"+`
	TaskState            *string `+"`json:\"task_state\"`"+`
	CreatedBy            int64   `+"`json:\"created_by\"`"+`
	UpdatedBy            *int64  `+"`json:\"updated_by\"`"+`
	CreatedAt            string  `+"`json:\"created_at\"`"+`
	UpdatedAt            string  `+"`json:\"updated_at\"`"+`
	DeletedAt            *string `+"`json:\"deleted_at\"`"+`
}

type %sStore struct {
	db *sql.DB
}

func (s *%sStore) Create(ctx context.Context, model *%s) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *%sStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *%sStore) Update(ctx context.Context, model *%s) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *%sStore) create(ctx context.Context, tx *sql.Tx, model *%s) error {
	//model.Version = %d
	//model.ProcessDefinitionKey = %d
	model.ResourceName = "%s"

	query := %s
		INSERT INTO %s (
			process_definition_key, version, 
			resource_name, process_instance_key, created_by
		) VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5
		) RETURNING 
		 	id, process_definition_key, version, resource_name, process_instance_key, created_by, updated_by,
			created_at, updated_at
		%s
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ProcessInstanceKey,
		model.CreatedBy,
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *%sStore) GetByID(ctx context.Context, id int64) (*%s, error) {
	query := %s
		SELECT id, process_definition_key, version, 
			resource_name, process_instance_key, 
			task_definition_id, task_state,
			created_by, updated_by, created_at, updated_at
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model %s
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.TaskDefinitionId,
		&model.TaskState,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &model, nil
}

func (s *%sStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := %sUPDATE %s SET deleted_at = NOW() WHERE id = $1;%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *%sStore) update(ctx context.Context, tx *sql.Tx, model *%s) error {
	query := %s
		UPDATE %s
		SET process_definition_key = $1, 
			version = $2, 
			resource_name = $3, 
			process_instance_key = $4, 
			updated_by = $5, 
			updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, process_definition_key, 
			version, 
			resource_name, 
			process_instance_key, 
			created_by, 
			updated_by, 
			created_at, 
			updated_at
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ProcessInstanceKey,
		model.UpdatedBy,
		model.ID,
	).Scan(&model.ID, 
		&model.ProcessDefinitionKey, 
		&model.Version, 
		&model.ResourceName, 
		&model.ProcessInstanceKey, 
		&model.CreatedBy, 
		&model.UpdatedBy, 
		&model.CreatedAt, 
		&model.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *%sStore) Search(ctx context.Context, pq PaginatedQuery) (map[string]interface{}, error) {
	sortOrder := "DESC"
	if pq.Sort == "desc" || pq.Sort == "DESC" {
		sortOrder = "DESC"
	}
	if pq.Sort == "asc" || pq.Sort == "ASC" {
		sortOrder = "ASC"
	}

	// Base Query
	query := %s
	 	SELECT p.id, p.process_definition_key, p.version,
            p.resource_name, p.process_instance_key,
            p.task_definition_id, p.task_state,
            p.created_by, p.updated_by, p.created_at, p.updated_at
        FROM pembuatan_media_berita_technology p
        LEFT JOIN users u ON p.created_by = u.id
        WHERE p.deleted_at IS NULL
    %s

	var params []interface{}
	params = append(params, pq.Limit, pq.Offset)

	if pq.Search != "" {
		query += %s
		AND (
				p.process_definition_key::text ILIKE '%s' || $3 || '%s' OR 
				p.resource_name ILIKE '%s' || $3 || '%s' OR
				p.process_instance_key::text ILIKE '%s' || $3 || '%s' OR
				p.task_definition_id::text ILIKE '%s' || $3 || '%s' OR
				p.task_state ILIKE '%s' || $3 || '%s' OR
				u.email ILIKE '%s' || $3 || '%s' OR
				u.username ILIKE '%s' || $3 || '%s'
			)
		%s
		params = append(params, pq.Search)
	}

	if pq.Since != "" && pq.Until != "" {
		query += %s
			AND p.created_at BETWEEN $4 AND $5
        %s
		params = append(params, pq.Since, pq.Until)
	}

	query += %s
        ORDER BY p.created_at %s
        LIMIT $1 OFFSET $2
    %s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var results []*%s
	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item PembuatanMediaBeritaTechnology
		if err := rows.Scan(
			&item.ID,
			&item.ProcessDefinitionKey,
			&item.Version,
			&item.ResourceName,
			&item.ProcessInstanceKey,
			&item.TaskDefinitionId,
			&item.TaskState,
			&item.CreatedBy,
			&item.UpdatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := %s
		SELECT COUNT(*)
        FROM %s p
        LEFT JOIN users u ON p.created_by = u.id
        WHERE p.deleted_at IS NULL
    %s
	var args []interface{}
	if pq.Search != "" {
		countQuery += %s
			%s
		%s
		args = append(args, pq.Search)
	}

	var totalCount int
	err = s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pq.Limit - 1) / pq.Limit

	response := map[string]interface{}{
		"content":      results,
		"totalElement": totalCount,
		"totalPages":   totalPages,
		"limit":        pq.Limit,
		"offset":       pq.Offset,
		"sort":         pq.Sort,
		"search":       pq.Search,
		"since":        pq.Since,
		"until":        pq.Until,
	}

	return response, nil
}
`,
		processName,
		version,
		processName,
		processDefinitionKey,
		processName,
		resourceName,
		processName,
		processName,

		processName,
		processName,
		processName,
		processName,
		processName,

		processName,
		processName,
		version,
		processDefinitionKey,
		resourceName,
		"`",
		tableName,
		"`",
		processName,
		processName,
		"`",
		tableName,
		"`",
		processName,
		processName,
		"`",
		tableName,
		"`",
		processName,
		processName,
		"`",
		tableName,
		"`",

		processName,
		"`", "`", "`",
		"%", "%", "%", "%", "%", "%", "%", "%", "%", "%", "%", "%", "%", "%",
		"`", "`", "`", "`", "` + sortOrder + `", "`",
		processName,
		"`", tableName, "`", "`",
		`
		AND (
                p.process_definition_key::text ILIKE '%' || $1 || '%' OR 
                p.resource_name ILIKE '%' || $1 || '%' OR
                p.process_instance_key::text ILIKE '%' || $1 || '%' OR
                p.task_definition_id::text ILIKE '%' || $1 || '%' OR
                p.task_state ILIKE '%' || $1 || '%' OR
                u.email ILIKE '%' || $1 || '%' OR
                u.username ILIKE '%' || $1 || '%'
            )
		`, "`",
	)

	err = os.WriteFile(filePathStore, []byte(modelCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	handlerCode := fmt.Sprintf(`package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"%s/internal/store"
	"github.com/go-chi/chi/v5"
)

type Create%sPayload struct {
	Variables   		 *map[string]string  %sjson:"variables,omitempty"%s
}
type Update%sPayload struct {
	Variables            *map[string]string  %sjson:"variables,omitempty"%s
}
type DataStore%sWrapper struct {
	Data store.%s `+"`json:\"data\"`"+`
	Message string    	 `+"`json:\"message\"`"+`
	Status  int          `+"`json:\"status\"`"+`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create %s godoc
//
//	@Summary		Create %s
//	@Description	Create %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			payload	body		Create%sPayload		true	"%s Payload"
//	@Success		201		{object}	DataStore%sWrapper	"%s Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s  [post]
func (app *application) create%sHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	var payload Create%sPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// TODO: Change in this code
	// TODO: ADD to storage interface for use create store
	
	ctx := r.Context()
	variables := make(map[string]interface{})
	variables["created_by"] = map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if payload.Variables != nil {
		for k, v := range *payload.Variables {
			variables[k] = v
		}
	}

	resp, err := app.zeebeClient.StartWorkflow(ctx, store.%sProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.%s{
		ProcessDefinitionKey: resp.GetProcessDefinitionKey(),
		Version:              resp.GetVersion(),
		ProcessInstanceKey:   resp.GetProcessInstanceKey(),
		ResourceName:         store.%sResourceName,
		CreatedBy:            user.ID,
		TaskState:            StringPtr("CREATED"),
	}

	if err := app.store.%s.Create(ctx, model); err != nil {
		if err := app.zeebeClient.CancelWorkflow(ctx, model.ProcessInstanceKey); err != nil {
			app.internalServerError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}


// Cancel %s godoc
//
//	@Summary		Cancel %s
//	@Description	Cancel %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"ProcessInstanceKey"
//	@Success		200		{string}	string	"%s Canceled"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/{id}  [delete]
func (app *application) cancel%sHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.get%s(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	// delete model
	if err := app.store.%s.Delete(ctx, model.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: add rollback if failed to cancel in zeebe
	if err := app.zeebeClient.CancelWorkflow(ctx, model.ProcessInstanceKey); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// delete cache
	app.cacheStorage.%s.Delete(ctx, model.ID)

	if err := app.jsonResponse(w, http.StatusOK, "success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}


// GetById %s godoc
//
//	@Summary		GetById %s
//	@Description	GetById %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"ID from table"
//	@Success		200					{string}	string	"%s GetById"
//	@Failure		400					{object}	error
//	@Failure		404					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/{id}  [get]
func (app *application) getById%sHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.get%s(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) get%s(ctx context.Context, modelID int64) (*store.%s, error) {
	if !app.config.redisCfg.enabled {
		return app.store.%s.GetByID(ctx, modelID)
	}

	model, err := app.cacheStorage.%s.Get(ctx, modelID)
	if err != nil {
		return nil, err
	}

	if model == nil {
		model, err = app.store.%s.GetByID(ctx, modelID)
		if err != nil {
			return nil, err
		}
		if err := app.cacheStorage.%s.Set(ctx, model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

// GetHistoryById %s godoc
//
//	@Summary		GetHistoryById %s
//	@Description	GetHistoryById %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			id				path		int		true	"ID from table"
//	@Param			size			query		string	false	"Size 50"
//	@Param			order			query		string	false	"Order DESC ASC"
//
// @Param			type			query		string	false	"Type USER_TASK"
// @Param			state			query		string	false	"State ACTIVE"
//
//	@Param			sort			query		string	false	"Sort startDate"
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200	{string}	string	"%s GetHistoryById"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/{id}/history  [get]
func (app *application) getHistoryById%sHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	flowNodeQueryParams, err := getFlowNodeQueryParams(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if flowNodeQueryParams.Sort == "" {
		flowNodeQueryParams.Sort = "startDate"
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.get%s(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}
	
	jsonTemplate := %s
{
    "filter": {
        "processInstanceKey": %s
    }, 
    "size": %s, 
    "sort": [{"field": "%s", "order": "%s"}]
    %s
}%s

	var searchAfterStr, searchBeforeStr string

	if flowNodeQueryParams.SearchAfter != "" {
		searchAfterStr = fmt.Sprintf(%s, "searchAfter": %s%s, flowNodeQueryParams.SearchAfter)
	}

	if flowNodeQueryParams.SearchBefore != "" {
		searchBeforeStr = fmt.Sprintf(%s, "searchBefore": %s%s, flowNodeQueryParams.SearchBefore)
	}

	searchParams := searchAfterStr + searchBeforeStr

	body := []byte(fmt.Sprintf(
		jsonTemplate,
		model.ProcessInstanceKey,
		flowNodeQueryParams.Size,
		flowNodeQueryParams.Sort,
		flowNodeQueryParams.Order,
		searchParams,
	))

	// get history from operate api rest api
	url := fmt.Sprintf("%s%s/search", app.config.camundaRest.camundaOperateBaseUrl, V1FlowNodeUrl)
	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	var jsonData interface{}
	if err := json.Unmarshal(resp, &jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Search %s godoc
//
//	@Summary		Search %s
//	@Description	Search %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			limit	query		string	true	"Limit 20"
//	@Param			page	query		string	true	"Page 1"
//	@Param			search	query		string	false	"Search string"
//	@Param			sort	query		string	false	"Sort desc"
//	@Param			since	query		string	false	"Since desc"
//	@Param			until	query		string	false	"Until desc"
//	@Success		200		{string}	string	"%s Search"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s  [get]
func (app *application) search%sHandler(w http.ResponseWriter, r *http.Request) {
	pq := store.PaginatedQuery{
		Limit: 20,
		Page:  1,
		Sort:  "desc",
	}
	if err := pq.Parse(r); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(pq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// get pagination from store
	models, err := app.store.%s.Search(ctx, pq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, models); err != nil {
		app.internalServerError(w, r, err)
	}
}


// Update %s godoc
//
//	@Summary		Update %s
//	@Description	Update %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			id		path		int												true	"ID from table"
//	@Param			payload	body		Update%sPayload		true	"%s Payload"
//	@Success		200		{object}	DataStore%sWrapper	"%s Updated"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/{id}  [patch]
func (app *application) update%sHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload Update%sPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := GetUserFromContext(r)

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.get%s(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	variables := make(map[string]interface{})
	variables["updated_by"] = map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"updated_at": time.Now().Unix(),
	}

	if payload.Variables != nil {
		for k, v := range *payload.Variables {
			variables[k] = v
		}
	}

	if err := app.zeebeClient.UpdateProcessInstance(ctx, model.ProcessInstanceKey, variables); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.%s.Update(ctx, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// delete cache
	app.cacheStorage.%s.Delete(ctx, model.ID)

	if err := app.jsonResponse(w, http.StatusOK, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetProcessIncidents %s godoc
//
//	@Summary		GetProcessIncidents %s
//	@Description	GetProcessIncidents %s
//	@Tags			bpmn/%s
//	@Accept			json
//	@produce		json
//	@Param			id				path		int		true	"ID from table"
//	@Param			size			query		string	false	"Size 50"
//	@Param			order			query		string	false	"Order DESC ASC"
//
// @Param			type			query		string	false	"Type USER_TASK"
// @Param			state			query		string	false	"State ACTIVE"
//
//	@Param			sort			query		string	false	"Sort startDate"
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200	{string}	string	"%s GetProcessIncidents"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s/{id}/history  [get]
func (app *application) getIncidentsById%sHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	flowNodeQueryParams, err := getFlowNodeQueryParams(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if flowNodeQueryParams.Sort == "" {
		flowNodeQueryParams.Sort = "creationTime"
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.get%s(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}
	
	jsonTemplate := %s
{
    "filter": {
        "processInstanceKey": %s
    }, 
    "size": %s, 
    "sort": [{"field": "%s", "order": "%s"}]
    %s
}%s

	var searchAfterStr, searchBeforeStr string

	if flowNodeQueryParams.SearchAfter != "" {
		searchAfterStr = fmt.Sprintf(%s, "searchAfter": %s%s, flowNodeQueryParams.SearchAfter)
	}

	if flowNodeQueryParams.SearchBefore != "" {
		searchBeforeStr = fmt.Sprintf(%s, "searchBefore": %s%s, flowNodeQueryParams.SearchBefore)
	}

	searchParams := searchAfterStr + searchBeforeStr

	body := []byte(fmt.Sprintf(
		jsonTemplate,
		model.ProcessInstanceKey,
		flowNodeQueryParams.Size,
		flowNodeQueryParams.Sort,
		flowNodeQueryParams.Order,
		searchParams,
	))

	// get incidents from operate api rest api
	url := fmt.Sprintf("%s%s/search", app.config.camundaRest.camundaOperateBaseUrl, V1IncidentUrl)
	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	var jsonData interface{}
	if err := json.Unmarshal(resp, &jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, jsonData); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

`,
		moduleName, processName, "`", "`", processName, "`", "`",
		processName, processName, processName, processName,
		processName, processName,

		processName, processName, processName, processName,
		strings.ReplaceAll(tableName, " ", "_"), processName,
		processName, processName, processName, processName, processName,

		// cancel
		processName, processName, processName, processName,
		processName, strings.ReplaceAll(tableName, " ", "_"),
		processName, processName, processName, processName,

		// get by id
		processName, processName, processName, processName, processName,
		strings.ReplaceAll(tableName, " ", "_"), processName, processName,
		// get
		processName, processName, processName, processName, processName, processName,

		// get history
		processName, processName, processName, processName, processName,
		strings.ReplaceAll(tableName, " ", "_"), processName, processName,
		"`", "%d", "%s", "%s", "%s", "%s", "`", "`", "%s", "`",
		"`", "%s", "`", "%s", "%s",

		// search
		processName, processName, processName, processName, processName,
		strings.ReplaceAll(tableName, " ", "_"),
		processName, processName,

		// update
		processName, processName, processName, processName, processName,
		processName, processName, processName, strings.ReplaceAll(tableName, " ", "_"),
		processName, processName, processName, processName, processName,

		// incidents
		processName, processName, processName, processName, processName,
		strings.ReplaceAll(tableName, " ", "_"), processName, processName,
		"`", "%d", "%s", "%s", "%s", "%s", "`", "`", "%s", "`",
		"`", "%s", "`", "%s", "%s",
	)

	err = os.WriteFile(filePathHandler, []byte(handlerCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	modelCacheCode := fmt.Sprintf(`package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"%s/internal/store"
	"github.com/go-redis/redis/v8"
)

type %sStore struct {
	rdb *redis.Client
}

const %sExpTime = time.Hour * 24 * 7
	
func (s *%sStore) Get(ctx context.Context, modelID int64) (*store.%s, error) {
	cacheKey := fmt.Sprintf("%s-%s", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.%s
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *%sStore) Set(ctx context.Context, model *store.%s) error {
	cacheKey := fmt.Sprintf("%s-%s", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, %sExpTime).Err()
}

func (s *%sStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("%s-%s", modelID)
	s.rdb.Del(ctx, cacheKey)
}
`,
		moduleName, processName, processName, processName,
		processName, processName, "%v", processName,
		processName, processName, processName, "%v",
		processName, processName, processName, "%v",
	)

	err = os.WriteFile(filePathStoreCache, []byte(modelCacheCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// edit file storage
	generateCodeStorage := fmt.Sprintf(`
	%s interface {
		Create(context.Context, *%s) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*%s, error)
		Update(context.Context, *%s) error
		Search(context.Context, PaginatedQuery) (map[string]interface{}, error)
	}
`,
		processName, processName, processName, processName,
	)
	generateCodeConstructor := fmt.Sprintf(`
		%s:   &%sStore{db},
`, processName, processName)

	// edit file routes
	generateCodeRoutes := fmt.Sprintf(`
			r.Route("/%s", func(r chi.Router) {
				r.Get("/", app.search%sTechnologyHandler)
				r.Post("/", app.create%sHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.getById%sHandler)
					r.Delete("/", app.cancel%sHandler)
					r.Patch("/", app.update%sHandler)
					r.Get("/history", app.getHistoryById%sHandler)
					r.Get("/incidents", app.getIncidentsById%sHandler)
				})
			})	
`,
		strings.ReplaceAll(tableName, " ", "_"),
		processName, processName,
		processName, processName,
		processName, processName, processName,
	)

	// edit file cache storage
	generateCodeCacheStorage := fmt.Sprintf(`
	%s interface {
		Get(context.Context, int64) (*store.%s, error)
		Set(context.Context, *store.%s) error
		Delete(context.Context, int64)
	}
`,
		processName,
		processName,
		processName,
	)

	generateCodeCacheInterface := fmt.Sprintf(`
		%s: &%sStore{
			rdb: rbd,
		},
`,
		processName,
		processName,
	)

	err = insertGeneratedCode(filePathEditStorage, generateCodeStorage, "// GENERATED CODE INTERFACE")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditStorage, generateCodeConstructor, "// GENERATED CODE CONSTRUCTOR")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditRoutes, generateCodeRoutes, "// GENERATE ROUTES API")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditCacheStorage, generateCodeCacheStorage, "// GENERATED CACHE CODE INTERFACE")
	if err != nil {
		return err
	}

	err = insertGeneratedCode(filePathEditCacheStorage, generateCodeCacheInterface, "// GENERATED CACHE CODE CONSTRUCTOR")
	if err != nil {
		return err
	}

	return nil
}
