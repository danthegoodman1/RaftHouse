package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/danthegoodman1/RaftHouse/ch"
	"github.com/danthegoodman1/RaftHouse/observability"
	"github.com/danthegoodman1/RaftHouse/utils"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// _ "github.com/danthegoodman1/RaftHouse/ch"
	"github.com/danthegoodman1/RaftHouse/gologger"
	"github.com/danthegoodman1/RaftHouse/http_server"
)

var logger = gologger.NewLogger()

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			logger.Error().Err(err).Msg("error loading .env file, exiting")
			os.Exit(1)
		}
	}
	logger.Debug().Msg("Starting RaftHouse")

	prometheusReporter := observability.NewPrometheusReporter()
	go func() {
		err := observability.StartInternalHTTPServer(":8042", prometheusReporter)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("internal server couldn't start")
			os.Exit(1)
		}
	}()

	httpServer := http_server.StartHTTPServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Warn().Msg("received shutdown signal!")

	// For AWS ALB needing some time to de-register pod
	// Convert the time to seconds
	sleepTime := utils.GetEnvOrDefaultInt("SHUTDOWN_SLEEP_SEC", 0)
	logger.Info().Msg(fmt.Sprintf("sleeping for %ds before exiting", sleepTime))

	time.Sleep(time.Second * time.Duration(sleepTime))
	logger.Info().Msg(fmt.Sprintf("slept for %ds, exiting", sleepTime))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to shutdown HTTP server")
	} else {
		logger.Info().Msg("successfully shutdown HTTP server")
	}

	if ch.KV != nil {
		if err := ch.KV.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close KV")
		}
	}
}
