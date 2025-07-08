package database

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/zsm/ecommerce-sys/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func StringToInt(s *string) int {
	if s == nil {
		return 0
	}
	i, err := strconv.Atoi(*s)
	if err != nil {
		return 0 // 或者你也可以 panic(err) 根据你项目的需要喵~
	}
	return i
}

var (
	ErrCantFindProduct    = errors.New("can't find the product")                    // 表示找不到产品的错误。
	ErrCantDecodeProducts = errors.New("can't find the product")                    // 表示解码产品失败的错误
	ErrUserIdIsNotValid   = errors.New("this user is not valid")                    // 表示用户 ID 无效的错误。
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")       // 表示无法更新用户的错误。
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")     // 表示无法从购物车中移除项的错误。
	ErrCantGetItem        = errors.New("was unnable to get the item from the cart") //表示无法从购物车中获取项的错误。
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")                // 表示无法更新购买的错误。
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	// 查找单个商品
	var product models.Product
	err := prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	//用户ID转换类型
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	//更新用户购物车，添加单个产品
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			primitive.E{Key: "usercart", Value: product},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	// 定义过滤条件和更新操作
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	//更新并且移除
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	//初始化订单状态
	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	// 聚合操作：计算购物车中所有商品的总金额
	unwind := bson.D{{Key: "$unwind", Value: bson.D{
		primitive.E{Key: "path", Value: "$usercart"},
	}}}

	grouping := bson.D{{Key: "$group", Value: bson.D{
		primitive.E{Key: "_id", Value: "$_id"},
		primitive.E{Key: "total", Value: bson.D{
			primitive.E{Key: "$sum", Value: "$usercart.price"},
		}},
	}}}

	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if errAll := currentresults.All(ctx, &getusercart); errAll != nil {
		panic(errAll)
	}

	// 计算总价格
	var total_price int32
	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	//将订单信息加到用户的订单列表里
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$push", Value: bson.D{
			primitive.E{Key: "order", Value: ordercart},
		}},
	}

	_, errr := userCollection.UpdateOne(ctx, filter, update)
	if errr != nil {
		log.Println(err)
	}

	//从用户文档中读取购物车内容
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id"}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}

	// 将购物车中的所有商品添加到订单列表中
	filterNew := bson.D{primitive.E{Key: "_id", Value: id}}
	updateNew := bson.M{"$push": bson.M{"order.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filterNew, updateNew)
	if err != nil {
		log.Println(err)
	}

	//清空用户的购物车
	usercart_empty := make([]models.ProductUser, 0)
	filterNewCart := bson.D{primitive.E{Key: "_id", Value: id}}
	updateNewCart := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "usercart", Value: usercart_empty},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filterNewCart, updateNewCart)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_detail models.Order

	//创建一个新的订单
	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true

	//拿详细信息
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	orders_detail.Price = StringToInt(product_details.Price)

	//更新用户集合，插入新订单
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	userCollection.UpdateOne(ctx, filter, update)

	//插入详细信息
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	return nil
}
