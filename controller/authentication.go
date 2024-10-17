package controller

import (
	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	userType := c.Param("type")
	if userType != "staff" && userType != "user" {
		response.BadRequest("Invalid user type").AbortWithError(c)
		return
	}

	lr := model.LoginRequest{}
	if err := c.ShouldBindJSON(&lr); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}

	account, err := service.ValidateCredential(lr)

	if err != nil {
		response.BadRequest("Invalid credential").AbortWithError(c)
		return
	}

	token, err := service.GenerateToken(userType, account)

	if err != nil {
		response.InternalServerError("Cannot generate token").AbortWithError(c)
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, false)
	c.JSON(201, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"token": token,
	}))
}
