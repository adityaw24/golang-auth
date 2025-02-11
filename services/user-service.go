package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adityaw24/golang-auth/model"
	"github.com/adityaw24/golang-auth/repository"
	"github.com/adityaw24/golang-auth/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	CreateUser(ctx context.Context, u model.RegisterUser) (string, error)
	GetUserList(ctx context.Context) ([]model.User, error)
	GetUserDetail(ctx context.Context, id uuid.UUID) (model.UserDetailModel, error)
}

type userService struct {
	userRepository repository.UserRepo
	timeoutContext time.Duration
	db             *sqlx.DB
}

func NewUserService(userRepo repository.UserRepo, timeoutContext time.Duration, db *sqlx.DB) UserService {
	return &userService{
		userRepository: userRepo,
		timeoutContext: timeoutContext,
		db:             db,
	}
}

func (service *userService) GetUserList(ctx context.Context) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeoutContext)
	defer cancel()

	var (
		list []model.User
		err  error
	)

	list, err = service.userRepository.GetUserList(ctx)
	if err != nil {
		utils.LogError("Services", "GetUserList", err)
		return list, err
	}
	return list, err
}

func (service *userService) GetUserDetail(ctx context.Context, id uuid.UUID) (model.UserDetailModel, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeoutContext)
	defer cancel()

	var (
		detail model.UserDetailModel
		err    error
	)
	detail, err = service.userRepository.GetUserDetail(ctx, id)
	if err != nil {
		utils.LogError("Services", "GetUserDetail", err)
		return detail, err
	}
	return detail, err
}

func (service *userService) CreateUser(ctx context.Context, u model.RegisterUser) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeoutContext)
	defer cancel()

	var (
		email string
		err   error
	)

	tx, err := service.db.Beginx()
	if err != nil {
		utils.LogError("Services", "CreateUser open tx", err)
		return email, err
	}

	user, err := service.userRepository.FindByEmail(ctx, u.Email)

	if user.Email != "" {
		err = errors.New("email address already registered")
		utils.LogError("Service", "CreateUser", err)
		utils.CommitOrRollback(tx, "Services CreateUser", err)
		return email, err
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.LogError("Service", "CreateUser", err)
		utils.CommitOrRollback(tx, "Services CreateUser", err)
		return email, err
	}

	//Hash plain password
	hashedPassword, err := utils.Hash(u.Password)
	if err != nil {
		utils.LogError("Service", "CreateUser Hash password", err)
		utils.CommitOrRollback(tx, "Services CreateUser hash password", err)
		return email, err
	}

	registeredData := model.User{
		Name:     u.Name,
		Username: u.Username,
		Password: hashedPassword,
		Email:    u.Email,
		Phone:    u.Phone,
	}

	email, err = service.userRepository.CreateUser(ctx, tx, registeredData)
	if err != nil {
		utils.LogError("Service", "CreateUser", err)
		utils.CommitOrRollback(tx, "Services CreateUser", err)
		return email, err
	}
	return email, nil
}
