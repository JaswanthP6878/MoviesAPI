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
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		// create a channel
		quit := make(chan os.Signal, 1)

		// Notify function relays the corresponding interupts into the
		// given channel (here the channel is "quit")
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// block till we recieve the quit signal.
		s := <-quit

		// once we recieve the quit signal we log the info
		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		//creating a context that ends in 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		app.logger.PrintInfo("Completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		// waiting for the emails to get sent
		app.wg.Wait()
		shutDownError <- nil

	}()

	app.logger.PrintInfo("Starting Server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	//When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS
	//immediately return ErrServerClosed.
	//Make sure the program doesn't exit and waits instead for Shutdown to return.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil

}
