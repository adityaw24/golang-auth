package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adityaw24/golang-auth/controller"
	"github.com/adityaw24/golang-auth/databases"
	"github.com/adityaw24/golang-auth/middleware"
	"github.com/adityaw24/golang-auth/repository"
	"github.com/adityaw24/golang-auth/services"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var (
	dbRepoConn databases.DatabaseRepo = databases.NewPostgresRepo()
	db         *sqlx.DB
)

func init() {
	viper.SetConfigFile(`.env`)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	dbHost := viper.GetString(`DB_HOST`)
	dbPort := viper.GetInt(`DB_PORT`)
	dbUser := viper.GetString(`DB_USER`)
	dbPass := viper.GetString(`DB_PASSWORD`)
	dbName := viper.GetString(`DB_NAME`)
	dbMigrateVersion := viper.GetUint(`DB_MIGRATE_VERSION`)
	runMigration := viper.GetBool(`DB_MIGRATE`)
	dbDriver := viper.GetString(`DB_DRIVER`)

	db, err = dbRepoConn.Connect(dbHost, dbPort, dbUser, dbPass, dbName, dbMigrateVersion, runMigration, dbDriver)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()

	appPort := viper.GetInt(`PORT`)
	secretKey := viper.GetString("JWT_SECRET")
	customJwt := services.NewJWTService(secretKey)
	timeoutCtx := time.Duration(viper.GetInt(`TIMEOUT_SECOND`)) * time.Second

	repoUser := repository.NewUserRepo(db)

	serviceAuth := services.NewAuthService(repoUser, timeoutCtx, db)
	serviceUser := services.NewUserService(repoUser, timeoutCtx, db)

	controllerAuth := controller.NewAuthController(serviceAuth, customJwt, serviceUser)
	controllerUser := controller.NewUserController(serviceUser)

	mw := middleware.InitCustomMiddleware(customJwt)

	// httpRouter := router.NewFiberRouter(app)
	// httpRouter.Use(mw.LoggerMiddleware(), mw.CORSMiddleware(), mw.MethodMiddleware(), mw.AuthorizeJWT())

	// version := httpRouter.Group(apiVersion)

	// version.Get("/", func(c *fiber.Ctx) error {
	// 	log.Println(c.OriginalURL())
	// 	return c.SendString("Hello, World!")
	// }).Name("index")

	// httpRouter.UserList(version, controllerUser)
	// httpRouter.UserDetail(version, controllerUser)

	// httpRouter.Login(version, controllerAuth)
	// httpRouter.Register(version, controllerAuth)

	// httpRouter.Run(appPort, "golang-auth")

	Serve(app, appPort, "golang-auth", mw, controllerAuth, controllerUser)

}

func Serve(app *fiber.App, port int, serviceName string, mw middleware.CustomMiddleware, controllerAuth controller.AuthController, controllerUser controller.UserController) error {
	apiVersion := viper.GetString(`API_VERSION`)
	version := app.Group(apiVersion)

	auth := version.Group("/auth")
	auth.Use(mw.LoggerMiddleware(), mw.CORSMiddleware(), mw.MethodMiddleware())
	auth.Post("/login", controllerAuth.Login())
	auth.Post("/register", controllerAuth.Register())

	user := version.Group("/user")
	user.Use(mw.LoggerMiddleware(), mw.CORSMiddleware(), mw.MethodMiddleware(), mw.AuthorizeJWT())
	user.Get("/user-list", controllerUser.GetUserList())
	user.Get("/user-detail/{id:[0-9]+}", controllerUser.GetUserDetail())

	log.Printf(serviceName+" - Fiber HTTP Server was running on port %d...\n", port)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error running Fiber HTTP Server on port %d: %v\n", port, err)
		return err
	}
	return err
}
