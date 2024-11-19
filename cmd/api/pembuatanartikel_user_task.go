package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damarteplok/social/internal/store"
)

type FormDataPembuatanArtikel struct {
	Title   string  `json:"title" validate:"required"`
	Content string  `json:"content,omitempty"`
	Version float64 `json:"version" validate:"required"`
	Tags    string  `json:"tags,omitempty"`
}

// GetUserTaskActive PembuatanArtikel godoc
//
//	@Summary		GetUserTaskActive PembuatanArtikel
//	@Description	GetUserTaskActive PembuatanArtikel
//	@Tags			bpmn/PembuatanArtikel
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
//	@Success		200				{string}	string	"PembuatanArtikel GetUserTaskActive"
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pembuatanartikel/search  [get]
func (app *application) getUserTaskActivePembuatanArtikelHandler(w http.ResponseWriter, r *http.Request) {
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

	jsonTemplate := `
{
	"taskDefinitionId": "%s",
	"state:": "%s",
    "pageSize": %d, 
    "sort": [{"field": "%s", "order": "%s"}]
    %s
}`

	var searchAfterStr, searchBeforeStr string

	if taskListQueryParams.SearchAfter != "" {
		searchAfterStr = fmt.Sprintf(`, "searchAfter": %s`, taskListQueryParams.SearchAfter)
	}

	if taskListQueryParams.SearchBefore != "" {
		searchBeforeStr = fmt.Sprintf(`, "searchBefore": %s`, taskListQueryParams.SearchBefore)
	}

	searchParams := searchAfterStr + searchBeforeStr

	body := []byte(fmt.Sprintf(
		jsonTemplate,
		store.PembuatanArtikelID,
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
}
