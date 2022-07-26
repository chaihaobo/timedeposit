// Package application
// @author： Boice
// @createTime：2022/7/22 16:19
package application

import (
	"gitlab.com/bns-engineering/td/application/engine"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service"
)

type Application struct {
	Repo    *repository.Repository
	Service *service.Service
	Engine  engine.Engine
}

func NewApplication(common *common.Common, repository *repository.Repository, service *service.Service) *Application {
	return &Application{
		Repo:    repository,
		Service: service,
		Engine:  engine.NewEngine(common, repository, service),
	}

}
