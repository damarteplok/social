package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict response", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charshet="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
