package routes

import (
	testController "gin-auth-mongo/controllers/test"

	"github.com/gin-gonic/gin"
)

func TestRoutes(r *gin.RouterGroup) {
	test := r.Group("/test")
	{
		test.POST("/insert", testController.InsertTestData)
		test.POST("/register", testController.Register)
		test.POST("/register2", testController.Register2)
		test.GET("/get", testController.GetTestData)
		test.GET("/get2", testController.GetTestData2)
		test.GET("/get3", testController.GetTestData3)
		test.GET("/get4", testController.GetTestData4)
		test.GET("/get5", testController.GetTestData5)
		test.GET("/get6", testController.GetTestData6)
		test.GET("/get7", testController.GetTestData7)
		test.GET("/get8", testController.GetTestData8)
		test.GET("/check-indexes", testController.CheckIndexes)
	}
}
