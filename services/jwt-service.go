package services

import (
	"errors"
	"time"

	"github.com/adityaw24/golang-auth/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uuid.UUID) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (j *jwtService) GenerateToken(userID uuid.UUID) string {
	claims := &jwtCustomClaim{
		userID.String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    userID.String(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signemodelkenAsString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return signemodelkenAsString
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(parsemodelken *jwt.Token) (interface{}, error) {
		if method, ok := parsemodelken.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.New("invalid signature method")
			utils.LogError("Services", "ValidateToken", err)
			return nil, err
		} else if method != jwt.SigningMethodHS256 {
			err := errors.New("invalid signature method")
			utils.LogError("Services", "ValidateToken", err)
			return nil, err
		} else {
			return []byte(j.secretKey), nil
		}
	})

	if err != nil {
		err := errors.New("token invalid")
		utils.LogError("Services", "ValidateToken", err)
		return nil, err
	}

	return token, nil
}
