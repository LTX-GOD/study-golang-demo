package controllers

//业务逻辑层骨架
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zsm/ecommerce-sys/database"
	"github.com/zsm/ecommerce-sys/models"
	generate "github.com/zsm/ecommerce-sys/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VertifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "用户名或密码错误"
		valid = false
	}

	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// 解析请求的 JSON 数据到 user 结构体
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 验证用户数据
		Validate := validator.New()
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr,
			})
			return
		}

		// 检查邮箱是否已被注册
		UserCollection := database.UserData(database.Client, "Users")
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Println("Email check failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "邮箱检查失败",
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "用户已存在！",
			})
			return
		}

		//检查手机号
		password := HashPassword(*user.Password)
		user.Password = &password

		//设置创建时间和更新时间
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		//生成用户ID和令牌
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken

		//初始化用户购物车、地址还有订单状态
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		//将用户插入数据库
		_, errorss := UserCollection.InsertOne(ctx, user)
		if errorss != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "用户创建失败",
			})
			return
		}

		c.JSON(http.StatusOK, "成功注册")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var founduser models.User
		var user models.User

		//解析请求的json数据岛user结构体
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//根据邮箱确定用户
		UserCollection := database.UserData(database.Client, "Users")
		err := UserCollection.FindOne(ctx, bson.M{
			"email": user.Email,
		}).Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "账号或密码错误",
			})
			return
		}

		//验证密码
		PasswordIsValid, msg := VertifyPassword(*founduser.Password, *user.Password)
		if !PasswordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": msg,
			})
			fmt.Println(msg)
			return
		}

		//生成新的token
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.Name, founduser.User_ID)

		//更新token
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		//返回信息和token
		c.JSON(http.StatusOK, gin.H{
			"user":         founduser,
			"token":        token,
			"refreshToken": refreshToken,
		})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//分配新的ID
		products.Product_ID = primitive.NewObjectID()

		//插入数据库
		productCollection := database.ProductData(database.Client, "Products")
		_, anyerr := productCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "未能插入产品",
			})
			return
		}

		c.JSON(http.StatusOK,"成功添加商品")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productlist []models.Product
		productCollection := database.ProductData(database.Client, "Products")

		// 查询所有商品
		cursor, err := productCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "查询商品时出错")
			return
		}

		// 解码查询结果
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if len(productlist) == 0 {
			log.Println("没有找到商品")
			c.JSON(http.StatusNotFound, gin.H{"error": "没有找到商品"})
			return
		}

		c.IndentedJSON(http.StatusOK, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productlist []models.Product
		productCollection := database.ProductData(database.Client, "Products")

		// 获取查询参数
		queryParam := c.Query("name")
		if queryParam == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "查询参数无效"})
			c.Abort()
			return
		}

		// 使用正则表达式进行模糊搜索
		filter := bson.M{
			"product_name": bson.M{
				"$regex":   queryParam,
				"$options": "i", // 不区分大小写
			},
		}

		// 执行查询
		cursor, err := productCollection.Find(ctx, filter)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "查询商品时出错")
			return
		}

		// 解码查询结果
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if len(productlist) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "没有找到匹配的商品"})
			return
		}

		c.IndentedJSON(http.StatusOK, productlist)
	}
}
