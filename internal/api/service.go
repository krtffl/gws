package api

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/oxtoacart/bpool"

	getwellsoon "github.com/krtffl/get-well-soon"
	"github.com/krtffl/get-well-soon/internal/config"
	"github.com/krtffl/get-well-soon/internal/http"
	"github.com/krtffl/get-well-soon/internal/http/webui"
	"github.com/krtffl/get-well-soon/internal/logger"
	"github.com/krtffl/get-well-soon/internal/repository"
)

type GWS struct {
	config     *config.Config
	httpServer *http.Server
	errCh      chan error
}

func New(cfg *config.Config) *GWS {
	db, err := NewDatabaseConnection(cfg.Database)
	if err != nil {
		logger.Fatal("[API - New] - "+
			"Failed to connect to database. %v", err)
	}

	repo := repository.NewGWSRepo(db)
	svc := webui.NewSvc(repo)

	// A buffer pool is created to safely check template
	// execution and properly handle the errors
	bpool := bpool.NewBufferPool(64)
	handler := webui.NewHandler(svc, bpool, cfg.Challenge)

	httpServer := http.New(
		cfg.Port,
		handler,
	)
	return &GWS{
		config:     cfg,
		httpServer: httpServer,
		errCh:      make(chan error, 1),
	}
}

func (api *GWS) Run() {
	go func() { api.errCh <- api.httpServer.Run() }()
	func() {
		for err := range api.errCh {
			if err != nil {
				logger.Fatal("[API Service - Run] - Couldn't run API. %v", err)
			}
		}
	}()
}

func (api *GWS) Shutdown() {
	api.httpServer.Shutdown()
}

func NewDatabaseConnection(cfg config.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)

	// creates the connection but does not validate it
	dbConnection, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// validates the connection is ok
	if err := dbConnection.Ping(); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	d, err := iofs.New(getwellsoon.Migrations, "migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err == migrate.ErrNoChange {
		v, _, err := m.Version()
		if err != nil {
			return nil, err
		}

		logger.Info(
			"[API Service - NewDatabaseConnection] - Database already at latest version: %d",
			v,
		)

		return dbConnection, nil
	}
	if err != nil {
		return nil, err
	}

	newV, _, err := m.Version()
	if err != nil {
		return nil, err
	}

	logger.Info("[API Service - NewDatabaseConnection] - Migrated database to version %d", newV)

	return dbConnection, nil
}
