package middleware

import (
	"fmt"
	"gin/database"
	"gin/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	//_ "gorm.io/driver/mysql"
)

func Authorization(c *gin.Context) {

	tokenString, errcookie := c.Cookie("Authorization")
	if errcookie != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
				var user model.User
				db := database.Database()
				db.First(&user, claims["sub"])
				if user.ID == 0 {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
				c.Set("user", user)
				c.Next()
			} else {
				fmt.Println(err)
			}
		}
	}
}
