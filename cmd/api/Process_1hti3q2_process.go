package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateProcess1hti3q2Payload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type UpdateProcess1hti3q2Payload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type DataStoreProcess1hti3q2Wrapper struct {
	Data store.Process1hti3q2 `json:"data"`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create Process1hti3q2 godoc
//
//	@Summary		Create Process1hti3q2
//	@Description	Create Process1hti3q2
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreateProcess1hti3q2Payload		true	"Process1hti3q2 Payload"
//	@Success		201		{object}	DataStoreProcess1hti3q2Wrapper	"Process1hti3q2 Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/Process_1hti3q2  [post]
func (app *application) createProcess1hti3q2Handler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload CreateProcess1hti3q2Payload
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
	variables["user"] = user

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if payload.Variables != nil {
		for k, v := range *payload.Variables {
			variables[k] = v
		}
	}

	resp, err := app.zeebeClient.StartWorkflow(ctx, store.Process1hti3q2ProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.Process1hti3q2{
		ProcessDefinitionKey: resp.GetProcessDefinitionKey(),
		Version:              resp.GetVersion(),
		ProcessInstanceKey:   resp.GetProcessInstanceKey(),
		ResourceName:         store.Process1hti3q2ResourceName,
		CreatedBy:            user.ID,
	}

	if err := app.store.Process1hti3q2.Create(ctx, model); err != nil {
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

// Cancel Process1hti3q2 godoc
//
//	@Summary		Cancel Process1hti3q2
//	@Description	Cancel Process1hti3q2
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"id"
//	@Success		200	{string}	string	"Process1hti3q2 Canceled"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/Process_1hti3q2/{id}  [delete]
func (app *application) cancelProcess1hti3q2Handler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.store.Process1hti3q2.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// delete model
	if err := app.store.Process1hti3q2.Delete(ctx, model.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: add rollback if failed to cancel in zeebe
	if err := app.zeebeClient.CancelWorkflow(ctx, model.ProcessInstanceKey); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetById Process1hti3q2 godoc
//
//	@Summary		GetById Process1hti3q2
//	@Description	GetById Process1hti3q2
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			id	path		int		true	"ID from table"
//	@Success		200	{string}	string	"Process1hti3q2 GetById"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/Process_1hti3q2/{id}  [get]
func (app *application) getByIdProcess1hti3q2Handler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	model, err := app.store.Process1hti3q2.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, model); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
