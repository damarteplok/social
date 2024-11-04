package main

import "net/http"

type DeployBpmnPayload struct {
	ResourceName string `json:"resource_name" validate:"required"`
}
type CrudPayload struct {
	ProcessName          string `json:"process_name" validate:"required,max=255"`
	ResourceName         string `json:"resource_name" validate:"required,max=255"`
	Version              int32  `json:"version" validate:"required"`
	ProcessDefinitionKey int64  `json:"process_definition_key" validate:"required"`
}

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
//	@Router			/camunda/deploy-crud  [post]
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

	response, err := app.zeebeClient.DeployProcessDefinition(payload.ResourceName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.zeebeClient.GenerateCRUDHandlers(response); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	jsonResponse := map[string]interface{}{
		"processDefinitionKey": response.ProcessDefinitionKey,
		"ppmnProcessId":        response.BpmnProcessId,
		"version":              response.Version,
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
//	@Router			/camunda/deploy  [post]
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

	response, err := app.zeebeClient.DeployProcessDefinition(payload.ResourceName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	jsonResponse := map[string]interface{}{
		"processDefinitionKey": response.ProcessDefinitionKey,
		"ppmnProcessId":        response.BpmnProcessId,
		"version":              response.Version,
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
//	@Router			/camunda/crud  [post]
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

	if err := app.zeebeClient.GenerateCRUDFromPayloadHandlers(payload.ProcessName, payload.ResourceName, payload.Version, payload.ProcessDefinitionKey); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "ok"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
