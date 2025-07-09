package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zsm/CurrencyExchangeApp/global"
	"github.com/zsm/CurrencyExchangeApp/models"
	"github.com/zsm/CurrencyExchangeApp/utils"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := global.Db.Where("username=?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd, err := utils.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashedPwd

	token, err := utils.GenerateJWT(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
