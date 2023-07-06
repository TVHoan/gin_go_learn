package auth

import (
	_ "fmt"
	"gin/database"
	"gin/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"

	//_ "gorm.io/driver/mysql"
	_ "gin/model"
	_ "github.com/dgrijalva/jwt-go"
	_ "gorm.io/driver/postgres"
	//_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

type LoginDto struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type RegisterDto struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func Login(c *gin.Context) {
	secret := os.Getenv("SECRET")
	db := database.Database()
	var user model.User
	var input LoginDto
	if error := c.BindJSON(&input); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Username and Password is required "})
	}
	userrecord := db.Select("id", "password_hash").Where("user_name = ? ", input.Username).First(&user)
	if userrecord.Error != nil {
		c.JSON(http.StatusFound, gin.H{"Error": " Not Found UserName"})
	}
	if error := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Username or Password",
		})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 24*60, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{"token": tokenString,
			"expired": time.Now().Add(time.Minute * 24 * 60)})

	}

}
func Register(c *gin.Context) {
	db := database.Database()
	var body RegisterDto
	c.BindJSON(&body)
	hash, errhash := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if errhash != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hash Password faile",
		})
	}
	user := model.User{UserName: body.Username, PasswordHash: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, user)
	}

}
func ValidateToken(c *gin.Context) {

}
