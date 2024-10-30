package main

import (
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// HealthMonitoring godoc
//
//	@Summary		Fetches health status api
//	@Description	Fetches health status api
//	@Tags			health
//	@Accept			json
//	@produce		json
//	@Success		200	{object}	HealthResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		BasicAuth
//	@Router			/health  [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthResponse{
		Status:  "ok",
		Env:     app.config.env,
		Version: version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
