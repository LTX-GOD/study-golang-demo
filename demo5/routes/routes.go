package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zsm/ecommerce-sys/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controllers.SignUp())//注册
	incomingRoutes.POST("/user/login", controllers.Login())//登陆
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())// 管理员浏览商品
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())   // 查询所有商品
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery()) // 通过 ID 查询商品
}
