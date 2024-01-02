// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/StartOpsz/pedant/internal/biz"
	"github.com/StartOpsz/pedant/internal/conf"
	"github.com/StartOpsz/pedant/internal/data"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func initApp(confData *conf.Data, pedant *conf.Pedant, llm *conf.Llm, logger *zap.Logger) (*app, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	sessionRepo := data.NewSessionDataSource(dataData)
	localCacheRepo := data.NewLocalCacheDataSource(dataData)
	sessionUseCase := biz.NewSessionUseCase(sessionRepo, localCacheRepo, pedant, llm, logger)
	multiModalRepo := data.NewMultiModalDataSource(dataData)
	multiModalUseCase := biz.NewMultiModalUseCase(multiModalRepo, localCacheRepo, pedant, llm, logger)
	imageRepo := data.NewImageDataSource(dataData)
	imageUseCase := biz.NewImageUseCase(imageRepo, localCacheRepo, pedant, llm, logger)
	mainApp := newApp(sessionUseCase, multiModalUseCase, imageUseCase)
	return mainApp, func() {
		cleanup()
	}, nil
}
