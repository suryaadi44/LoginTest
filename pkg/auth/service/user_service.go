package service

import (
	. "login/pkg/auth/entity"
	. "login/pkg/auth/repository"
)

type UserService struct {
	repository *UserRepository
}

func (service *UserService) NewUser(data User) error {
	return service.repository.NewUser(data)
}

func (service *UserService) FindUser(user string) (User, error) {
	return service.repository.FindUser(user)
}

func NewUserService(repository *UserRepository) *UserService {
	return &UserService{repository: repository}
}