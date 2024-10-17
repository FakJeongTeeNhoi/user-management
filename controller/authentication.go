package controller

import (
	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
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

func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, false)
	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}

func VerifyHandler(c *gin.Context) {
	info, err := c.Get("user")
	if !err {
		response.InternalServerError("Cannot get account info").AbortWithError(c)
		return
	}

	accountInfo := info.(jwt.MapClaims)
	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(accountInfo))
}

func RegisterHandler(c *gin.Context) {
	userType := c.Param("type")
	if userType != "staff" && userType != "user" {
		response.BadRequest("Invalid user type").AbortWithError(c)
		return
	}

	var receiverEmail, emailSubject, emailBody string

	if userType == "staff" {
		scr := model.StaffCreateRequest{}
		if err := c.ShouldBindJSON(&scr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		// random 6 letter password
		password := service.RandomString(8)
		scr.Password = password

		staff, err := createStaffHandler(scr)
		if err != nil {
			response.InternalServerError("Failed to create staff").AbortWithError(c)
			return
		}

		token, err := service.GenerateToken(userType, staff.Account)

		if err != nil {
			response.InternalServerError("Cannot generate token").AbortWithError(c)
			return
		}

		receiverEmail = scr.Email
		emailSubject = "Staff Registration"
		emailBody = "Your account has been created. Your password is " + password +
			". Please validate your account by clicking the link below: " +
			os.Getenv("FRONTEND_URL") +
			os.Getenv("STAFF_VALIDATE_PATH") +
			"?token=" + token
	} else {
		ucr := model.UserCreateRequest{}
		if err := c.ShouldBindJSON(&ucr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		password := service.RandomString(8)
		ucr.Password = password

		user, err := registerUserHandler(ucr)
		if err != nil {
			response.InternalServerError("Failed to create user").AbortWithError(c)
			return
		}

		token, err := service.GenerateToken(userType, user.Account)
		if err != nil {
			response.InternalServerError("Cannot generate token").AbortWithError(c)
			return
		}

		receiverEmail = ucr.Email
		emailSubject = "User Registration"
		emailBody = "Your account has been created. Your password is " + password +
			". Please validate your account by clicking the link below: " +
			os.Getenv("FRONTEND_URL") +
			os.Getenv("STAFF_VALIDATE_PATH") +
			"?token=" + token
	}

	err := service.SendMail(receiverEmail, emailSubject, emailBody)
	if err != nil {
		response.InternalServerError("Failed to send email").AbortWithError(c)
		return
	}

	c.JSON(201, response.CommonResponse{
		Success: true,
	})
}

func ChangePasswordHandler(c *gin.Context) {
	info, exists := c.Get("user")
	if !exists {
		response.InternalServerError("Cannot get account info").AbortWithError(c)
		return
	}

	accountInfo := info.(jwt.MapClaims)
	email := accountInfo["email"].(string)

	cpr := model.ChangePasswordRequest{}
	if err := c.ShouldBindJSON(&cpr); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}

	account, err := service.ValidateCredential(model.LoginRequest{
		Email:    email,
		Password: cpr.OldPassword,
	})
	if err != nil {
		response.BadRequest("Invalid old password").AbortWithError(c)
		return
	}

	encryptedPassword, err := service.EncryptPassword(cpr.NewPassword)
	if err != nil {
		response.InternalServerError("Failed to encrypt password").AbortWithError(c)
		return
	}

	account.Password = encryptedPassword
	if err := account.Update(); err != nil {
		response.InternalServerError("Failed to change password").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}
