package service

import (
	"github.com/arden/easy/example/gorm/model"
	"github.com/arden/easy/example/gorm/repository"
)

// 中间件管理服务
var User = new(service)

type service struct{

}

func (s *service) GetByPhone(phone string) (*model.User, error) {
	user, err := repository.User.GetByPhone(phone)
	return user, err
}

