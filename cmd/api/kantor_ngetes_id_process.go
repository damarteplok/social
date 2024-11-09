package main

import (
    "net/http"
	"github.com/damarteplok/social/internal/store"
)

type CreateKantorNgetesIdPayload struct {
	Variables   		 *map[string]string  `json:"variables,omitempty"`
}
type UpdateKantorNgetesIdPayload struct {
	Variables            *map[string]string  `json:"variables,omitempty"`
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

	//model := &store.KantorNgetesId{}

	//ctx := r.Context()

	// if err := app.store.KantorNgetesId.Create(ctx, model); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	//if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
	//	app.internalServerError(w, r, err)
	//	return
	//}
}
