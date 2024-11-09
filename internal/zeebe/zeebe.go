package zeebe

import (
	"context"
	"embed"
	"encoding/xml"
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

func (c *Client) StartWorkflow(ctx context.Context, processDefinitionKey int64, variables map[string]interface{}) (*pb.CreateProcessInstanceWithResultResponse, error) {
	request, err := c.client.NewCreateInstanceCommand().ProcessDefinitionKey(processDefinitionKey).VariablesFromMap(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow instance request: %w", err)
	}

	result, err := request.WithResult().Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start workflow: %w", err)
	}

	return result, nil
}

func (c *Client) CancelWorkflow(ctx context.Context, processInstanceKey int64) error {
	_, err := c.client.NewCancelInstanceCommand().ProcessInstanceKey(processInstanceKey).Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to cancel workflow: %w", err)
	}
	return nil
}

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

func (c *Client) GenerateCRUDFromPayloadHandlers(processName, resourceName string, version int32, processDefinitionKey int64) error {
	processNameTitle := toCamelCase(processName)
	tableName := processName
	if err := generateCrudProcess(processNameTitle, resourceName, tableName, version, processDefinitionKey); err != nil {
		return err
	}

	return nil
}

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

func MustReadFile(resourceFile string) ([]byte, error) {
	contents, err := res.ReadFile("resources/" + resourceFile)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func toCamelCase(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func unMarshalBpmn(bpmnContent []byte) ([]BPMNProcess, error) {
	var bpmn BPMNDocument
	if err := xml.Unmarshal(bpmnContent, &bpmn); err != nil {
		return nil, fmt.Errorf("failed to parse BPMN content: %w", err)
	}
	return bpmn.Processes, nil
}

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

	filePathScripts := fmt.Sprintf("./scripts/%s_user_task.sql", nameFile)
	scriptCode := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE
);

DROP TABLE IF EXISTS %s;
	`, nameFile, nameFile)

	err := os.WriteFile(filePathScripts, []byte(scriptCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	filePathStore := fmt.Sprintf("./internal/store/%s_user_task.go", nameFile)
	modelCode := fmt.Sprintf(`package store
import (
	"context"
	"database/sql"
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
    ID int64 `+"`json:\"id\"`"+`
	Name string  `+"`json:\"name\"`"+`
	CreatedAt string `+"`json:\"created_at\"`"+`
	UpdatedAt string `+"`json:\"updated_at\"`"+`
	DeletedAt *string `+"`json:\"deleted_at\"`"+`
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
	query := %s
		INSERT INTO %s (name)
		VALUES (
			$1
		) RETURNING 
		 	id, name, 
			created_at, updated_at
		%s
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
	).Scan(
		&model.ID,
		&model.Name,
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
		SELECT id, name, created_at, updated_at
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model %s
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
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
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, name, created_at updated_at;
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.ID,
	).Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt)
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

	return nil
}

func generateCrudProcess(processName, resourceName, tableName string, version int32, processDefinitionKey int64) error {
	filePathHandler := fmt.Sprintf("./cmd/api/%s_process.go", tableName)
	filePathStore := fmt.Sprintf("./internal/store/%s_process.go", tableName)
	filePathScripts := fmt.Sprintf("./scripts/%s_process.sql", tableName)

	scriptCode := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id BIGSERIAL PRIMARY KEY,
	process_definition_key BIGINT NOT NULL,
	version INT NOT NULL,
	resource_name VARCHAR(256) NOT NULL,
	process_instance_key BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE
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
    ID int64 `+"`json:\"id\"`"+`
	ProcessDefinitionKey int64  `+"`json:\"process_definition_key\"`"+`
	Version int32 `+"`json:\"version\"`"+`
	ResourceName string `+"`json:\"resource_name\"`"+`
	ProcessInstanceKey int64 `+"`json:\"process_instance_key\"`"+`
	CreatedAt string `+"`json:\"created_at\"`"+`
	UpdatedAt string `+"`json:\"updated_at\"`"+`
	DeletedAt *string `+"`json:\"deleted_at\"`"+`
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
		INSERT INTO %s (process_definition_key, version, resource_name, process_instance_key)
		VALUES (
			$1, 
			$2, 
			$3,
			$4
		) RETURNING 
		 	id, process_definition_key, version, resource_name, process_instance_key,
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
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
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
		SELECT id, process_definition_key, version, resource_name, process_instance_key, created_at, updated_at
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
		SET process_definition_key = $1, version = $2, resource_name = $3, process_instance_key = $4, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, process_definition_key, version, resource_name, process_instance_key, created_at updated_at;
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
		model.ID,
	).Scan(&model.ID, &model.ProcessDefinitionKey, &model.Version, &model.ResourceName, &model.ProcessInstanceKey, &model.CreatedAt, &model.UpdatedAt)
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
	)

	err = os.WriteFile(filePathStore, []byte(modelCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	handlerCode := fmt.Sprintf(`package main

import (
    "net/http"
	"github.com/damarteplok/social/internal/store"
)

type Create%sPayload struct {
	Variables   		 *map[string]string  %sjson:"variables,omitempty"%s
}
type Update%sPayload struct {
	Variables            *map[string]string  %sjson:"variables,omitempty"%s
}
type DataStore%sWrapper struct {
	Data store.%s `+"`json:\"data\"`"+`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create %s godoc
//
//	@Summary		Create %s
//	@Description	Create %s
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			payload	body		Create%sPayload		true	"%s Payload"
//	@Success		201		{object}	DataStore%sWrapper	"%s Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/%s  [post]
func (app *application) create%sHandler(w http.ResponseWriter, r *http.Request) {
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

	//model := &store.%s{}

	//ctx := r.Context()

	// if err := app.store.%s.Create(ctx, model); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	//if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
	//	app.internalServerError(w, r, err)
	//	return
	//}
}
`,
		processName,
		"`",
		"`",
		processName,
		"`",
		"`",
		processName,
		processName,
		processName,
		processName,
		processName,
		processName,
		processName,
		processName,
		processName,
		tableName,
		processName,
		processName,
		processName,
		processName,
	)

	err = os.WriteFile(filePathHandler, []byte(handlerCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	return nil
}
