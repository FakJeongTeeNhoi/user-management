package controller

import (
	"github.com/FakJeongTeeNhoi/user-management/model"
	"github.com/FakJeongTeeNhoi/user-management/model/response"
	"github.com/FakJeongTeeNhoi/user-management/service"
	"github.com/gin-gonic/gin"
)

func createStaffHandler(stf model.StaffCreateRequest) (model.Staff, error) {
	encryptedPassword, err := service.EncryptPassword(stf.Password)
	if err != nil {
		return model.Staff{}, err
	}
	stf.Password = encryptedPassword

	staff := stf.ToStaff()
	err = staff.Create()
	if err != nil {
		return model.Staff{}, err
	}

	return staff, nil
}

func GetAllStaffHandler(c *gin.Context) {
	staffs := model.Staffs{}
	if err := staffs.GetAll(nil); err != nil {
		response.InternalServerError("Failed to get all staffs").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"count":  len(staffs),
		"staffs": staffs,
	}))
}

func GetStaffHandler(c *gin.Context) {
	accountId := c.Param("accountId")
	staff := model.Staff{}

	if err := staff.GetOne(map[string]interface{}{"account_id": service.ParseToUint(accountId)}); err != nil {
		response.NotFound("Staff not found").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"staff": staff,
	}))
}

func UpdateStaffHandler(c *gin.Context) {
	sur := model.StaffUpdateRequest{}
	if err := c.ShouldBindJSON(&sur); err != nil {
		response.BadRequest("Invalid request").AbortWithError(c)
		return
	}

	staff := model.Staff{}
	if err := staff.GetOne(map[string]interface{}{"account_id": sur.ID}); err != nil {
		response.NotFound("Staff not found").AbortWithError(c)
		return
	}

	staff = sur.ToStaff(staff)
	if err := staff.Update(); err != nil {
		response.InternalServerError("Failed to update staff").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}

func DeleteStaffHandler(c *gin.Context) {
	accountId := c.Param("accountId")

	staff := model.Staff{}
	if err := staff.GetOne(map[string]interface{}{"account_id": service.ParseToUint(accountId)}); err != nil {
		response.NotFound("Staff not found").AbortWithError(c)
		return
	}

	if err := staff.Delete(); err != nil {
		response.InternalServerError("Failed to delete staff").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	})
}
