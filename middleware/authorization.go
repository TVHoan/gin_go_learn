package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	//_ "gorm.io/driver/mysql"
)

func Authorization() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println(c)

		c.Next()

	}
}
