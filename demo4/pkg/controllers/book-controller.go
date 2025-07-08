package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zsm/go-bookstore/pkg/models"
)

// NewBook 是一个用来创建新书的结构体
var NewBook models.Book

// GetBookTest 函数返回所有的书籍
func GetBookTest(c *gin.Context) {
	newBooks := models.GetAllBooks()
	c.JSON(http.StatusOK, newBooks)
}

// GetBookByIdTest 函数返回指定ID的书籍
func GetBookByIdTest(c *gin.Context) {
	bookId := c.Param("bookId")
	ID, _ := strconv.ParseInt(bookId, 0, 0)
	bookDetails, _ := models.GetBookById(ID)
	// 以JSON格式返回书籍详情
	c.JSON(http.StatusOK, bookDetails)
}

// CreateBookTest 函数创建一个新的书籍并返回详细信息
func CreateBookTest(c *gin.Context) {
	// 初始化一个新书籍结构体
	var CreateBook = &models.Book{}
	// 从请求体中解析书籍信息
	c.ShouldBindJSON(CreateBook)
	// 创建书籍并保存到数据库
	b := CreateBook.CreateBook()
	// 以JSON格式返回书籍详情
	c.JSON(http.StatusOK, b)
}

func DeleteBookTest(c *gin.Context) {
	bookId := c.Param("bookId")
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Printf("解析错误！")
	}
	// 删除之前获取书籍详情, 并删除书籍
	book, _ := models.GetBookById(ID)
	models.DeleteBook(ID)
	c.JSON(http.StatusOK, book)
}

func UpdateBookTest(c *gin.Context) {
	var updateBook = &models.Book{}
	c.ShouldBindJSON(updateBook)
	bookId := c.Param("bookId")
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("解析错误！")
	}
	bookDetails, db := models.GetBookById(ID)
	// 如果Name、Author、Publication字段有更新，则更新数据库对应数据
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	// 保存更新后的书籍信息到数据库
	db.Save(&bookDetails)
	c.JSON(http.StatusOK, bookDetails)
}
