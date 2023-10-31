// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/toomanysource/atreus/app/comment/service/internal/biz"
	"github.com/toomanysource/atreus/app/comment/service/internal/conf"
	"github.com/toomanysource/atreus/app/comment/service/internal/data"
	"github.com/toomanysource/atreus/app/comment/service/internal/server"
	"github.com/toomanysource/atreus/app/comment/service/internal/service"

	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, registry *conf.Registry, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewMysqlConn(confData, logger)
	client := data.NewRedisConn(confData, logger)
	writer := data.NewKafkaWriter(confData, logger)
	dataData, cleanup, err := data.NewData(db, client, writer, logger)
	if err != nil {
		return nil, nil, err
	}
	discovery := server.NewDiscovery(registry)
	userServiceClient := server.NewUserClient(discovery, logger)
	commentRepo := data.NewCommentRepo(dataData, userServiceClient, logger)
	commentUseCase := biz.NewCommentUseCase(commentRepo, logger)
	commentService := service.NewCommentService(commentUseCase, logger)
	grpcServer := server.NewGRPCServer(confServer, commentService, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, commentService, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
