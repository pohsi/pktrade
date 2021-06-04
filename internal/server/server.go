// Package server provides running server and handle API request
package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/lib/pq"
	"github.com/pohsi/pktrade/internal/auth"
	"github.com/pohsi/pktrade/internal/config"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/internal/healthcheck"
	"github.com/pohsi/pktrade/internal/trade"
	"github.com/pohsi/pktrade/pkg/accesslog"
	"github.com/pohsi/pktrade/pkg/dbconnection"
	"github.com/pohsi/pktrade/pkg/log"
)

// Server is the major be responsible for run and handle rest request
type Server interface {
	Run() error

	Port() int

	Logger() log.Logger
}

type concreteServer struct {
	logger  log.Logger
	config  config.Config
	version string
}

// New creates server instance which takes custom logger
func New(cfg config.Config, logger log.Logger, version string) (Server, error) {

	if logger == nil {
		logger = log.New()
	}

	return &concreteServer{logger: logger, config: cfg}, nil
}

func (c *concreteServer) Run() error {

	db, err := dbx.MustOpen("postgres", c.config.DSN)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	db.QueryLogFunc = logDBQuery(c.logger)
	db.ExecLogFunc = logDBExec(c.logger)
	defer func() {
		if err := db.Close(); err != nil {
			c.logger.Error(err)
		}
	}()

	address := fmt.Sprintf(":%v", c.Port())

	s := &http.Server{
		Addr:    address,
		Handler: buildHandler(c.logger, dbconnection.New(db), &c.config, c.version),
	}

	go routing.GracefulShutdown(s, 10*time.Second, c.logger.Infof)

	c.logger.Infof("server is running at %v", address)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		c.logger.Error(err)
		return err
	}

	return nil
}

func (c *concreteServer) Port() int {
	return c.config.ServerPort
}

func (c *concreteServer) Logger() log.Logger {
	return c.logger
}

func buildHandler(logger log.Logger, db dbconnection.DB, cfg *config.Config, response string) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.NewHandler(logger),
		errors.NewHandler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, response)

	rg := router.Group("/v1")

	authHandler := auth.NewHandler(cfg.JWTSigningKey)

	trade.RegisterHandlers(rg.Group(""),
		trade.NewService(trade.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	return router
}

func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}

// func initDB(db dbconnection.DB) {
// 	// CREATE TABLE `users` (`id` int primary key, `name` varchar(255))
// 	q := db.DB().CreateTable("users", map[string]string{
// 		"id":   "int primary key",
// 		"name": "varchar(255)",
// 	})
// 	err := q.Execute()
// }
