package user_service

import (
	. "login/pkg/user/entity"
	. "login/pkg/user/repository"
)

type UserService struct {
	repository *UserRepository
}

func (service *UserService) NewUser(data User) {
	service.repository.NewUser(data)
}

func (service *UserService) FindUser(user string) User {
	return service.repository.FindUser(user)
}

func NewUserService(repository *UserRepository) *UserService {
	return &UserService{repository: repository}
}
