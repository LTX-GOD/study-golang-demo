package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zsm/ecommerce-sys/database"
	"github.com/zsm/ecommerce-sys/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//业务逻辑层骨架

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "用户ID无效"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.IndentedJSON(500, "用户ID格式错误")
			return
		}

		var addresses models.Address
		addresses.Address_id = primitive.NewObjectID()

		// 解析请求体中的地址信息
		if err := c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 更新用户的地址列表
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		userCollection := database.UserData(database.Client, "Users")
		pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}

		// 限制地址数量
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "地址数量已达上限")
			return
		}

		defer pointcursor.Close(ctx)
		c.IndentedJSON(200, "成功添加地址")
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "用户ID无效"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.IndentedJSON(500, "用户ID格式错误")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 更新家庭地址（第一个地址）
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			primitive.E{Key: "address.0.house_name", Value: editaddress.House},
			primitive.E{Key: "address.0.street_name", Value: editaddress.Street},
			primitive.E{Key: "address.0.city_name", Value: editaddress.City},
			primitive.E{Key: "address.0.postalcode", Value: editaddress.PostalCode},
		}}}

		userCollection := database.UserData(database.Client, "Users")
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "无法更新家庭地址")
			return
		}

		c.IndentedJSON(200, "成功更新家庭地址")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "用户ID无效"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.IndentedJSON(500, "用户ID格式错误")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 更新工作地址（第二个地址）
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			primitive.E{Key: "address.1.house_name", Value: editaddress.House},
			primitive.E{Key: "address.1.street_name", Value: editaddress.Street},
			primitive.E{Key: "address.1.city_name", Value: editaddress.City},
			primitive.E{Key: "address.1.postalcode", Value: editaddress.PostalCode},
		}}}

		userCollection := database.UserData(database.Client, "Users")
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "无法更新工作地址")
			return
		}

		c.IndentedJSON(200, "成功更新工作地址")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "用户ID无效"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.IndentedJSON(500, "用户ID格式错误")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 清空用户的所有地址
		emptylist := make([]models.Address, 0)
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: emptylist}}}}

		userCollection := database.UserData(database.Client, "Users")
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "无法删除地址")
			return
		}

		c.IndentedJSON(200, "成功删除所有地址")
	}
}
