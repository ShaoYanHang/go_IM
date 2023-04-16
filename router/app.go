package router

import (
	"IM/docs"
	"IM/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)


func Router() *gin.Engine{
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/index", service.GetIndex)
	r.GET("/user/GetUserList", service.GetUserList)
	r.POST("/user/CreateUser", service.CreateUser)
	r.DELETE("/user/DeleteUser", service.DeleteUser)
	r.POST("/user/UpdateUser", service.UpdateUser)
	r.POST("/user/FindUserByNameAndPwd", service.FindUserByNameAndPwd)

	r.GET("/user/SendMsg", service.SendMsg)
	return r
}
