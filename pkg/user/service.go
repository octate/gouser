package user

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Service struct {
	conf *viper.Viper
	log  *logrus.Logger
	Repo Repository
}

// NewService returns a user service object.
func NewService(conf *viper.Viper, log *logrus.Logger, Repo Repository) *Service {
	return &Service{conf: conf, log: log, Repo: Repo}
}

func (s Service) CreateUser(ctx context.Context, user *User) (err error) {
	return s.Repo.CreateUser(ctx, user)
}

func (s Service) FetchUserByID(ctx context.Context, userID int) (user *User, err error) {

	return s.Repo.Fetch(ctx, userID)
}

func (s Service) UpdateUser(ctx context.Context, user *User) (err error) {
	return s.Repo.UpdateUser(ctx, user)
}

func (s *Service) FetchAllUsers(ctx context.Context, filter *UserRequest) (users []User, pagination Pagination, err error) {

	return s.Repo.FetchAllUsers(ctx, filter)
}

func (s *Service) FetchByMobileNumber(dCtx context.Context, mobile string) (user *User, err error) {

	return s.Repo.FetchByMobileNumber(dCtx, mobile)
}
