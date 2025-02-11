package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/adityaw24/golang-auth/model"
	"github.com/adityaw24/golang-auth/repository"
	"github.com/adityaw24/golang-auth/utils"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredentials(ctx context.Context, d model.UserLogin) (model.UserResponse, error)
}

type authService struct {
	userRepo       repository.UserRepo
	timeoutContext time.Duration
	db             *sqlx.DB
}

func NewAuthService(userRepo repository.UserRepo, timeoutContext time.Duration, db *sqlx.DB) AuthService {
	return &authService{
		userRepo:       userRepo,
		timeoutContext: timeoutContext,
		db:             db,
	}
}

func (service *authService) VerifyCredentials(ctx context.Context, d model.UserLogin) (model.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeoutContext)
	defer cancel()

	var (
		user model.UserResponse
		err  error
	)

	user, err = service.userRepo.FindByEmail(ctx, d.Email)
	if err != nil {
		utils.LogError("Services", "VerifyCredentials", err)
		return user, err
	}

	//Compare password
	isValid := comparePassword(user.Password, d.Password)
	if !isValid {
		err = errors.New("failed to login. check your credential")
		utils.LogError("Services", "VerifyCredentials", err)
		return user, err
	}
	//Return login user data
	return user, nil
}

// Compare plain hashed password retrieved from db against user-entered password
func comparePassword(hashedPass string, plainPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		utils.LogError("Services", "VerifyCredentials", err)
		return false
	}
	log.Println("| Password Matched.")
	return true
}
