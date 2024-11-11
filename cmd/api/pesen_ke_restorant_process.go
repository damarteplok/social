package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePesenKeRestorantPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type UpdatePesenKeRestorantPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type DataStorePesenKeRestorantWrapper struct {
	Data store.PesenKeRestorant `json:"data"`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create PesenKeRestorant godoc
//
//	@Summary		Create PesenKeRestorant
//	@Description	Create PesenKeRestorant
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreatePesenKeRestorantPayload		true	"PesenKeRestorant Payload"
//	@Success		201		{object}	DataStorePesenKeRestorantWrapper	"PesenKeRestorant Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pesen_ke_restorant  [post]
func (app *application) createPesenKeRestorantHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	var payload CreatePesenKeRestorantPayload
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

	resp, err := app.zeebeClient.StartWorkflow(ctx, store.PesenKeRestorantProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.PesenKeRestorant{
		ProcessDefinitionKey: resp.GetProcessDefinitionKey(),
		Version:              resp.GetVersion(),
		ProcessInstanceKey:   resp.GetProcessInstanceKey(),
		ResourceName:         store.PesenKeRestorantResourceName,
		CreatedBy:            user.ID,
	}

	if err := app.store.PesenKeRestorant.Create(ctx, model); err != nil {
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

// Cancel PesenKeRestorant godoc
//
//	@Summary		Cancel PesenKeRestorant
//	@Description	Cancel PesenKeRestorant
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			processInstanceKey	path		int		true	"ProcessInstanceKey"
//	@Success		200					{string}	string	"PesenKeRestorant Canceled"
//	@Failure		400					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/pesen_ke_restorant/{processInstanceKey}  [delete]
func (app *application) cancelPesenKeRestorantHandler(w http.ResponseWriter, r *http.Request) {
	processInstanceKey, err := strconv.ParseInt(chi.URLParam(r, "processinstanceKey"), 10, 64)
	if err != nil || processInstanceKey < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if err := app.zeebeClient.CancelWorkflow(ctx, processInstanceKey); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "success"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
