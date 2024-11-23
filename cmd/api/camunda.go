package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-chi/chi/v5"
)

const (
	BucketBPMN             = "bpmn"
	BucketForm             = "form"
	StateCreated           = "CREATED"
	StateCompleted         = "COMPLETED"
	StateCanceled          = "CANCELED"
	StateFailed            = "FAILED"
	ResourceUrl            = "/v2/resources"
	ProcessInstanceUrl     = "/v2/process-instances"
	V1TasklistUrl          = "/v1/tasks"
	V1FlowNodeUrl          = "/v1/flownode-instances"
	V1ProcessInstanceUrl   = "/v1/process-instances"
	V1ProcessDefinitionUrl = "/v1/process-definitions"
	V1IncidentUrl          = "/v1/incidents"
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

// upload upload godoc
//
//	@Summary		Upload Bpmn Camunda
//	@Description	Upload Bpmn Camunda by Name In Folder Resources
//	@Tags			camunda/minio
//	@Accept			multipart/form-data
//	@produce		json
//	@Param			file	formData	file	true	"File Upload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/minio/upload  [post]
func (app *application) uploadCamundaHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	tempFile, err := os.CreateTemp("", header.Filename)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, file); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	uploadInfo, err := app.minioClient.UploadBpmnOrForm(ctx, tempFile, header.Filename)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTypeNotAllowed):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"bucket":   uploadInfo.Bucket,
		"Key":      uploadInfo.Key,
		"Location": uploadInfo.Location,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// uploadMultiple godoc
//
//	@Summary		Upload Multiple Bpmn Camunda
//	@Description	Upload Multiple Bpmn Camunda by Name In Folder Resources
//	@Tags			camunda/minio
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			files	formData	file	true	"Files Upload" collectionFormat(multi)
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/minio/upload-multiple [post]
func (app *application) uploadMultipleCamundaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		app.badRequestResponse(w, r, errors.New("no files found"))
		return
	}

	// check files type
	for _, fileHeader := range files {
		if filepath.Ext(fileHeader.Filename) != ".bpmn" && filepath.Ext(fileHeader.Filename) != ".form" {
			app.badRequestResponse(w, r, errors.New("file type not allowed"))
			return
		}
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	var uploadInfos []map[string]interface{}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("", fileHeader.Filename)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		defer os.Remove(tempFile.Name())

		if _, err := io.Copy(tempFile, file); err != nil {
			app.internalServerError(w, r, err)
			return
		}

		uploadInfo, err := app.minioClient.UploadBpmnOrForm(ctx, tempFile, fileHeader.Filename)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrTypeNotAllowed):
				app.badRequestResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		uploadInfos = append(uploadInfos, map[string]interface{}{
			"bucket":   uploadInfo.Bucket,
			"Key":      uploadInfo.Key,
			"Location": uploadInfo.Location,
		})
	}

	if err := app.jsonResponse(w, http.StatusCreated, uploadInfos); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Deploy godoc
//
//	@Summary		Deploy Bpmn Camunda and Create CRUD in Store And Handler File
//	@Description	Deploy Bpmn Camunda by Name From MINIO And Create CRUD in Store And Handler File
//	@Tags			camunda/minio
//	@Accept			json
//	@produce		json
//	@Param			payload	body		DeployBpmnPayload	true	"Deploy Bpmn Payload"
//	@Success		201		{string}	string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/minio/deploy-crud  [post]
func (app *application) getObjectFromMinioThanUseItHandler(w http.ResponseWriter, r *http.Request) {
	var payload DeployBpmnPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	fileResource, err := app.minioClient.GetObject(ctx, BucketBPMN, payload.ResourceName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	defer fileResource.Close()

	var formFiles []*os.File
	var tmpFiles []*os.File
	if payload.FormResources != nil {
		for _, formResource := range payload.FormResources {
			formFile, err := app.minioClient.GetObject(ctx, BucketForm, formResource)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}
			defer formFile.Close()

			tmpFile, err := os.CreateTemp("", "form-*.tmp")
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}
			defer tmpFile.Close()

			if _, err := io.Copy(tmpFile, formFile); err != nil {
				app.internalServerError(w, r, err)
				return
			}

			formFiles = append(formFiles, tmpFile)
			tmpFiles = append(tmpFiles, tmpFile)
		}
	}

	response, bpmnProcess, err := app.zeebeClient.DeployProcessDefinitionFromFiles(fileResource, formFiles)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	for _, tmpFile := range tmpFiles {
		os.Remove(tmpFile.Name())
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
//
//	@Param			bpmnProcessId				query		string	false	"Bpmn Process Id"
//	@Param			processDefinitionKey		query		string	false	"Process Definition Key"
//	@Param			parentProcessInstanceKey	query		string	false	"Parent Process Instance Key"
//	@Param			startDate					query		string	false	"Start Date"
//	@Param			endDate						query		string	false	"End Date"
//	@Param			state						query		string	false	"State"
//
//	@Param			size						query		string	false	"Size 50"
//	@Param			searchAfter					query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore				query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200							{string}	string	"search process instance"
//	@Failure		400							{object}	error
//	@Failure		500							{object}	error
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
	"filter": {
		%s
    }, 
    "size": %s
    %s
}`
	var bpmnProcessIdStr, processDefinitionKeyStr, parentProcessInstanceKeyStr, startDateStr, endDateStr, stateStr string
	if flowNodeQueryParams.BpmnProcessId != "" {
		bpmnProcessIdStr = fmt.Sprintf(`, "bpmnProcessId": "%s"`, flowNodeQueryParams.BpmnProcessId)
	}
	if flowNodeQueryParams.ProcessDefinitionKey != "" {
		processDefinitionKeyStr = fmt.Sprintf(`, "processDefinitionKey": %s`, flowNodeQueryParams.ProcessDefinitionKey)
	}
	if flowNodeQueryParams.ParentProcessInstanceKey != "" {
		parentProcessInstanceKeyStr = fmt.Sprintf(`, "parentProcessInstanceKey": %s`, flowNodeQueryParams.ParentProcessInstanceKey)
	}
	if flowNodeQueryParams.StartDate != "" {
		startDateStr = fmt.Sprintf(`, "startDate": "%s"`, flowNodeQueryParams.StartDate)
	}
	if flowNodeQueryParams.EndDate != "" {
		endDateStr = fmt.Sprintf(`, "endDate": "%s"`, flowNodeQueryParams.EndDate)
	}
	if flowNodeQueryParams.State != "" {
		stateStr = fmt.Sprintf(`, "state": "%s"`, flowNodeQueryParams.State)
	}

	filterForm := bpmnProcessIdStr + processDefinitionKeyStr + parentProcessInstanceKeyStr + startDateStr + endDateStr + stateStr

	if filterForm != "" {
		filterForm = filterForm[2:]
	}
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
		filterForm,
		flowNodeQueryParams.Size,
		searchParams,
	))

	log.Println(string(body))

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

// GetTEXT/XML godoc
//
//	@Summary		Get TEXT/XML from rest api
//	@Description	Get TEXT/XML from rest api
//	@Tags			camunda/resource
//	@Accept			json
//	@produce		text/xml
//
//	@Param			processDefinitionKey	path		int		true	"Process Definition Key"
//
//	@Success		200						{string}	string	"search process instance"
//	@Failure		400						{object}	error
//	@Failure		500						{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/resource/{processDefinitionKey}/xml  [get]
func (app *application) xmlCamundaHandler(w http.ResponseWriter, r *http.Request) {
	processDefinitionKey, err := strconv.ParseInt(chi.URLParam(r, "processDefinitionKey"), 10, 64)
	if err != nil || processDefinitionKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	url := fmt.Sprintf("%s%s/%d/xml",
		app.config.camundaRest.camundaOperateBaseUrl,
		V1ProcessDefinitionUrl,
		processDefinitionKey,
	)
	log.Println(url)
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

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		app.internalServerError(w, r, err)
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
