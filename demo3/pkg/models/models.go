package models

import (
	"zsm/pkg/config"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	// 使用 GORM 的 AutoMigrate 方法自动迁移 Book 结构体。
	// 自动迁移会创建或更新数据库表，使其与 Book 结构体匹配。如果表不存在，则创建表；如果表已存在，则更新表结构以匹配 Book 结构体的定义。
	db.AutoMigrate(&Book{})
}

// CreateBook 方法用于创建 Book 结构体的实例并插入到数据库中。
func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

// GetAllBooks 方法用于从数据库中获取所有 Book 结构体的实例。
func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

// GetBookById 方法用于从数据库中获取指定 ID 的 Book 结构体的实例。
func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", Id).Find(&getBook)
	return &getBook, db
}

// DeleteBook 方法用于从数据库中删除指定 ID 的 Book 结构体的实例。
func DeleteBook(ID int64) Book {
	var book Book
	db.Where("ID=?", ID).Delete(book)
	return book
}
