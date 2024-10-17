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

func RegisterHandler(c *gin.Context) {
	userType := c.Param("type")
	if userType != "staff" && userType != "user" {
		response.BadRequest("Invalid user type").AbortWithError(c)
		return
	}

	if userType == "staff" {
		scr := model.StaffCreateRequest{}
		if err := c.ShouldBindJSON(&scr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		// random 6 letter password
		password := service.RandomString(8)
		scr.Password = password

		_, err := createStaffHandler(scr)
		if err != nil {
			response.InternalServerError("Failed to create staff").AbortWithError(c)
			return
		}

		if err := service.SendMail(
			scr.Email,
			"Staff Registration",
			"Your account has been created. Your password is "+password,
		); err != nil {
			response.InternalServerError("Failed to send email").AbortWithError(c)
			return
		}

		c.JSON(201, response.CommonResponse{
			Success: true,
		})
	} else {
		ucr := model.UserCreateRequest{}
		if err := c.ShouldBindJSON(&ucr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		password := service.RandomString(8)
		ucr.Password = password

		_, err := registerUserHandler(ucr)
		if err != nil {
			response.InternalServerError("Failed to create user").AbortWithError(c)
			return
		}

		if err := service.SendMail(
			ucr.Email,
			"User Registration",
			"Your account has been created. Your password is "+password,
		); err != nil {
			response.InternalServerError("Failed to send email").AbortWithError(c)
			return
		}

		c.JSON(201, response.CommonResponse{
			Success: true,
		})
	}
}
