package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/AnkitNayan83/houseBank/api"
	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/gapi"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(store, config)
	runGRPCServer(store, config)

}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal("cannot create http server: ", err)
	}

	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal("cannot start http server: ", err)
	}
}

func runGRPCServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHouseBankServer(grpcServer, server)
	reflection.Register(grpcServer) // this will allow the client to see the available grpc services and how to call them

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("cannot create gRPC listener: ", err)
	}

	log.Printf("starting gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}

func runGatewayServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
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
		log.Fatal("cannot register grpc gateway server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)

	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("starting HTTP Gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server: ", err)
	}
}
