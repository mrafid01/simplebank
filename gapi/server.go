package gapi

import (
	db "github.com/mrafid01/simplebank/db/sqlc"
	"github.com/mrafid01/simplebank/pb"
	"github.com/mrafid01/simplebank/token"
	"github.com/mrafid01/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
