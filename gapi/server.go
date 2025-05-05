package gapi

import (
	"fmt"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/token"
	"github.com/AnkitNayan83/houseBank/util"
)

type Server struct {
	pb.UnimplementedHouseBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
