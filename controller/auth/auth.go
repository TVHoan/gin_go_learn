package auth

import (
	_ "fmt"
	"gin/database"
	"gin/model"
	"github.com/gin-gonic/gin"
	//_ "gorm.io/driver/mysql"
	_ "gin/model"
	_ "gorm.io/driver/postgres"
	//_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

func Login(c *gin.Context) {
	db := database.Database()
	username := c.Param("username")
	password := c.Param("password")
	var user model.User
	if err := db.Where("Username = ? and Password = ? ", username, password).First(&user); err != nil {
		c.AbortWithStatus(500)
	} else {
		c.JSON(200, user)
	}
}
