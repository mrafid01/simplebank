package gapi

import (
	db "github.com/mrafid01/simplebank/db/sqlc"
	"github.com/mrafid01/simplebank/pb"
	"github.com/mrafid01/simplebank/token"
	"github.com/mrafid01/simplebank/util"
	"github.com/mrafid01/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
