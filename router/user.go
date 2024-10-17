package router

import (
	"github.com/FakJeongTeeNhoi/user-management/controller"
	"github.com/gin-gonic/gin"
)

func UserRouterGroup(server *gin.RouterGroup) {
	user := server.Group("/user")

	// TODO: Add middleware here

	user.GET("/", controller.GetAllUsersHandler)
	user.GET("/:accountId", controller.GetUserHandler)
	user.PUT("/", controller.UpdateUserHandler)
	user.DELETE("/:accountId", controller.DeleteUserHandler)
}
