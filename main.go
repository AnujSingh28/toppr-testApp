package main

import (
	"test-app/Config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

// Chapter model
type Chapters struct {
	ID      uint   `json:"chapter_id"`
	Class   uint   `json:"class"`
	Subject uint   `json:"subject"` // phy:1, chem:2, math:3, bio:4
	Chapter string `json:"chapter"`
}

// Question model
type Questions struct {
	ID       uint   `json:"id"`
	Chapter  string `json:"chapter_id"`
	Question string `json:"question"`
	Op1      string `json:"op1"`
	Op2      string `json:"op2"`
	Op3      string `json:"op3"`
	Op4      string `json:"op4"`
	Ans      uint   `json:"ans"`
}

func main() {
	db, err = gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Chapters{}, &Questions{})

	c1 := Chapters{Class: 9, Subject: 1, Chapter: "Laws of Motion"}
	c2 := Chapters{Class: 10, Subject: 1, Chapter: "Optics"}
	db.Create(&c1)
	db.Create(&c2)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/classes", AllClasses)
	r.GET("/classes/:class_id", GetSubjects)
	r.GET("/classes/:class_id/:subject_id", GetChapters)
	r.GET("/test/:chapter_id", GetTest)

	r.Use((cors.Default()))
	r.Run(":8080")

}

func AllClasses(c *gin.Context) {
	c.HTML(200, "classes.html", gin.H{"message": "rendering html"})
	c.JSON(200, gin.H{"message": "all classes"})
}

func GetChapters(c *gin.Context) {
	c.JSON(200, gin.H{"message": "all chapters"})
}

func GetSubjects(c *gin.Context) {
	c.HTML(200, "subjects.html", gin.H{"message": "rendering html"})
	c.JSON(200, gin.H{"message": "all subjects"})
}

func GetTest(c *gin.Context) {
	c.JSON(200, gin.H{"message": "test for the choosen"})
}
