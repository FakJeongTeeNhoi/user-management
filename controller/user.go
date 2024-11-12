package controller

import (
	"strconv"

	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"
	"github.com/gin-gonic/gin"
)

func registerUserHandler(ucr model.UserCreateRequest) (model.User, error) {
	encryptedPassword, err := service.EncryptPassword(ucr.Password)
	if err != nil {
		return model.User{}, err
	}
	ucr.Password = encryptedPassword

	user := ucr.ToUser()
	err = user.Create()
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func GetAllUsersHandler(c *gin.Context) {
	users := model.Users{}
	if err := users.GetAll(nil); err != nil {
		response.InternalServerError("Failed to get users").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"count": len(users),
		"users": users,
	}))
}

func GetUsersHandler(c *gin.Context) {
	uif := model.UsersInfoRequest{}
	
	if err := c.ShouldBindJSON(&uif); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}
	users := model.Users{}
	
	for _, participant := range uif.UserList {
		user := model.User{}
		if err := user.GetOne(map[string]interface{}{"user_id": strconv.FormatUint(uint64(participant), 10)}); err != nil {
			response.NotFound("User not found").AbortWithError(c)
			return
		}
		users = append(users, user)
	}
	
	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"count": len(users),
		"users": users,
	}))
}

func GetUserHandler(c *gin.Context) {
	accountId := c.Param("accountId")
	user := model.User{}

	if err := user.GetOne(map[string]interface{}{"account_id": service.ParseToUint(accountId)}); err != nil {
		response.NotFound("User not found").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"user": user,
	}))
}

func UpdateUserHandler(c *gin.Context) {
	uur := model.UserUpdateRequest{}
	if err := c.ShouldBindJSON(&uur); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}

	user := model.User{}
	if err := user.GetOne(map[string]interface{}{"account_id": uur.ID}); err != nil {
		response.NotFound("User not found").AbortWithError(c)
		return
	}

	user = uur.ToUser(user)
	if err := user.Update(); err != nil {
		response.InternalServerError("Failed to update user").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}

func DeleteUserHandler(c *gin.Context) {
	accountId := c.Param("accountId")

	user := model.User{}
	if err := user.GetOne(map[string]interface{}{"account_id": service.ParseToUint(accountId)}); err != nil {
		response.NotFound("User not found").AbortWithError(c)
		return
	}

	if err := user.Delete(); err != nil {
		response.InternalServerError("Failed to delete user").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}
