package main

// only need mysql OR sqlite
// both are included here for reference
import (
	"fmt"
	"gin/controller/auth"
	"gin/database"
	"gin/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"

	//_ "gorm.io/driver/sqlite"
	"gin/middleware"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.Database()
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&model.User{})
	secretKey := os.Getenv("SECRET")
	println(secretKey)
	r := gin.New()
	r.POST("/login/", auth.Login)
	r.POST("/register/", auth.Register)
	r.GET("/people/", middleware.Authorization, GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)

	r.Run(":8005")
}

func DeletePerson(c *gin.Context) {
	db := database.Database()
	id := c.Params.ByName("id")
	var person Person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdatePerson(c *gin.Context) {
	db := database.Database()
	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)

	db.Save(&person)
	c.JSON(200, person)

}

func CreatePerson(c *gin.Context) {
	db := database.Database()
	var person Person
	c.BindJSON(&person)

	if result := db.Create(&person).Error; result != nil {
		c.JSON(400, result.Error)
	}
	c.JSON(200, person)
}

func GetPerson(c *gin.Context) {
	db := database.Database()
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}
func GetPeople(c *gin.Context) {
	db := database.Database()
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}

}
