package main

import (
	"context"
	"net/http"
	"time"

	"github.com/damarteplok/social/internal/store"
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

	variables := make(map[string]interface{})
	variables["user"] = user

	if payload.Variables != nil {
		for k, v := range *payload.Variables {
			variables[k] = v
		}
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := app.zeebeClient.StartWorkflow(ctx, store.KantorNgetesIdProcessDefinitionKey, variables)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	model := &store.KantorNgetesId{
		ProcessDefinitionKey: result.GetProcessDefinitionKey(),
		Version:              result.GetVersion(),
		ProcessInstanceKey:   result.GetProcessInstanceKey(),
		ResourceName:         store.KantorNgetesIdResourceName,
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
