package router

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type fiberRouter struct {
	app *fiber.App
	// version                 string
	// controllerAuth          controller.AuthController
	// controllerLeaveBalance  controller.LeaveBalanceController
	// controllerLeaveRecord   controller.LeaveRecordController
	// controllerPayrollRecord controller.PayrollRecordController
	// controllerUser          controller.UserController
}

// Router interface
type Router interface {
	App() *fiber.App
	Group(group string) fiber.Router
	Use(mw ...interface{}) fiber.Router
	Run(Port int, serviceName string) error
	UserRouter
	AuthRouter
}

func NewFiberRouter(
	app *fiber.App,
	// version string,
	// controllerAuth controller.AuthController,
	// controllerLeaveBalance controller.LeaveBalanceController,
	// controllerLeaveRecord controller.LeaveRecordController,
	// controllerPayrollRecord controller.PayrollRecordController,
	// controllerUser controller.UserController,
) Router {

	return &fiberRouter{
		app,
		// version,
		// controllerAuth,
		// controllerLeaveBalance,
		// controllerLeaveRecord,
		// controllerPayrollRecord,
		// controllerUser,
	}
}

func (r *fiberRouter) App() *fiber.App {
	return r.app
}

func (r *fiberRouter) Group(group string) fiber.Router {
	return r.app.Group(group)
}

func (r *fiberRouter) Use(mw ...interface{}) fiber.Router {
	return r.app.Use(mw...)
}

func (r *fiberRouter) Run(port int, serviceName string) error {
	log.Printf(serviceName+" - Fiber HTTP Server was running on port %d...\n", port)

	err := r.app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error running Fiber HTTP Server on port %d: %v\n", port, err)
		return err
	}
	return err
}
