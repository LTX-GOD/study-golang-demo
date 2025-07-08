package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zsm/go-bookstore/pkg/controllers"
)

func Router() *gin.Engine {
	r := gin.Default()
	book := r.Group("/book")
	{
		book.GET("/", controllers.GetBookTest)
		book.GET("/:bookId", controllers.GetBookByIdTest)
		book.POST("/", controllers.CreateBookTest)
		book.PUT("/:bookId", controllers.UpdateBookTest)
		book.DELETE("/:bookId", controllers.DeleteBookTest)
	}
	return r
}
