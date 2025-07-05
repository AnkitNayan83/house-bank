package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/AnkitNayan83/houseBank/api"
	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/gapi"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/util"
	"github.com/AnkitNayan83/houseBank/workers"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/AnkitNayan83/houseBank/doc/statik" // this is the generated statik file
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db:")
	}

	// Run db migrations
	runDbMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := workers.NewRedisTaskDistributor(&redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGinServer(store, config)
	go runGatewayServer(store, config, taskDistributor)
	runGRPCServer(store, config, taskDistributor)

}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskPorcessor := workers.NewRedisTaskPorcessor(&redisOpt, store)
	log.Info().Msg("starting task processor ⌛⌛")
	err := taskPorcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failde to start task processor ❌❌")
	}

	log.Info().Msg("task processor started successfully ✅✅")
}

func runDbMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migrate up:")
	}

	log.Info().Msg("db migrated successfully")
}

func runGinServer(store db.Store, config util.Config) {
	// set gin mode
	gin.SetMode(gin.ReleaseMode)

	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gin http server:")
	}

	log.Info().Msg("starting gin server at " + config.GinServerAddress)

	err = server.Start(config.GinServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gin http server:")
	}
}

func runGRPCServer(store db.Store, config util.Config, taskDistributor workers.TaskDistributor) {
	server, err := gapi.NewServer(store, config, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterHouseBankServer(grpcServer, server)
	reflection.Register(grpcServer) // this will allow the client to see the available grpc services and how to call them

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC listener:")
	}

	log.Info().Msg("starting gRPC server at " + listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server:")
	}
}

func runGatewayServer(store db.Store, config util.Config, taskDistributor workers.TaskDistributor) {
	server, err := gapi.NewServer(store, config, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc server:")
	}

	jsonOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpt)

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	err = pb.RegisterHouseBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register grpc gateway server: ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.GRPCGatewayAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener:")
	}

	log.Info().Msg("starting HTTP Gateway server at " + listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start http gateway server:")
	}
}
