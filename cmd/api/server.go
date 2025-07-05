package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("caught signal", map[string]string{
			"signal": s.String(),
		})

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		cancel()
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.bot.Run(ctx)

	}()

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.eventWorker.Start(ctx)
	}()

	app.logger.Info("starting server", "details", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.Env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "details", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
