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

type CreatePembuatanMediaBeritaTechnologyPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type UpdatePembuatanMediaBeritaTechnologyPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type DataStorePembuatanMediaBeritaTechnologyWrapper struct {
	Data    store.PembuatanMediaBeritaTechnology `json:"data"`
	Message string                               `json:"message"`
	Status  int                                  `json:"status"`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create PembuatanMediaBeritaTechnology godoc
//
//	@Summary		Create PembuatanMediaBeritaTechnology
//	@Description	Create PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreatePembuatanMediaBeritaTechnologyPayload		true	"PembuatanMediaBeritaTechnology Payload"
//	@Success		201		{object}	DataStorePembuatanMediaBeritaTechnologyWrapper	"PembuatanMediaBeritaTechnology Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology  [post]
func (app *application) createPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	var payload CreatePembuatanMediaBeritaTechnologyPayload
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

	resp, err := app.zeebeClient.StartWorkflow(ctx, store.PembuatanMediaBeritaTechnologyProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.PembuatanMediaBeritaTechnology{
		ProcessDefinitionKey: resp.GetProcessDefinitionKey(),
		Version:              resp.GetVersion(),
		ProcessInstanceKey:   resp.GetProcessInstanceKey(),
		ResourceName:         store.PembuatanMediaBeritaTechnologyResourceName,
		CreatedBy:            user.ID,
		TaskState:            StringPtr("CREATED"),
	}

	if err := app.store.PembuatanMediaBeritaTechnology.Create(ctx, model); err != nil {
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

// Cancel PembuatanMediaBeritaTechnology godoc
//
//	@Summary		Cancel PembuatanMediaBeritaTechnology
//	@Description	Cancel PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"ProcessInstanceKey"
//	@Success		200	{string}	string	"PembuatanMediaBeritaTechnology Canceled"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology/{id}  [delete]
func (app *application) cancelPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.getPembuatanMediaBeritaTechnology(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	// delete model
	if err := app.store.PembuatanMediaBeritaTechnology.Delete(ctx, model.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: add rollback if failed to cancel in zeebe
	if err := app.zeebeClient.CancelWorkflow(ctx, model.ProcessInstanceKey); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// delete cache
	app.cacheStorage.PembuatanMediaBeritaTechnology.Delete(ctx, model.ID)

	if err := app.jsonResponse(w, http.StatusOK, "success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetById PembuatanMediaBeritaTechnology godoc
//
//	@Summary		GetById PembuatanMediaBeritaTechnology
//	@Description	GetById PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"ID from table"
//	@Success		200	{string}	string	"PembuatanMediaBeritaTechnology GetById"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology/{id}  [get]
func (app *application) getByIdPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.getPembuatanMediaBeritaTechnology(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	// get zeebe client untuk mendapatkan detail task
	url := fmt.Sprintf("%s%s/%d", app.config.camundaRest.camundaOperateBaseUrl, V1ProcessInstanceUrl, model.ProcessInstanceKey)
	resp, err := app.zeebeClientRest.SendRequest(
		ctx,
		http.MethodGet,
		url,
		bytes.NewBuffer([]byte("{}")),
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

	if err := app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"model":   model,
		"camunda": jsonData,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPembuatanMediaBeritaTechnology(ctx context.Context, modelID int64) (*store.PembuatanMediaBeritaTechnology, error) {
	if !app.config.redisCfg.enabled {
		return app.store.PembuatanMediaBeritaTechnology.GetByID(ctx, modelID)
	}

	model, err := app.cacheStorage.PembuatanMediaBeritaTechnology.Get(ctx, modelID)
	if err != nil {
		return nil, err
	}

	if model == nil {
		model, err = app.store.PembuatanMediaBeritaTechnology.GetByID(ctx, modelID)
		if err != nil {
			return nil, err
		}
		if err := app.cacheStorage.PembuatanMediaBeritaTechnology.Set(ctx, model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

// GetHistoryById PembuatanMediaBeritaTechnology godoc
//
//	@Summary		GetHistoryById PembuatanMediaBeritaTechnology
//	@Description	GetHistoryById PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			id				path		int		true	"ID from table"
//	@Param			size			query		string	false	"Size 50"
//	@Param			order			query		string	false	"Order DESC ASC"
//
//	@Param			type			query		string	false	"Type USER_TASK"
//	@Param			state			query		string	false	"State ACTIVE"
//
//	@Param			sort			query		string	false	"Sort startDate"
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200				{string}	string	"PembuatanMediaBeritaTechnology GetHistoryById"
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology/{id}/history  [get]
func (app *application) getHistoryByIdPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
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

	model, err := app.getPembuatanMediaBeritaTechnology(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	jsonTemplate := `
{
    "filter": {
        "processInstanceKey": %d
    }, 
    "size": %s, 
    "sort": [{"field": "%s", "order": "%s"}]
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

// Search PembuatanMediaBeritaTechnology godoc
//
//	@Summary		Search PembuatanMediaBeritaTechnology
//	@Description	Search PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			limit	query		string	true	"Limit 20"
//	@Param			page	query		string	true	"Page 1"
//	@Param			search	query		string	false	"Search string"
//	@Param			sort	query		string	false	"Sort desc"
//	@Param			since	query		string	false	"Since desc"
//	@Param			until	query		string	false	"Until desc"
//	@Success		200		{string}	string	"PembuatanMediaBeritaTechnology Search"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology  [get]
func (app *application) searchPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
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
	models, err := app.store.PembuatanMediaBeritaTechnology.Search(ctx, pq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, models); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Update PembuatanMediaBeritaTechnology godoc
//
//	@Summary		Update PembuatanMediaBeritaTechnology
//	@Description	Update PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			id		path		int												true	"ID from table"
//	@Param			payload	body		UpdatePembuatanMediaBeritaTechnologyPayload		true	"PembuatanMediaBeritaTechnology Payload"
//	@Success		200		{object}	DataStorePembuatanMediaBeritaTechnologyWrapper	"PembuatanMediaBeritaTechnology Updated"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology/{id}  [patch]
func (app *application) updatePembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload UpdatePembuatanMediaBeritaTechnologyPayload
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

	model, err := app.getPembuatanMediaBeritaTechnology(ctx, id)
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

	if err := app.store.PembuatanMediaBeritaTechnology.Update(ctx, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// delete cache
	app.cacheStorage.PembuatanMediaBeritaTechnology.Delete(ctx, model.ID)

	if err := app.jsonResponse(w, http.StatusOK, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetProcessIncidents PembuatanMediaBeritaTechnology godoc
//
//	@Summary		GetProcessIncidents PembuatanMediaBeritaTechnology
//	@Description	GetProcessIncidents PembuatanMediaBeritaTechnology
//	@Tags			bpmn/PembuatanMediaBeritaTechnology
//	@Accept			json
//	@produce		json
//	@Param			id				path		int		true	"ID from table"
//	@Param			size			query		string	false	"Size 50"
//	@Param			order			query		string	false	"Order DESC ASC"
//
//	@Param			type			query		string	false	"Type USER_TASK"
//	@Param			state			query		string	false	"State ACTIVE"
//
//	@Param			sort			query		string	false	"Sort startDate"
//	@Param			searchAfter		query		string	false	"SearchAfter 1731486859777,2251799814109407"
//	@Param			searchBefore	query		string	false	"SearchBefore 1731486859777,2251799814109407"
//	@Success		200				{string}	string	"PembuatanMediaBeritaTechnology GetProcessIncidents"
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatan_media_berita_technology/{id}/incidents  [get]
func (app *application) getIncidentsByIdPembuatanMediaBeritaTechnologyHandler(w http.ResponseWriter, r *http.Request) {
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

	model, err := app.getPembuatanMediaBeritaTechnology(ctx, id)
	if err != nil {
		app.handleRequestError(w, r, err)
		return
	}

	jsonTemplate := `
{
    "filter": {
        "processInstanceKey": %d
    }, 
    "size": %s, 
    "sort": [{"field": "%s", "order": "%s"}]
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
		model.ProcessInstanceKey,
		flowNodeQueryParams.Size,
		flowNodeQueryParams.Sort,
		flowNodeQueryParams.Order,
		searchParams,
	))

	// get history from operate api rest api
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
