package api

import (
	v1 "two_services/api/api/handlers/v1"
	"two_services/api/api/middleware"
	"two_services/api/api/token"
	"two_services/api/config"
	"two_services/api/pkg/logger"
	"two_services/api/services"
	"two_services/api/storage/repo"

	_ "two_services/api/api/docs"

	"github.com/casbin/casbin/v2"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf            config.Config
	Logger          logger.Logger
	ServiceManager  services.IServiceManager
	InMemoryStorage repo.InMemoryStorageI
	CasbinEnforcer  *casbin.Enforcer
}

// @BasePath /v1
// @version 1.0
// @description this is a user and task services api
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization

func New(option Option) *gin.Engine {
	router := gin.New()
	jwtHandler := token.JWTHandler{
		SigninKey: option.Conf.SigninKey,
		Log:       option.Logger,
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.NewAuthorizer(option.CasbinEnforcer, jwtHandler, option.Conf))

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		JwtHandler: jwtHandler,
		
	})

	task := router.Group("/v1")
	task.POST("/tasks", handlerV1.CreateTask)
	task.GET("/tasks/:id", handlerV1.GetTask)
	task.PUT("/tasks", handlerV1.UpdateTask)
	task.DELETE("/tasks/:id", handlerV1.DeleteTask)
	task.GET("/tasks", handlerV1.ListTasks)
	task.GET("/tasksoverdue", handlerV1.ListTasksOverdue)

	task.GET("/users/:id", handlerV1.GetUser)
	task.DELETE("/users/:id", handlerV1.DeleteUser)
	task.PUT("/users", handlerV1.UpdateUser)
	task.POST("/users/filtr", handlerV1.Filtr)

	task.POST("/users", handlerV1.RegistrUser)
	task.POST("/users/login", handlerV1.Login)
	task.POST("/users/verify", handlerV1.Verify)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
	return router
}
