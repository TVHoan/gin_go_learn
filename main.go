package main

// only need mysql OR sqlite
// both are included here for reference
import (
	"fmt"
	"github.com/gin-gonic/gin"
	//_ "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	//_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

func database() *gorm.DB {
	//db *gorm.DB
	//err error
	dsn := "host=dpg-cif4npp5rnujc4p57seg-a.singapore-postgres.render.com user=root password=nCWzlSEHnHroRmNBcN3l4q5TgEqVdFxh dbname=gorm_ny4t port=5432 TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	return db
}
func main() {
	// NOTE: See we're using = to assign the global var
	//         	instead of := which would assign it only in this function
	// db, err = gorm.Open("sqlite3", "./gorm.db")
	// db, _ = gorm.Open("mysql", "root:1@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local")
	db := database()
	db.AutoMigrate(&Person{})

	r := gin.Default()
	r.GET("/people/", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)

	r.Run(":8005")
}

func DeletePerson(c *gin.Context) {
	db := database()
	id := c.Params.ByName("id")
	var person Person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdatePerson(c *gin.Context) {
	db := database()
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
	db := database()
	var person Person
	c.BindJSON(&person)

	if result := db.Create(&person).Error; result != nil {
		c.JSON(400, result.Error)
	}
	c.JSON(200, person)
}

func GetPerson(c *gin.Context) {
	db := database()
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
	db := database()
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}

}
