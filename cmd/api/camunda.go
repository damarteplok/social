package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-chi/chi/v5"
)

const (
	StateCreated         = "CREATED"
	StateCompleted       = "COMPLETED"
	StateCanceled        = "CANCELED"
	StateFailed          = "FAILED"
	ResourceUrl          = "/v2/resources"
	ProcessInstanceUrl   = "/v2/process-instances"
	V1TasklistUrl        = "/v1/tasks"
	V1FlowNodeUrl        = "/v1/flownode-instances"
	V1ProcessInstanceUrl = "/v1/process-instances"
	V1IncidentUrl        = "/v1/incidents"
)

// Deploy godoc
//
//	@Summary		Deploy Bpmn Camunda and Create CRUD in Store And Handler File
//	@Description	Deploy Bpmn Camunda by Name In Folder Resources And Create CRUD in Store And Handler File
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		json
//	@Param			payload	body		DeployBpmnPayload	true	"Deploy Bpmn Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/deploy-crud  [post]
func (app *application) deployCamundaHandler(w http.ResponseWriter, r *http.Request) {
	var payload DeployBpmnPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	response, bpmnProcess, err := app.zeebeClient.DeployProcessDefinition(payload.ResourceName, payload.FormResources)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.zeebeClient.GenerateCRUDUserTaskServiceTaskHandler(&bpmnProcess); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	processes := make([]map[string]interface{}, len(response))
	for i, process := range response {
		if err := app.zeebeClient.GenerateCRUDHandlers(process); err != nil {
			app.internalServerError(w, r, err)
			return
		}
		processes[i] = map[string]interface{}{
			"processDefinitionKey": process.ProcessDefinitionKey,
			"bpmnProcessId":        process.BpmnProcessId,
			"version":              process.Version,
		}
	}

	jsonResponse := map[string]interface{}{
		"processes": processes,
	}

	if err := app.jsonResponse(w, http.StatusCreated, jsonResponse); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Deploy Only godoc
//
//	@Summary		Deploy Only Bpmn Camunda
//	@Description	Deploy Only Bpmn Camunda by Name In Folder Resources
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		json
//	@Param			payload	body		DeployBpmnPayload	true	"Deploy Bpmn Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/deploy  [post]
func (app *application) deployOnlyCamundaHandler(w http.ResponseWriter, r *http.Request) {
	var payload DeployBpmnPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	response, _, err := app.zeebeClient.DeployProcessDefinition(payload.ResourceName, payload.FormResources)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	processes := make([]map[string]interface{}, len(response))
	for i, process := range response {
		processes[i] = map[string]interface{}{
			"processDefinitionKey": process.ProcessDefinitionKey,
			"bpmnProcessId":        process.BpmnProcessId,
			"version":              process.Version,
		}
	}

	jsonResponse := map[string]interface{}{
		"processes": processes,
	}

	if err := app.jsonResponse(w, http.StatusCreated, jsonResponse); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// CRUD godoc
//
//	@Summary		CRUD Store And Handler
//	@Description	CRUD Store And Handler from Payload
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CrudPayload	true	"Crud Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/crud  [post]
func (app *application) crudCamundaHandler(w http.ResponseWriter, r *http.Request) {
	var payload CrudPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.zeebeClient.GenerateCRUDFromPayloadHandlers(
		payload.ProcessName,
		payload.ResourceName,
		payload.Version,
		payload.ProcessDefinitionKey,
	); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "ok"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Create Proses Instance godoc
//
//	@Summary		Create Proses Instance from rest api
//	@Description	Create Proses Instance from rest api
//	@Tags			camunda/process-instance
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreateProcessInstancePayload	true	"Create Proses Instance Payload"
//	@Success		200		{object}	CreateProcessInstancesResponse
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/process-instance  [post]
func (app *application) createProsesInstance(w http.ResponseWriter, r *http.Request) {
	var payload CreateProcessInstancePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		app.config.camundaRest.zeebeRestAddress+ProcessInstanceUrl,
		bytes.NewBuffer(body),
	)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var processInstanceResp CreateProcessInstancesResponse
	if err := json.Unmarshal(resp, &processInstanceResp); err != nil {
		app.internalServerError(w, r, fmt.Errorf("failed to unmarshal response: %w", err))
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, processInstanceResp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Cancel Proses Instance godoc
//
//	@Summary		Cancel Proses Instance from rest api
//	@Description	Cancel Proses Instance from rest api
//	@Tags			camunda/process-instance
//	@Accept			json
//	@produce		json
//	@Param			processinstanceKey	path		int	true	"processinstanceskey"
//	@Success		204					{string}	string
//	@Failure		400					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/process-instance/{processinstanceKey}/cancel  [post]
func (app *application) cancelProcessInstance(w http.ResponseWriter, r *http.Request) {
	processInstanceKey, err := strconv.ParseInt(chi.URLParam(r, "processinstanceKey"), 10, 64)
	if err != nil || processInstanceKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	url := fmt.Sprintf("%s/%d/cancellation", ProcessInstanceUrl, processInstanceKey)
	_, err = app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		app.config.camundaRest.zeebeRestAddress+url,
		bytes.NewBufferString("{}"),
	)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "cancelled success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Search TaskList godoc
//
//	@Summary		Search TaskList from rest api v1
//	@Description	Search TaskList from rest api v1
//	@Tags			camunda/user-task
//	@Accept			json
//	@produce		json
//	@Param			payload	body		SearchTaskListPayload	true	"Search TaskList Payload"
//	@Success		200		{object}	SearchTaskListPayload
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/user-task  [post]
func (app *application) searchTaskListHandler(w http.ResponseWriter, r *http.Request) {
	var payload SearchTaskListPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.setDefaultSort(&payload)
	app.setDefaultState(&payload)

	body, err := json.Marshal(payload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

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
}

// Delete godoc
//
//	@Summary		Delete Bpmn Camunda
//	@Description	Delete Bpmn Camunda
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		json
//
//	@Param			processDefinitionKey	path		int	true	"Process Definition Key"
//
//	@Success		200						{string}	string
//
//	@Failure		404						{object}	error
//
//	@Failure		400						{object}	error
//	@Failure		500						{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/{processDefinitionKey}/delete  [delete]
func (app *application) deleteCamundaHandler(w http.ResponseWriter, r *http.Request) {
	processDefinitionKey, err := strconv.ParseInt(chi.URLParam(r, "processDefinitionKey"), 10, 64)
	if err != nil || processDefinitionKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	_, err = app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%d/deletion", app.config.camundaRest.zeebeRestAddress+ResourceUrl, processDefinitionKey),
		nil,
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "deleted success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Resolve Incident godoc
//
//	@Summary		Resolve Incident Bpmn Camunda
//	@Description	Resolve Incident Bpmn Camunda
//	@Tags			camunda/incident
//	@Accept			json
//	@produce		json
//
//	@Param			incidentKey	path		int	true	"incidentKey"
//
//	@Success		200			{string}	string
//
//	@Failure		404			{object}	error
//
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/incident/{incidentKey}/resolve  [post]
func (app *application) resolveIncidentHandler(w http.ResponseWriter, r *http.Request) {
	incidentKey, err := strconv.ParseInt(chi.URLParam(r, "incidentKey"), 10, 64)
	if err != nil || incidentKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	_, err = app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%d/resolution", app.config.camundaRest.zeebeRestAddress+"/v2/incidents", incidentKey),
		nil,
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "resolved success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Search User Task godoc
//
//	@Summary		Search User Task from rest api v2 must enabled in camunda-platform config first
//	@Description	Search User Task from rest api v2 must enabled in camunda-platform config first
//	@Tags			camunda/user-task
//	@Accept			json
//	@produce		json
//	@Param			payload	body		QueryUserTaskPayload	true	"Query User Task Payload"
//	@Success		200		{object}	QueryUserTaskPayload
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/user-task/search  [post]
func (app *application) searchUserTaskHandler(w http.ResponseWriter, r *http.Request) {
	var payload QueryUserTaskPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	url := fmt.Sprintf("%s/v2/user-tasks/search", app.config.camundaRest.zeebeRestAddress)
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

// GetProsesInstance godoc
//
//	@Summary		Get Proses Instance from rest api
//	@Description	Get Proses Instance from rest api
//	@Tags			camunda/process-instance
//	@Accept			json
//	@produce		json
//	@Param			size			query		string	false	"Size 50"
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200				{string}	string	"search process instance"
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/process-instance  [get]
func (app *application) searchProcessInstance(w http.ResponseWriter, r *http.Request) {
	flowNodeQueryParams, err := getFlowNodeQueryParams(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	jsonTemplate := `
{
    "size": %s
    %s
}`

	var searchAfterStr, searchBeforeStr string

	if flowNodeQueryParams.SearchAfter != "" {
		searchAfterStr = fmt.Sprintf(`, "searchAfter": %s`, flowNodeQueryParams.SearchAfter)
	}

	if flowNodeQueryParams.SearchBefore != "" {
		searchBeforeStr = fmt.Sprintf(`, "searchBefore": %s`, flowNodeQueryParams.SearchBefore)
	}

	searchParams := searchAfterStr + searchBeforeStr

	body := []byte(fmt.Sprintf(
		jsonTemplate,
		flowNodeQueryParams.Size,
		searchParams,
	))

	url := fmt.Sprintf("%s%s/search", app.config.camundaRest.camundaOperateBaseUrl, V1ProcessInstanceUrl)
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

// GetOperateApi godoc
//
//	@Summary		Get Information operate statistics camunda from rest api
//	@Description	Get Information operate statistics camunda from rest api
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		json
//	@Success		200	{string}	string	"search process instance"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/operate/statistics  [get]
func (app *application) operateStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// get core statistics
	url := fmt.Sprintf("%s/api/process-instances/core-statistics", app.config.camundaRest.camundaOperateBaseUrl)
	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	var coreStats OperateCoreStats
	if err := json.Unmarshal(resp, &coreStats); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// get process
	urlProcess := fmt.Sprintf("%s/api/incidents/byProcess", app.config.camundaRest.camundaOperateBaseUrl)
	respProcess, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodGet,
		urlProcess,
		nil,
	)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	var processStats []OperateProcessStats
	if err := json.Unmarshal(respProcess, &processStats); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"stats":   coreStats,
		"process": processStats,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) setDefaultSort(payload *SearchTaskListPayload) {
	if len(payload.Sort) > 1 {
		for i := range payload.Sort {
			if payload.Sort[i].Field == "" {
				payload.Sort[i].Field = "creationTime"
			}
			if payload.Sort[i].Order == "" {
				payload.Sort[i].Order = "DESC"
			}
		}
	} else {
		payload.Sort = append(payload.Sort, SortSearchTasklist{
			Field: "creationTime",
			Order: "DESC",
		})
	}
}

func (app *application) setDefaultState(payload *SearchTaskListPayload) {
	if payload.State == "" {
		payload.State = "CREATED"
	}
}
