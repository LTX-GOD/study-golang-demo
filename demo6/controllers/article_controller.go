package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/zsm/CurrencyExchangeApp/global"
	"github.com/zsm/CurrencyExchangeApp/models"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(c *gin.Context) {
	var article models.Article

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func GetArticles(c *gin.Context) {

	cachedData, err := global.RedisDB.Get(cacheKey).Result()

	if err == redis.Nil {
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, articles)

	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		var articles []models.Article

		if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, articles)
	}
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	var article models.Article

	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, article)
}
