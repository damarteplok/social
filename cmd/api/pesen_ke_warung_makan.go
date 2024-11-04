package main

import (
	"net/http"

	"github.com/damarteplok/social/internal/store"
)

type (
	CreatePesenKeWarungMakanPayload    struct{}
	UpdatePesenKeWarungMakanPayload    struct{}
	DataStorePesenKeWarungMakanWrapper struct {
		Data store.PesenKeWarungMakan `json:"data"`
	}
)

// Create PesenKeWarungMakan godoc
//
//	@Summary		Create PesenKeWarungMakan
//	@Description	Create PesenKeWarungMakan
//	@Tags			camunda
//	@Accept			json
//	@produce		json
//	@Param			payload	body		CreatePesenKeWarungMakanPayload		true	"PesenKeWarungMakan Payload"
//	@Success		201		{object}	DataStorePesenKeWarungMakanWrapper	"PesenKeWarungMakan Created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/camunda/pesen_ke_warung_makan  [post]
func (app *application) createPesenKeWarungMakanHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePesenKeWarungMakanPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// TODO: Change in this code

	// model := &store.PesenKeWarungMakan{}

	// ctx := r.Context()

	// if err := app.store.PesenKeWarungMakan.Create(ctx, model); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	//if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
	//	app.internalServerError(w, r, err)
	//	return
	//}
}
