package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zsm/ecommerce-sys/database"
	"github.com/zsm/ecommerce-sys/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckEmptyParam(c *gin.Context, paramValue, paramName string) bool {
	if paramValue == "" {
		log.Printf("%s is empty, please provide a valid %s\n", paramName, paramName)
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s is empty", paramName))
		return true
	}
	return false
}

//业务逻辑层骨架

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	if prodCollection == nil || userCollection == nil {
		log.Fatal("prodCollection or userCollection is nil")
	}
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddtoCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		if CheckEmptyParam(c, productQueryID, "product id") {
			return
		}
		if CheckEmptyParam(c, userQueryID, "user id") {
			return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Invalid product ID format:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//调用数据库添加商品
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil {
			log.Println("Error adding product to cart:", err)
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully added product to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		if CheckEmptyParam(c, productQueryID, "product id") {
			return
		}
		if CheckEmptyParam(c, userQueryID, "user id") {
			return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Invalid product ID format:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//调用数据库删除商品
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil {
			log.Println("Error removing cart item:", err)
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// 成功移除商品，返回 200 状态和成功消息
		c.IndentedJSON(http.StatusOK, "Successfully removed the item")

	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if CheckEmptyParam(c, user_id, "user id") {
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User

		//查找用户
		var UserCollection = database.UserData(database.Client, "Users")
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		//查询
		//先找用户id
		filter_match := bson.D{{Key: "$match", Value: bson.D{
			primitive.E{
				Key:   "_id",
				Value: usert_id,
			},
		}}}

		//将用户购物车数组拆分多个
		unwind := bson.D{{Key: "$unwind", Value: bson.D{
			primitive.E{
				Key:   "path",
				Value: "$usercart",
			},
		}}}

		//按照用户ID分组，并且求和
		grouping := bson.D{{Key: "$group", Value: bson.D{
			primitive.E{
				Key:   "_id",
				Value: "$_id",
			},
			{
				Key: "total",
				Value: bson.D{
					primitive.E{
						Key:   "$sum",
						Value: "$usercart.price",
					},
				},
			},
		}}}

		//聚合查询
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{
			filter_match, unwind, grouping,
		})
		if err != nil {
			log.Println(err)
		}

		//保存查询结果
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		//遍历结果发给客户端
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, json["UserCart"])
		}

		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if CheckEmptyParam(c, userQueryID, "user id") {
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.IndentedJSON(200, "成功下单")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		if CheckEmptyParam(c, productQueryID, "product id") {
			return
		}
		if CheckEmptyParam(c, userQueryID, "user id") {
			return
		}

		_, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Error processing instant buy:", err)
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, "Successfully placed the order")
	}
}
