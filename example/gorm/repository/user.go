package repository

import (
	"errors"
	"github.com/arden/easy"
	"github.com/arden/easy/example/gorm/model"
	"github.com/arden/easy/gorm/repository"
	"github.com/gogf/gf/os/glog"
	"gorm.io/gorm"
)

var (
	User = newUserRepository()
)

type UserRepository interface {
	Insert(user *model.User) error
	GetByPhone(phone string) (*model.User, error)
}

type userRepository struct {
	repository.TransactionRepository
}

func newUserRepository(name...string) UserRepository {
	db := easy.Gorm(name...)
	userRepository := &userRepository{
		repository.NewGormRepository(db, glog.DefaultLogger(), ""),
	}
	userRepository.AutoMigrateOrWarn(&model.User{})
	return userRepository
}

// Insert insert a new user
func (r userRepository) Insert(user *model.User) error {
	return r.Create(user)
}

// GetByEmail get a user by email
func (r userRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	result := r.DB().First(&user, model.User{Phone: phone})
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
