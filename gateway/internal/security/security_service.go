package security

import "Booking_system/gateway/model"

type ISecurityRepository interface {
	Login(user model.UserCredentials) (model.Tokens, error)
	Refresh(token string) (model.Tokens, error)
	Create(createUser model.UserCredentials) (string, error)
	CheckUserExistence(username model.Username) (string, error)
}

type SecurityService struct {
	repo ISecurityRepository
}

func NewSecurityService(repo ISecurityRepository) *SecurityService {
	return &SecurityService{repo: repo}
}

func (svc *SecurityService) Login(loginUser model.UserCredentials) (model.Tokens, error) {
	return svc.repo.Login(loginUser)
}

func (svc *SecurityService) Refresh(token string) (model.Tokens, error) {
	return svc.repo.Refresh(token)
}

func (svc *SecurityService) Create(createUser model.UserCredentials) (string, error) {
	return svc.repo.Create(createUser)
}

func (svc *SecurityService) CheckUserExistence(username model.Username) (string, error) {
	return svc.repo.CheckUserExistence(username)
}
