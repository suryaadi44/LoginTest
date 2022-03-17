package service

import (
	. "login/pkg/session/entity"
	. "login/pkg/session/repository"
)

type SessionService struct {
	repository *SessionRepository
}

func (service *SessionService) NewSession(data Session) error {
	return service.repository.NewSession(data)
}

func (service *SessionService) FindSession(uid string) (Session, error) {
	return service.repository.FindSession(uid)
}

func (service *SessionService) DeleteSession(uid string) error {
	return service.repository.DeleteSession(uid)
}

func NewSessionService(repository *SessionRepository) *SessionService {
	return &SessionService{repository: repository}
}
