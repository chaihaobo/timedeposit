// Package service
// @author： Boice
// @createTime：2022/7/26 11:25
package service

import (
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service/mambu"
)

type Service struct {
	Mambu *mambu.Mambu
}

func NewService(common *common.Common, repo *repository.Repository) *Service {
	return &Service{
		Mambu: mambu.NewMambuService(common, repo),
	}
}
