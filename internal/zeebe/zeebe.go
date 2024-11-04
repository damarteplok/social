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

func (c *Client) StartWorkflow(ctx context.Context, workflowName string, variables map[string]interface{}) (string, error) {
	request, err := c.client.NewCreateInstanceCommand().BPMNProcessId(workflowName).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return "", fmt.Errorf("failed to create workflow instance request: %w", err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start workflow: %w", err)
	}

	return result.String(), nil
}

func (z *Client) DeployProcessDefinition(resourceName string) (*pb.ProcessMetadata, error) {
	definition, err := MustReadFile(resourceName)
	if err != nil {
		return nil, errors.New("failed to read file")
	}
	command := z.client.NewDeployResourceCommand().AddResource(definition, resourceName)

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	resource, err := command.Send(ctx)
	if err != nil {
		return nil, err
	}

	if len(resource.GetDeployments()) < 1 {
		return nil, errors.New("failed to deploy model; nothing was deployed")
	}

	process := resource.GetDeployments()[0].GetProcess()
	if process == nil {
		return nil, errors.New("failed to deploy; the deployment was successful, but no process was returned")
	}

	return resource.GetDeployments()[0].GetProcess(), nil
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
	if err := generateCrud(processNameTitle, resourceName, tableName, version, processDefinitionKey); err != nil {
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

	if err := generateCrud(processName, resourceName, tableName, version, processDefinitionKey); err != nil {
		return err
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
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func generateCrud(processName, resourceName, tableName string, version int32, processDefinitionKey int64) error {
	filePathHandler := fmt.Sprintf("./cmd/api/%s.go", tableName)
	filePathStore := fmt.Sprintf("./internal/store/%s.go", tableName)

	modelCode := fmt.Sprintf(`package store
import (
	"context"
	"database/sql"
	"errors"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type %s struct {
    ID int64 `+"`json:\"id\"`"+`
	ProcessDefinitionKey int64  `+"`json:\"process_definition_key\"`"+`
	Version int32 `+"`json:\"version\"`"+`
	ResourceName string `+"`json:\"resource_name\"`"+`
	CreatedAt string `+"`json:\"created_at\"`"+`
	UpdatedAt string `+"`json:\"updated_at\"`"+`
}

type %sStore struct {
	db *sql.DB
}
	
func (s *%sStore) Create(ctx context.Context, model *%s) error {
	model.Version = %d
	model.ProcessDefinitionKey = %d
	model.ResourceName = "%s"

	query := %s
		INSERT INTO %s (process_definition_key, version, resource_name)
		VALUES (
			$1, 
			$2, 
			$3
		) RETURNING 
		 	id, process_definition_key, version, resource_name, 
			created_at, updated_at
		%s
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
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
		SELECT id, process_definition_key, version, resource_name, created_at, updated_at
		FROM %s
		WHERE id = $1;
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model %s
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
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

func (s *%sStore) Delete(ctx context.Context, id int64) error {
	query := %sDELETE FROM %s WHERE id = $1;%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
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

func (s *%sStore) Update(ctx context.Context, model *%s) error {
	query := %s
		UPDATE %s
		SET process_definition_key = $1, version = $2, resource_name = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, process_definition_key, version, resource_name, created_at updated_at;
	%s

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ID,
	).Scan(&model.ID, &model.ProcessDefinitionKey, &model.Version, &model.ResourceName, &model.CreatedAt, &model.UpdatedAt)
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

	err := os.WriteFile(filePathStore, []byte(modelCode), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	handlerCode := fmt.Sprintf(`package main

import (
    "net/http"
	"github.com/damarteplok/social/internal/store"
)

type Create%sPayload struct {}
type Update%sPayload struct {}
type DataStore%sWrapper struct {
	Data store.%s `+"`json:\"data\"`"+`
}

// Create %s godoc
//
//	@Summary		Create %s
//	@Description	Create %s
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		Create%sPayload		true	"%s Payload"
//	@Success		201		{object}	DataStore%sWrapper	"%s Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/%s  [post]
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
		processName,
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
