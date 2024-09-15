package router

import (
	"github.com/FakJeongTeeNhoi/user-management/controller"
	"github.com/gin-gonic/gin"
)

func UserRouterGroup(server *gin.RouterGroup) {
	user := server.Group("/user")
	{
		user.POST("/register", controller.RegisterUserHandler)
	}
}
