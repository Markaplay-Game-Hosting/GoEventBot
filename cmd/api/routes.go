package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.requireAuthenticatedUser(app.updateUserPasswordHandler))

	router.HandlerFunc(http.MethodGet, "/oauth/authenticate", app.authenticateHandler)
	router.HandlerFunc(http.MethodPost, "/oauth/callback", app.callbackHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Events routes
	router.HandlerFunc(http.MethodPost, "/v1/events", app.requireAuthenticatedUser(app.createEventHandler))
	router.HandlerFunc(http.MethodGet, "/v1/events/:id", app.getEventHandler)
	router.HandlerFunc(http.MethodGet, "/v1/events", app.requireAuthenticatedUser(app.getAllEventsHandler))
	router.HandlerFunc(http.MethodPut, "/v1/events/:id", app.requireAuthenticatedUser(app.updateEventHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/events/:id", app.requireAuthenticatedUser(app.deleteEventHandler))

	// Tags routes
	router.HandlerFunc(http.MethodPost, "/v1/tags", app.requireAuthenticatedUser(app.createTagHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tags/:id", app.getTagHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tags", app.listTagsHandler)
	router.HandlerFunc(http.MethodPut, "/v1/tags/:id", app.requireAuthenticatedUser(app.updateTagHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tags/:id", app.requireAuthenticatedUser(app.deleteTagHandler))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.setTracingId(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))))
}
