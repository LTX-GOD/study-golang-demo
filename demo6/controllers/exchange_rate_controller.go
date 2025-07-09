package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zsm/CurrencyExchangeApp/global"
	"github.com/zsm/CurrencyExchangeApp/models"
	"gorm.io/gorm"
)

func CreateExchangeRate(c *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := c.ShouldBindJSON(&exchangeRate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exchangeRate)
}

func GetExchangeRates(c *gin.Context) {
	var exchangeRates []models.ExchangeRate

	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, exchangeRates)
}
