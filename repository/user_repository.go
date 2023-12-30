package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

// DI
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// emailが紐づくユーザー情報が存在する場合は書き換える
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	// ユーザー情報の作成に成功した場合はuserオブジェクトを上書き、エラーの場合はエラー情報を出力
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}