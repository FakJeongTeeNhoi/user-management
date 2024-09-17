package router

import (
	"github.com/FakJeongTeeNhoi/user-management/controller"
	"github.com/gin-gonic/gin"
)

func StaffRouterGroup(server *gin.RouterGroup) {
	staff := server.Group("/staff")
    {
        staff.POST("/create-staff", controller.CreateStaffHandler)
		staff.GET("/", controller.GetAllStaffHandler)
		staff.GET("/:accountId", controller.GetStaffHandler)
		staff.PUT("/", controller.UpdateStaffHandler)
		staff.DELETE("/:accountId", controller.DeleteStaffHandler)
    }
}