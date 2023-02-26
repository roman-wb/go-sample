package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	customConfig "github.com/roman-wb/go-sample/internal/config"
	"github.com/roman-wb/go-sample/internal/flags"
	httpV1 "github.com/roman-wb/go-sample/internal/servers/http/v1"
	"github.com/roman-wb/go-sample/pkg/clients/flipt"
	"github.com/roman-wb/go-sample/pkg/config"
	"github.com/roman-wb/go-sample/pkg/env"
	"github.com/roman-wb/go-sample/pkg/logger"
	"github.com/roman-wb/go-sample/pkg/postgres"
	"github.com/roman-wb/go-sample/pkg/redis"
	"github.com/roman-wb/go-sample/pkg/shutdowner"

	"github.com/rs/zerolog"
)

func main() {
	if err := runApp(); err != nil {
		log.Fatal(err)
	}
}

//nolint:funlen
func runApp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// env
	if err := env.Init(env.EnvFile); err != nil {
		return fmt.Errorf("env init: %w", err)
	}

	// config
	config, err := config.New[customConfig.Config]()
	if err != nil {
		return fmt.Errorf("config new: %w", err)
	}

	// logger
	logger := logger.New(zerolog.InfoLevel)
	// logger.Info().Interface("config", config).Msg("")

	// postgres
	db, dbCleanUp, err := postgres.New(ctx, logger, config.Config.PgURL)
	if err != nil {
		return fmt.Errorf("postgres new: %w", err)
	}
	_ = db //nolint:wsl

	// redis
	rdb, rCleanUp, err := redis.New(ctx, config.Config.RedisURL)
	if err != nil {
		return fmt.Errorf("redis new: %w", err)
	}
	_ = rdb //nolint:wsl

	// flipt - feature flags
	fliptClient, fliptCleanUp, err := flipt.New(
		flipt.WithServer(config.Config.FliptServer),
		flipt.WithInsecure(config.Config.FliptInsecure),
	)
	if err != nil {
		return fmt.Errorf("flipt new: %w", err)
	}

	flags := flags.New(fliptClient, config.Config.FliptTimeout)
	_ = flags

	// OS signal or cancel service context
	signalCtx, signalCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	// http server
	server := httpV1.New(logger, config.Config)
	go startHTTPServer(signalCancel, logger, server)

	<-signalCtx.Done()

	shutdown(
		ctx,
		logger,
		config.Config.ShutdownTimeout,
		server,
		dbCleanUp,
		rCleanUp,
		fliptCleanUp,
	)

	return nil
}

func startHTTPServer(cancel context.CancelFunc, logger *zerolog.Logger, server *http.Server) {
	defer cancel()

	logger.Info().
		Msg("start server")

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error().
			Err(err).
			Msg("listen and serve")
	}
}

func shutdown(
	ctx context.Context,
	logger *zerolog.Logger,
	timeout time.Duration,
	server *http.Server,
	dbCleanUp postgres.CleanUpFunc,
	rCleanUp redis.CleanUpFunc,
	fliptCleanUp flipt.CleanUpFunc,
) {
	logger.Info().
		Str("timeout", timeout.String()).
		Msg("starting shutdown")

	afterHook := func(name string, err error) {
		if err != nil {
			logger.Error().
				Err(err).
				Msgf("%s shutdown with error", name)
		} else {
			logger.Info().
				Msgf("%s shutdown successfully", name)
		}
	}

	shutdowner.NewShutdowner(
		afterHook,
		shutdowner.WithProcess(
			"http server",
			func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		),
	).Wait(ctx, timeout)

	shutdowner.NewShutdowner(
		afterHook,
		shutdowner.WithProcess(
			"postgres",
			func(ctx context.Context) error {
				dbCleanUp()

				return nil
			},
		),
		shutdowner.WithProcess(
			"redis",
			func(ctx context.Context) error {
				return rCleanUp()
			},
		),
		shutdowner.WithProcess(
			"flipt",
			func(ctx context.Context) error {
				return fliptCleanUp()
			},
		),
	).Wait(ctx, timeout)

	logger.Info().
		Msg("shutdown finished")
}
