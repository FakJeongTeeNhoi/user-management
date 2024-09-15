package controller

import (
	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(c *gin.Context) {
	ucr := model.UserCreateRequest{}
	if err := c.ShouldBindJSON(&ucr); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}

	user := ucr.ToUser()
	user, err := user.Create()
	if err != nil {
		response.InternalServerError("Failed to create user").AbortWithError(c)
		return
	}

	c.JSON(201, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"user": user,
	}))
}
