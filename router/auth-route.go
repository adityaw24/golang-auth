package router

import (
	"github.com/adityaw24/golang-auth/controller"
	"github.com/gofiber/fiber/v2"
)

type AuthRouter interface {
	Login(group fiber.Router, controller controller.AuthController) fiber.Router
	Register(group fiber.Router, controller controller.AuthController) fiber.Router
}

func (r *fiberRouter) Login(group fiber.Router, controller controller.AuthController) fiber.Router {
	return group.Post("/login", controller.Login())
}

func (r *fiberRouter) Register(group fiber.Router, controller controller.AuthController) fiber.Router {
	return group.Post("/register", controller.Register())
}
