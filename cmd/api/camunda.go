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
	StateCreated       = "CREATED"
	StateCompleted     = "COMPLETED"
	StateCanceled      = "CANCELED"
	StateFailed        = "FAILED"
	ProcessInstanceUrl = "/v2/process-instances"
	V1TasklistUrl      = "/v1/tasks"
)

// Deploy godoc
//
//	@Summary		Deploy Bpmn Camunda
//	@Description	Deploy Bpmn Camunda by Name In Folder Resources And Create CRUD in Store And Handler File
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		DeployBpmnPayload	true	"Deploy Bpmn Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		BasicAuth
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
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		DeployBpmnPayload	true	"Deploy Bpmn Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		BasicAuth
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
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CrudPayload	true	"Crud Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		BasicAuth
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
//	@Summary		Create Proses Instance form rest api
//	@Description	Create Proses Instance form rest api
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreateProcessInstancePayload	true	"Create Proses Instance Payload"
//	@Success		200		{object}	CreateProcessInstancesResponse
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		BasicAuth
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
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resp, err := app.zeebeClientRest.SendRequest(ctx, "POST", ProcessInstanceUrl, bytes.NewBuffer(body))
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
//	@Summary		Cancel Proses Instance form rest api
//	@Description	Cancel Proses Instance form rest api
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			processinstanceKey	path		int	true	"processinstanceskey"
//	@Success		204					{string}	string
//	@Failure		400					{object}	error
//	@Failure		500					{object}	error
//	@Security		BasicAuth
//	@Router			/camunda/process-instance/{processinstanceKey}/cancel  [post]
func (app *application) cancelProcessInstance(w http.ResponseWriter, r *http.Request) {
	processInstanceKey, err := strconv.ParseInt(chi.URLParam(r, "processinstanceKey"), 10, 64)
	if err != nil || processInstanceKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	url := ProcessInstanceUrl + "/" + strconv.Itoa(int(processInstanceKey)) + "/cancellation"
	_, err = app.zeebeClientRest.SendRequest(ctx, "POST", url, bytes.NewBufferString("{}"))
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
//	@Summary		Search TaskList form rest api
//	@Description	Search TaskList form rest api
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		SearchTaskListPayload	true	"Search TaskList Payload"
//	@Success		200		{object}	SearchTaskListPayload
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		BasicAuth
//	@Router			/camunda/tasklist  [post]
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

	body, err := json.Marshal(payload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	url := V1TasklistUrl + "/search"
	resp, err := app.zeebeClientRest.SendRequest(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
