package user

import (
	client2 "github.com/codingXiang/go-harbor-client/client"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

const (
	repo     = "api.user."
	root     = repo + "root"
	base     = repo + "base"
	current  = repo + "current"
	pwd      = repo + "password"
	sysadmin = repo + "sysadmin"
)

type Service interface {
	List() ([]User, *gorequest.Response, []error)
	Get(id int) (User, *gorequest.Response, []error)
	Create(user *User) (*gorequest.Response, []error)
	Update(id int, user *User) (*gorequest.Response, []error)
	Delete(id int) (*gorequest.Response, []error)
	Current() (User, *gorequest.Response, []error)
	ChangeSysadmin(id int, role UpdateRole) (*gorequest.Response, []error)
	ChangePassword(id int, password UpdatePassword) (*gorequest.Response, []error)
}

type UserService struct {
	client client2.ClientInterface
}

func NewUserService(client client2.ClientInterface) Service {
	return &UserService{client: client}
}

func (s *UserService) getConfigString(key string) string {
	return s.client.GetConfig().GetString(key)
}

func (s *UserService) List() ([]User, *gorequest.Response, []error) {
	var v []User
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, s.getConfigString(root)).
		EndStruct(&v)
	return v, &resp, errs
}

func (s *UserService) Get(id int) (User, *gorequest.Response, []error) {
	var v User
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(base), id)).
		EndStruct(&v)
	return v, &resp, errs
}

func (s *UserService) Create(user *User) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, s.getConfigString(root)).
		Send(*user).
		End()
	return &resp, errs
}
func (s *UserService) Update(id int, user *User) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf(s.getConfigString(base), id)).
		Send(*user).
		End()
	return &resp, errs
}
func (s *UserService) Delete(id int) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(base), id)).
		End()
	return &resp, errs
}
func (s *UserService) Current() (User, *gorequest.Response, []error) {
	var v User
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, s.getConfigString(current)).
		EndStruct(&v)
	return v, &resp, errs
}

func (s *UserService) ChangePassword(id int, password UpdatePassword) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(pwd), id)).
		Send(password).
		End()
	return &resp, errs
}

func (s *UserService) ChangeSysadmin(id int, role UpdateRole) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(sysadmin), id)).
		Send(role).
		End()
	return &resp, errs
}
