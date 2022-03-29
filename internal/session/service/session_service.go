package service

import (
	. "github.com/suryaadi44/LoginTest/internal/session/entity"
	. "github.com/suryaadi44/LoginTest/internal/session/repository"
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
