package controller

import (
	"fmt"
	"os"

	"github.com/FakJeongTeeNhoi/user-management/middleware"
	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *gin.Context) {
	userType := c.Query("type")
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

	userInfo, _ := service.GetInfoFromToken(token)

	c.SetCookie("token", token, 3600, "/", "", false, false)
	c.JSON(201, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"token":    token,
		"userInfo": userInfo,
	}))
}

func LogoutHandler(c *gin.Context) {
	token := middleware.ExtractToken(c.GetHeader("authorization"))
	middleware.BlackListToken = append(middleware.BlackListToken, token)
	c.SetCookie("token", "", -1, "/", "", false, false)
	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}

func VerifyHandler(c *gin.Context) {
	info, exists := c.Get("user")
	if !exists {
		response.InternalServerError("Cannot get account info").AbortWithError(c)
		return
	}

	accountInfo := info.(jwt.MapClaims)
	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(accountInfo))
}

func RegisterHandler(c *gin.Context) {
	userType := c.Query("type")
	if userType != "staff" && userType != "user" {
		response.BadRequest("Invalid user type").AbortWithError(c)
		return
	}

	var receiverEmail, emailSubject, emailBody string
	var user model.User
	var staff model.Staff

	password := service.RandomString(8)
	fmt.Println(password)

	if userType == "staff" {
		scr := model.StaffCreateRequest{}
		if err := c.ShouldBindJSON(&scr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		scr.Password = password

		var err error
		staff, err = createStaffHandler(scr)
		if err != nil {
			response.InternalServerError("Failed to create staff").AbortWithError(c)
			return
		}

		token, err := service.GenerateToken(userType, staff.Account)

		if err != nil {
			_ = staff.Delete()

			response.InternalServerError("Cannot generate token").AbortWithError(c)
			return
		}

		receiverEmail = scr.Email
		emailSubject = "Staff Registration"
		emailBody = "Your account has been created. Your password is <b>" + password +
			"</b>. <br> Please validate your account by clicking the link below: <a href='" +
			os.Getenv("FRONTEND_URL") +
			os.Getenv("STAFF_VERIFY_PATH") +
			"?name=" + staff.Account.Name + "&token=" + token + "'>Validate</a>"
	} else {
		ucr := model.UserCreateRequest{}
		if err := c.ShouldBindJSON(&ucr); err != nil {
			response.BadRequest("Invalid request").AbortWithError(c)
			return
		}

		ucr.Password = password

		var err error
		user, err = registerUserHandler(ucr)
		if err != nil {
			response.InternalServerError("Failed to create user").AbortWithError(c)
			return
		}

		token, err := service.GenerateToken(userType, user.Account)
		if err != nil {
			_ = user.Delete()

			response.InternalServerError("Cannot generate token").AbortWithError(c)
			return
		}

		receiverEmail = ucr.Email
		emailSubject = "User Registration"
		emailBody = "Your account has been created. Your password is <b>" + password +
			"</b>. <br> Please validate your account by clicking the link below: <a href='" +
			os.Getenv("FRONTEND_URL") +
			os.Getenv("USER_VERIFY_PATH") +
			"?name=" + user.Account.Name + "&token=" + token + "'>Validate</a>"
	}

	err := service.SendMail(receiverEmail, emailSubject, emailBody)
	if err != nil {

		if userType == "staff" {
			_ = staff.Delete()
		} else {
			_ = user.Delete()
		}

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

func GetAccountInfoHandler(c *gin.Context) {
	info, exists := c.Get("user")
	if !exists {
		response.InternalServerError("Cannot get account info").AbortWithError(c)
		return
	}

	accountInfo := info.(jwt.MapClaims)
	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(accountInfo))
}
