package tokens

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/zsm/ecommerce-sys/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string
	Name  string
	Uid   string
	jwt.StandardClaims
}

// UserData 是存储用户数据的 MongoDB 集合引用
var UserData *mongo.Collection = database.UserData(database.Client, "Users")

// 从环境变量中读取JWT的签名和认证
var SECRET_KEY = os.Getenv("SECRET_KEY")

// TokenGenerator 生成一个签名的访问令牌和一个签名的刷新令牌。
func TokenGenerator(email string, name string, uid string) (signedtoken string, signedrefeshtoken string, err error) {
	claims := &SignedDetails{
		Email: email,
		Name:  name,
		Uid:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(), // 令牌有效期为24小时
		},
	}

	//创建一个仅包含过期时间的声明，用来刷新令牌
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(), //刷新有效七天
		},
	}

	//HS256访问令牌
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	//刷新
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshtoken, err
}

// ValidateToken 验证给定的签名令牌是否有效，并返回其声明。
func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	// 解析并验证签名令牌，使用提供的密钥和声明类型
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil // 使用SECRET_KEY作为签名密钥
	})
	if err != nil {
		msg = err.Error() // 如果解析过程中出现错误，设置错误信息并返回
		return
	}

	// 断言token.Claims为*SignedDetails类型，并进行类型检查
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Invalid token" // 如果断言失败，说明令牌无效，设置错误信息并返回
		return
	}

	// 检查令牌的过期时间
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token expired" // 如果令牌已过期，设置错误信息并返回
		return
	}

	// 如果所有检查都通过，返回令牌中的声明和一个空消息
	return claims, ""
}

// UpdateAllTokens 更新用户的访问令牌和刷新令牌，并记录更新时间。
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {

	// 创建一个带有超时的上下文，超时时间为100秒
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel() // 确保函数返回时取消上下文

	var updateobj primitive.D

	// 构建更新对象，包括访问令牌、刷新令牌和更新时间
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) // 格式化当前时间为RFC3339格式

	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: updated_at})

	// 设置Upsert选项，表示如果用户不存在则插入新记录
	upsert := true
	filter := bson.M{"user_id": userid} // 设置过滤条件，匹配指定的用户ID
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	// 执行更新操作，将更新对象应用到符合过滤条件的文档中
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	}, &opt)

	// 处理更新操作中的错误
	if err != nil {
		log.Panic(err) // 记录错误并引发恐慌
		return
	}
}
