package controller

import (
	"strconv"
	"test-app/cache"
	"test-app/entity"
	"test-app/service"

	"github.com/gin-gonic/gin"
)

type TestController interface {
	AllClasses() []entity.Classes
	AllSubject(ctx *gin.Context) []entity.Subjects
	FetchChapter(ctx *gin.Context) []entity.Chapters
	FetchQuestions(ctx *gin.Context) ([]entity.Questions, uint64)
	NextQuestion(queNum int) entity.Questions
	UpdateScore(ctx *gin.Context, queNum int)
	UpdateTest() int
	Signup(ctx *gin.Context) entity.User
	Login(ctx *gin.Context) entity.User
}

type testController struct {
	service service.TestService
	cache   cache.TestCache
}

func NewTestController(service service.TestService, cache cache.TestCache) TestController {
	return &testController{
		service: service,
		cache:   cache,
	}
}

func (con *testController) Signup(ctx *gin.Context) entity.User {
	var user entity.User
	ctx.BindJSON(&user)
	return con.service.Signup(ctx, user)
}

func (con *testController) Login(ctx *gin.Context) entity.User {
	var user entity.User
	ctx.BindJSON(&user)
	return con.service.Login(ctx, user)
}

func (con *testController) AllClasses() []entity.Classes {
	return con.service.AllClasses()
}

func (con *testController) AllSubject(ctx *gin.Context) []entity.Subjects {
	var class entity.Classes
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		panic(err)
	}
	class.ID = uint64(id)
	class.Name = ctx.PostForm("name")
	return con.service.AllSubject(class)
}

func (con *testController) FetchChapter(ctx *gin.Context) []entity.Chapters {
	var subject entity.Subjects
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		panic(err)
	}
	subject.ID = uint64(id)
	subject.Name = ctx.PostForm("name")
	chapter := con.service.FetchChapter(subject)
	return chapter
}

func (con *testController) FetchQuestions(ctx *gin.Context) ([]entity.Questions, uint64) {
	var chapter entity.Chapters
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		panic(err)
	}
	chapter.ID = uint64(id)
	chapter.Name = ctx.PostForm("name")
	questions, testID := con.service.FetchQuestions(chapter)
	con.cache.SetTestID("TestID", testID) // Cache TestID
	testIDStr := strconv.Itoa(int(testID))
	scoreKey := "Score_" + testIDStr
	con.cache.SetScore(scoreKey, 0) // Cache Score = 0
	var questionKey string
	for i, val := range questions {
		questionKey = strconv.Itoa(i+1) + "_" + testIDStr
		con.cache.SetQuestion(questionKey, val) // Cache all the questions  (key, val) : (i_TestID, question[i])
	}
	return questions, testID
}

func (con *testController) NextQuestion(queNum int) entity.Questions {
	testID := con.cache.GetTestID("TestID")
	testIDStr := strconv.Itoa(int(testID))
	questionKey := strconv.Itoa(queNum) + "_" + testIDStr
	question := con.cache.GetQuestion(questionKey)
	return question
}

func (con *testController) UpdateScore(ctx *gin.Context, queNum int) {
	answerSubmitted, err := strconv.Atoi(ctx.PostForm("answerSubmitted"))
	if err != nil {
		panic("Failed to convert string to int")
	}
	testID := con.cache.GetTestID("TestID")
	testIDStr := strconv.Itoa(int(testID))
	questionKey := strconv.Itoa(queNum) + "_" + testIDStr
	answer := con.cache.GetQuestion(questionKey).Answer
	if answerSubmitted == answer {
		scoreKey := "Score_" + testIDStr
		con.cache.UpdateScore(scoreKey)
	}
}

func (con *testController) UpdateTest() int {
	testID := con.cache.GetTestID("TestID")
	scoreKey := "Score_" + strconv.Itoa(int(testID))
	score := con.cache.GetScore(scoreKey)
	con.service.UpdateTest(testID, score)
	con.cache.FlushAll()
	return score
}
