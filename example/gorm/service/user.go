package service

import (
	"github.com/arden/gf-plus/example/gorm/model"
	"github.com/arden/gf-plus/example/gorm/repository"
)

// 中间件管理服务
var User = new(service)

type service struct{

}

func (s *service) GetByPhone(phone string) (*model.User, error) {
	user, err := repository.User.GetByPhone(phone)
	return user, err
}

