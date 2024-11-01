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
	authentication.GET("/verify",
		middleware.Authorize(),
		middleware.SetAccountInfo(),
		controller.VerifyHandler)

	authentication.POST("/change-password",
		middleware.Authorize(),
		middleware.SetAccountInfo(),
		controller.ChangePasswordHandler)

	authentication.GET("/account-info",
		middleware.SetAccountInfo(),
		controller.GetAccountInfoHandler)
}
