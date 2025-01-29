package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Server Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad Request Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Not Found Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Conflict Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Unauthorized Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	// When a browser sees the WWW-Authenticate header in a 401 Unauthorized response, it automatically shows a login popup to the user.
	// This popup allows the user to enter their credentials (username and password)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) rateLimitExceededError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJsonError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

// Warnings
func (app *application) forbidenWarning(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden", "method", r.Method, "path", r.URL.Path, "error")
	writeJsonError(w, http.StatusForbidden, "forbidden")
}
