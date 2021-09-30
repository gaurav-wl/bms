package server

import (
	dbPackage "github.com/gauravcoco/bms/db"
	"github.com/gauravcoco/bms/providers"
	"github.com/gauravcoco/bms/providers/configProvider"
	"github.com/gauravcoco/bms/providers/converter"
	"github.com/gauravcoco/bms/providers/dbHelpProvider"
	"github.com/gauravcoco/bms/providers/dbProvider"
	"github.com/gauravcoco/bms/providers/keyProvider"
	"github.com/gauravcoco/bms/providers/middlewareProvider"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	PSQL        providers.PSQLProvider
	Config      providers.ConfigProvider
	Middlewares providers.MiddlewareProvider
	KeyProvider providers.KeyProvider
	DBHelper    providers.DBHelpProvider
	Converter   providers.Converter
}

func Run() {
	serverInit().Start()
}

func serverInit() *Server {
	var config providers.ConfigProvider

	// Loading Config
	config = configProvider.NewConfigProvider()
	err := config.Read()
	if err != nil {
		logrus.Fatalf("Error reading config %v", err)
	}

	// PSQL connection
	db := dbProvider.NewPSQLProvider(config.GetPSQLConnectionString(), config.GetPSQLMaxConnection(), config.GetPSQLMaxIdleConnection())


	// Migrations Up
	dbPackage.NewMigrationProvider(db.DB()).Up()

	return &Server{
		PSQL:        db,
		Config:      config,
		Middlewares: middlewareProvider.NewMiddleware(db.DB(), config.GetJWTKey()),
		KeyProvider: keyProvider.NewKeyProvider(),
		DBHelper:    dbHelpProvider.NewDBHelper(db.DB()),
		Converter:   converter.NewConverter(),
	}
}

func (srv *Server) Start() {
	addr := ":" + srv.Config.GetServerPort()

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.InjectRoutes(),
		ReadTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Minute,
	}

	logrus.Info("Server running at PORT ", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
		return
	}
}
