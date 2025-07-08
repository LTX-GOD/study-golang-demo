package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=6,max=30"`
	Password      *string            `json:"password" validate:"required,min=6,max=30"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token" `
	Refresh_Token *string            `json:"refresh_token"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updated_at"`
	User_ID       string             `json:"user_id"`
	// 切片本身已经是一个引用类型，能够提供对底层数据的引用，因此不加*号
	UserCart        []ProductUser `json:"usercart" bson:"usercart"`
	Address_Details []Address     `json:"address" bson:"address"`
	Order_Status    []Order       `json:"order" bson:"order"`
}

type Product struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *string            `json:"price"`
	Rating       *string            `json:"rating"`
	Image        *string            `json:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *string            `json:"price"`
	Rating       *string            `json:"rating"`
	Image        *string            `json:"image"`
}

type Address struct {
	Address_id primitive.ObjectID `bson:"_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	PostalCode *string            `json:"postalcode" bson:"postalcode"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list" bson:"order_list"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int                `json:"price" bson:"price"`
	Discount       *int               `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
