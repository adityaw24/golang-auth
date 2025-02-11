package controller

import (
	"errors"
	"net/http"

	"github.com/adityaw24/golang-auth/model"
	"github.com/adityaw24/golang-auth/services"
	"github.com/adityaw24/golang-auth/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login() fiber.Handler
	Register() fiber.Handler
}

type authController struct {
	authServ services.AuthService
	jwtServ  services.JWTService
	userServ services.UserService
}

func NewAuthController(authServ services.AuthService, jwtServ services.JWTService, userServ services.UserService) AuthController {
	return &authController{
		authServ: authServ,
		jwtServ:  jwtServ,
		userServ: userServ,
	}
}

func (c *authController) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userLoginData := model.UserResponse{}
		userLogin := model.UserLogin{}
		err := ctx.BodyParser(&userLogin)

		if err != nil {
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return err
		}

		userLoginData, err = c.authServ.VerifyCredentials(ctx.Context(), userLogin)
		if err != nil {
			errMsg := errors.New("incorrect email/password").Error()
			utils.BuildErrorResponse(ctx, http.StatusUnauthorized, errMsg)
			return err
		}
		token := c.jwtServ.GenerateToken(userLoginData.User_id)
		userLoginData.Token = token

		utils.BuildResponse(ctx, http.StatusOK, "success", userLoginData)
		return err
	}
}

func (c *authController) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var regUser model.RegisterUser
		var createdUser string

		err := ctx.BodyParser(&regUser)
		if err != nil {
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return err
		}

		createdUser, err = c.userServ.CreateUser(ctx.Context(), regUser)
		if err != nil {
			errMsg := err.Error()
			utils.BuildErrorResponse(ctx, http.StatusConflict, errMsg)
			return err
		}
		utils.BuildResponse(ctx, http.StatusOK, "success", createdUser)
		return err
	}
}
