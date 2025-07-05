package main

import (
	"expvar"
	"github.com/Markaplay-Game-Hosting/GoEventBot/ui"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	// ─── Static & Frontend ────────────────────────────────────────────────────
	if IsProduction() {
		app.logger.Info("Running in production mode")
		fs := http.FileServer(http.FS(ui.DistDirFS))

		router.Handler("GET", "/assets/*filepath", fs)
		router.Handler("GET", "/favicon.svg", fs)
		router.Handler("GET", "/build/*filepath", fs)

		// SPA entrypoint + fallback
		router.HandlerFunc("GET", "/", app.spaHandler)
		router.NotFound = http.HandlerFunc(app.spaHandler)
	} else {
		app.logger.Info("Running in development mode")
		router.NotFound = http.HandlerFunc(app.ServeFrontendProxy)
	}

	// ─── Public JSON Endpoints ────────────────────────────────────────────────
	router.HandlerFunc("GET", "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc("GET", "/oauth/authenticate", app.authenticateHandler)
	router.HandlerFunc("GET", "/oauth/callback", app.callbackHandler)

	// ─── /docs and /reference and /debug ──────────────────────────────────────
	router.ServeFiles("/docs/*filepath", http.Dir("./docs"))
	router.HandlerFunc("GET", "/reference", app.referenceHandler)
	router.Handler("GET", "/debug/vars", expvar.Handler())

	// ─── V1 API (all routes under /v1) ────────────────────────────────────────
	// Events
	router.HandlerFunc("POST", "/v1/events", app.requireAuthenticatedUser(app.createEventHandler))
	router.HandlerFunc("GET", "/v1/events/:id", app.requireAuthenticatedUser(app.getEventHandler))
	router.HandlerFunc("GET", "/v1/events", app.requireAuthenticatedUser(app.getAllEventsHandler))
	router.HandlerFunc("GET", "/v1/calendar", app.requireAuthenticatedUser(app.getAllEventInstancesHandler))
	router.HandlerFunc("PUT", "/v1/events/:id", app.requireAuthenticatedUser(app.updateEventHandler))
	router.HandlerFunc("DELETE", "/v1/events/:id", app.requireAuthenticatedUser(app.deleteEventHandler))

	// Tags
	router.HandlerFunc("POST", "/v1/tags", app.requireAuthenticatedUser(app.createTagHandler))
	router.HandlerFunc("GET", "/v1/tags/:id", app.requireAuthenticatedUser(app.getTagHandler))
	router.HandlerFunc("GET", "/v1/tags", app.requireAuthenticatedUser(app.listTagsHandler))
	router.HandlerFunc("PUT", "/v1/tags/:id", app.requireAuthenticatedUser(app.updateTagHandler))
	router.HandlerFunc("DELETE", "/v1/tags/:id", app.requireAuthenticatedUser(app.deleteTagHandler))

	// Discord interactions
	router.HandlerFunc("GET", "/v1/discord/:id/roles", app.requireAuthenticatedUser(app.listRolesHandler))
	router.HandlerFunc("GET", "/v1/discord/:id/channels", app.requireAuthenticatedUser(app.listChannelsHandler))
	router.HandlerFunc("GET", "/v1/discord/:id", app.requireAuthenticatedUser(app.getGuildInfo))

	// Jobs
	router.HandlerFunc("GET", "/v1/jobs", app.requireAuthenticatedUser(app.listJobs))
	router.HandlerFunc("GET", "/v1/jobs/:id", app.requireAuthenticatedUser(app.getJob))
	router.HandlerFunc("PATCH", "/v1/jobs/:id/cancel", app.requireAuthenticatedUser(app.skipJob))

	// Settings
	router.HandlerFunc("GET", "/v1/settings", app.requireAuthenticatedUser(app.listSettings))
	router.HandlerFunc("POST", "/v1/settings", app.requireAuthenticatedUser(app.addSettings))
	router.HandlerFunc("GET", "/v1/settings/:id", app.requireAuthenticatedUser(app.getSettings))
	router.HandlerFunc("PUT", "/v1/settings/:id", app.requireAuthenticatedUser(app.updateSetting))
	router.HandlerFunc("DELETE", "/v1/settings/:id", app.requireAuthenticatedUser(app.deleteSetting))

	// ─── Global middleware chain ───────────────────────────────────────────────
	// Order matters: Recover panic → Trace ID → CORS → Rate Limit → Metrics → Auth
	handler := app.recoverPanic(
		app.setTracingID(
			app.enableCORS(
				app.rateLimit(
					app.metrics(
						app.authenticate(router),
					)))))

	return handler
}
