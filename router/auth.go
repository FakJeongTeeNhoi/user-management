package router

import (
	"github.com/FakJeongTeeNhoi/user-management/controller"
	"github.com/FakJeongTeeNhoi/user-management/middleware"
	"github.com/gin-gonic/gin"
)

func AuthenticationGroup(server *gin.RouterGroup) {
	authentication := server.Group("/auth")

	// TODO: Add middleware here

	authentication.POST("/login", controller.LoginHandler)
	authentication.POST("/logout", controller.LogoutHandler)
	authentication.POST("/register", controller.RegisterHandler)
	authentication.GET("/verify", controller.VerifyHandler).
		Use(middleware.Authorize()).
		Use(middleware.SetAccountInfo())

	authentication.POST("change-password", controller.ChangePasswordHandler).
		Use(middleware.Authorize()).
		Use(middleware.SetAccountInfo())
}
