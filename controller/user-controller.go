package controller

import (
	"errors"
	"net/http"

	"github.com/adityaw24/golang-auth/services"
	"github.com/adityaw24/golang-auth/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController interface {
	GetUserList() fiber.Handler
	GetUserDetail() fiber.Handler
}

type userController struct {
	userService services.UserService
}

func NewUserController(userServ services.UserService) UserController {
	return &userController{
		userService: userServ,
	}
}

func (c *userController) GetUserList() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var users, err = c.userService.GetUserList(ctx.Context())
		if err != nil {
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return err
		}
		utils.BuildResponse(ctx, http.StatusOK, "success", users)

		return err
	}
}

func (c *userController) GetUserDetail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		v := ctx.Queries()
		id := uuid.MustParse(v["id"])

		var userDetail, err = c.userService.GetUserDetail(ctx.Context(), id)
		if err != nil {
			errMsg := errors.New("the server cannot find the requested resource").Error()
			utils.BuildErrorResponse(ctx, http.StatusNotFound, errMsg)
			return err
		}
		utils.BuildResponse(ctx, http.StatusOK, "success", userDetail)
		return err
	}
}
