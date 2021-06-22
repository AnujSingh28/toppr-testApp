package main

import (
	"context"
	"log"
	"net/http"
	"test-app/cache"
	"test-app/controller"
	"test-app/model"
	"test-app/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	testModel      model.TestModel
	testService    service.TestService
	testController controller.TestController
	queNum         int
	testCache      cache.TestCache
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*.html")

	//server.POST("/signup", Signup)
	server.GET("/start", Init)
	server.GET("/classes", Classes)
	server.POST("/subject", Subject)
	server.POST("/chapters", Chapter)
	server.POST("/score", Score)
	server.POST("/question", Question)
	server.POST("/nextQuestion", NextQuestion)

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}

func Init(ctx *gin.Context) {
	testModel = model.NewTestModel()
	testService = service.New(testModel)
	testCache = cache.NewRedisCache("localhost:6379", 0, 2100, context.Background())
	testController = controller.NewTestController(testService, testCache)
	queNum = 0
	data := gin.H{
		"title": "Instruction Page",
	}
	ctx.HTML(200, "start.html", data)
}

func Classes(ctx *gin.Context) {
	classes := testController.AllClasses()
	data := gin.H{
		"title":   "Select your class",
		"classes": classes,
	}
	ctx.HTML(200, "classes.html", data)
}

func Subject(ctx *gin.Context) {
	subjects := testController.AllSubject(ctx)
	data := gin.H{
		"title":    "Subject Page",
		"subjects": subjects,
	}
	ctx.HTML(200, "subjects.html", data)
}

func Chapter(ctx *gin.Context) {
	chapters := testController.FetchChapter(ctx)
	data := gin.H{
		"title":    "Chapter Page",
		"chapters": chapters,
	}
	ctx.HTML(200, "chapters.html", data)
}

func Question(ctx *gin.Context) {
	testController.FetchQuestions(ctx)
	queNum++
	question := testController.NextQuestion(queNum)
	data := gin.H{
		"title":    "Question Page",
		"question": question,
		"queNum":   queNum,
	}
	ctx.HTML(200, "questions.html", data)
}

func NextQuestion(ctx *gin.Context) {
	if queNum == 10 {
		ctx.Redirect(http.StatusTemporaryRedirect, "/score")
		return
	}
	queNum++
	testController.UpdateScore(ctx, queNum-1)
	question := testController.NextQuestion(queNum)
	data := gin.H{
		"title":    "Question Page",
		"question": question,
		"queNum":   queNum,
	}
	ctx.HTML(200, "questions.html", data)
}

func Score(ctx *gin.Context) {
	queNum := 10
	testController.UpdateScore(ctx, queNum)
	score := testController.UpdateTest()
	data := gin.H{
		"title": "Result Page",
		"score": score,
	}
	ctx.HTML(200, "score.html", data)
}
