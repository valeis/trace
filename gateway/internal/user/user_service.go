package user

import "Booking_system/gateway/model"

type IUserRepository interface {
	Delete(email string) (string, error)
	GetByEmail(email string) (model.User, error)
	GetById(id string) (model.User, error)
	UpdateUser(updateUser model.User) (string, error)
	GetAllUsers(page uint32, limit uint32) ([]model.User, string, error)
	AddAdmin(email string) (string, error)
	SetUserRole(email string) (string, error)
	CreateUser(user model.User) (string, error)
}

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Delete(email string) (string, error) {
	return svc.repo.Delete(email)
}

func (svc *UserService) GetByEmail(email string) (model.User, error) {
	return svc.repo.GetByEmail(email)
}

func (svc *UserService) GetById(id string) (model.User, error) {
	return svc.repo.GetById(id)
}

func (svc *UserService) UpdateUser(updateUser model.User) (string, error) {
	return svc.repo.UpdateUser(updateUser)
}

func (svc *UserService) GetAllUsers(page uint32, limit uint32) ([]model.User, string, error) {
	return svc.repo.GetAllUsers(page, limit)
}

func (svc *UserService) AddAdmin(email string) (string, error) {
	return svc.repo.AddAdmin(email)
}

func (svc *UserService) SetUserRole(email string) (string, error) {
	return svc.repo.SetUserRole(email)
}

func (svc *UserService) CreateUser(user model.User) (string, error) {
	return svc.repo.CreateUser(user)
}
