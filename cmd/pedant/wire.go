//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/StartOpsz/pedant/internal/biz"
	"github.com/StartOpsz/pedant/internal/conf"
	"github.com/StartOpsz/pedant/internal/data"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// initApp init kratos application.

func initApp(*conf.Data, *conf.Pedant, *conf.Llm, *zap.Logger) (*app, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, newApp))
}
