package main

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"

	"github.com/felixge/httpsnoop"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				app.logger.Error("panic recovered", "error", err, "stack", string(buf))
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.Limiter.Enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.RPS), app.config.Limiter.Burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		w.Header().Add("Vary", "Cookie")

		if auth := r.Header.Get("Authorization"); auth != "" {
			app.logger.Info("bearer auth")
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			ctx := r.Context()
			token := oauth2.Token{AccessToken: parts[1]}
			client := app.oauth2Config.Client(ctx, &token) //app.models.Users.GetForToken(data.ScopeAuthentication, parts[1])
			if client == nil {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}

			resp, err := client.Get("https://discord.com/api/users/@me")
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()

			}(resp.Body)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			var userInfo UserInfo
			if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			r = app.contextSetUser(r, nil, &userInfo)
		} else {
			app.logger.Info("cookie auth")
			du, err := app.getCookieHandler(r)
			if err == nil {
				app.logger.Info("cookie authenticated", "user", du)
				r = app.contextSetUser(r, nil, &du)
			} else if errors.Is(err, http.ErrNoCookie) {
				app.logger.Info("anonymous user", "user", du)
				r = app.contextSetUser(r, data.AnonymousUser, nil)
			} else {
				app.logger.Info("invalid cookie", "error", err)
				app.InvalidCookieResponse(w, r, err)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, du := app.contextGetUser(r)
		if (user == nil && du == nil) || user.IsAnonymous() {
			app.logger.Warn("user not authenticated", "user", user, "du", du)
			app.authenticationRequiredResponse(w, r)
			return
		}
		next(w, r)
	}
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := app.contextGetUser(r)
		perms, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !perms.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next(w, r)
	}
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.Cors.TrustedOrigins {
				if origin == app.config.Cors.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

						w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						w.Header().Set("Access-Control-Allow-Credentials", "true")

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) setTracingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "X-Trace-ID", traceID)
		w.Header().Set("X-Trace-ID", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
