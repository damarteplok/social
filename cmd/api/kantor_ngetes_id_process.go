package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateKantorNgetesIdPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type UpdateKantorNgetesIdPayload struct {
	Variables *map[string]string `json:"variables,omitempty"`
}
type DataStoreKantorNgetesIdWrapper struct {
	Data store.KantorNgetesId `json:"data"`
}

// TODO: U CAN ADD MORE HANDLER LIKE THIS EXAMPLE

// Create KantorNgetesId godoc
//
//	@Summary		Create KantorNgetesId
//	@Description	Create KantorNgetesId
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreateKantorNgetesIdPayload		true	"KantorNgetesId Payload"
//	@Success		201		{object}	DataStoreKantorNgetesIdWrapper	"KantorNgetesId Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/kantor_ngetes_id  [post]
func (app *application) createKantorNgetesIdHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	var payload CreateKantorNgetesIdPayload
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

	resp, err := app.zeebeClient.StartWorkflow(ctx, store.KantorNgetesIdProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.KantorNgetesId{
		ProcessDefinitionKey: resp.GetProcessDefinitionKey(),
		Version:              resp.GetVersion(),
		ProcessInstanceKey:   resp.GetProcessInstanceKey(),
		ResourceName:         store.KantorNgetesIdResourceName,
		CreatedBy:            user.ID,
	}

	if err := app.store.KantorNgetesId.Create(ctx, model); err != nil {
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

// Cancel KantorNgetesId godoc
//
//	@Summary		Cancel KantorNgetesId
//	@Description	Cancel KantorNgetesId
//	@Tags			bpmn
//	@Accept			json
//	@produce		json
//	@Param			processInstanceKey	path		int		true	"ProcessInstanceKey"
//	@Success		200					{string}	string	"KantorNgetesId Canceled"
//	@Failure		400					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bpmn/kantor_ngetes_id/{processInstanceKey}  [delete]
func (app *application) cancelKantorNgetesIdHandler(w http.ResponseWriter, r *http.Request) {
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
