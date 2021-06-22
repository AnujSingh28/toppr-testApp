package model

import (
	"fmt"
	"net/http"
	"test-app/Config"
	"test-app/entity"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type TestModel interface {
	AllClasses() []entity.Classes
	AllSubject(class entity.Classes) []entity.Subjects
	FetchChapter(subject entity.Subjects) []entity.Chapters
	FetchQuestions(chapter entity.Chapters) ([]entity.Questions, uint64)
	UpdateTest(testID uint64, score int)
	Signup(ctx *gin.Context, user entity.User) entity.User
	Login(ctx *gin.Context, user entity.User) entity.User
}

type database struct {
	connection *gorm.DB
}

func NewTestModel() TestModel {
	db, err := gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.User{}, &entity.Classes{}, &entity.Subjects{}, &entity.Chapters{}, &entity.Questions{}, &entity.Tests{},
		&entity.ClassSubs{}, &entity.SubChaps{}, &entity.ChapQues{}, &entity.ChapTests{})
	if err != nil {
		panic("Failed to Automigrate")
	}
	return &database{
		connection: db,
	}
}

func (db *database) Signup(ctx *gin.Context, user entity.User) entity.User {
	fmt.Println(ctx.Params.ByName("username"))
	if user.FirstName == "" || user.LastName == "" || user.Password == "" || user.Username == "" {
		ctx.AbortWithStatus(400)
	}
	rowsAffected := db.connection.Create(&user).RowsAffected
	if rowsAffected == 0 {
		panic("Error while adding user")
	}
	ctx.Header("access-control-allow-origin", "*")
	ctx.JSON(http.StatusOK, user)
	return user
}

func (db *database) Login(ctx *gin.Context, user entity.User) entity.User {
	var tempUser entity.User
	if err := db.connection.Where("username = ?", user.Username).First(&tempUser).Error; err != nil {
		ctx.Header("access-control-allow-origin", "*")
		ctx.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		if tempUser.Password == user.Password {
			fmt.Println("Authenticated")
			ctx.Header("access-control-allow-origin", "*")
			ctx.JSON(http.StatusOK, user)
		} else {
			ctx.Header("access-control-allow-origin", "*")
			ctx.JSON(401, gin.H{"authenticated": false})
		}
	}
	return user
}

func (db *database) AllClasses() []entity.Classes {
	var classes []entity.Classes
	rowsAffected := db.connection.Raw("select * from classes").Scan(&classes).RowsAffected
	if rowsAffected == 0 {
		panic("Error while fetching classes")
	}
	return classes
}

func (db *database) AllSubject(class entity.Classes) []entity.Subjects {
	var subjects []entity.Subjects
	rowsAffected := db.connection.Raw("select * from subjects inner join class_subs on id = subject_id where class_id = ?", class.ID).Scan(&subjects).RowsAffected
	if rowsAffected == 0 {
		panic("Error while fetching subjects")
	}
	return subjects
}

func (db *database) FetchChapter(subject entity.Subjects) []entity.Chapters {
	var chapters []entity.Chapters
	rowsAffected := db.connection.Raw("select * from chapters inner join sub_chaps on id = chapter_id where subject_id = ?", subject.ID).Scan(&chapters).RowsAffected
	if rowsAffected == 0 {
		panic("Error while fetching chapters")
	}
	return chapters
}

func (db *database) FetchQuestions(chapter entity.Chapters) ([]entity.Questions, uint64) {
	var questions []entity.Questions
	rowsAffected := db.connection.Raw("select * from questions inner join chap_ques on id = question_id where chapter_id = ? limit 10", chapter.ID).Scan(&questions).RowsAffected
	if rowsAffected != 10 {
		panic("Error while fetching questions")
	}
	tx := db.connection.Begin()
	for _, que := range questions {
		UpdateQuestion(que, tx)
	}
	testID := AddTest(chapter, tx)
	tx.Commit()
	return questions, testID
}

func UpdateQuestion(question entity.Questions, tx *gorm.DB) {
	rowsAffected := tx.Model(&entity.Questions{}).Where("id = ?", question.ID).Update("times_asked", question.TimesAsked+1).RowsAffected
	if rowsAffected != 1 {
		tx.Rollback()
		panic("Error while updating times_asked count in questions table")
	}
}

func AddTest(chapter entity.Chapters, tx *gorm.DB) uint64 {
	var testID uint64
	tx.Raw("insert into tests(score) values(0)")
	rowsAffected := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&testID).RowsAffected
	if rowsAffected != 1 {
		tx.Rollback()
		panic("Error while adding entry in tests table")
	}
	record := entity.ChapTests{ChapterID: chapter.ID, TestID: testID}
	rowsAffected = tx.Create(&record).RowsAffected
	if rowsAffected != 1 {
		tx.Rollback()
		panic("Error while adding entry in chap_tests table")
	}
	return testID
}

func (db *database) UpdateTest(testID uint64, score int) {
	rowsAffected := db.connection.Model(&entity.Tests{}).Where("id = ?", testID).Update("score", score).RowsAffected
	if rowsAffected != 1 {
		panic("Error while updating test score")
	}
}
