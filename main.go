package main

// only need mysql OR sqlite
// both are included here for reference
import (
	"fmt"
	"gin/controller/auth"
	"gin/model"

	"github.com/gin-gonic/gin"

	"gin/database"
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
	db := database.Database()
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&model.User{})

	r := gin.New()
	r.Use(middleware.Authorization())
	r.POST("/login/", auth.Login)
	r.GET("/people/", GetPeople)
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
