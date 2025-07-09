package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/zsm/CurrencyExchangeApp/global"
)

func LikeArticle(c *gin.Context) {
	articleID := c.Param("id")

	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

func GetArticleLikes(c *gin.Context) {
	articleID := c.Param("id")

	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}
