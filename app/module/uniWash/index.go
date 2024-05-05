package uniWash

import (
	"go-fiber-starter/app/module/uniWash/repository"
	"go-fiber-starter/app/module/uniWash/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(service.Service),
	fx.Provide(repository.Repository),
)
