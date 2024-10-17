package router

import (
	"github.com/FakJeongTeeNhoi/user-management/controller"

	"github.com/gin-gonic/gin"
)

func AuthenticationGroup(server *gin.RouterGroup) {
	authentication := server.Group("/authen")
	{
		authentication.POST("/login", controller.LoginHandler)
	}
}
