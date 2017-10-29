package main

import (
	"fmt"
	"github.com/gin-gonic/gin"  // Wen framework
	"github.com/jinzhu/gorm"  // ORM
	// (_) se utiliza para evadir la restricción de Go
	// con las variables o paquetes sin utilizar
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Person struct {
	// Se utilizan masyúsculas para setear como campos públicos.
	ID uint `json:"id"`
	FirstName string `json:"firstname"`
	LastName string `json:lastname`
	City string `json:city`
}

func main() {
	/* NOTA: Se utiliza = para asignar las variables globalmente
	// en lugar de utilizar := lo cual lo asigna solo al scope de la función
	*/
	db, err = gorm.Open("sqlite3", "./gorm.db")
	// Para utilizar MySQL
	// db, err = gorm.Open("mysql",
	// "user:pass@tcp(127.0.0.1:3306)/samples?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Permite crear o modificar registros sin afectar los ya existentes
	db.AutoMigrate(&Person{})

	// Creamos la instancia de Gin framework
	r := gin.Default()
	// Declaramos las rutas de nuestra API
	r.GET("/people/", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people/", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)

	// Iniciamos el servidor de nuestra API
	r.Run(":8080")
}

func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}

func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)

	db.Create(&person)
	c.JSON(200, person)
}

func UpdatePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.BindJSON(&person)
		db.Save(&person)
		c.JSON(200, person)
	}
}

func DeletePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")

	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
