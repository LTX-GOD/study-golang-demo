package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zsm/ecommerce-sys/controllers"
	"github.com/zsm/ecommerce-sys/database"
	"github.com/zsm/ecommerce-sys/middleware"
	"github.com/zsm/ecommerce-sys/routes"
)

func main() {
	// 获取环境变量PORT的值, 如果不存在则赋值8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// 创建应用程序实例
	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"),
		database.UserData(database.Client, "Users"),
	)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "token"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// 注册
	routes.UserRoutes(router) // 调用routs包中的UserRoutes函数,注册路由,并命名为router
	router.Use(middleware.Authentication())

	// 定义用户路由之外的路由
	router.GET("/addtocart", app.AddtoCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
