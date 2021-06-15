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
	Chapter_ID uint   `json:"chapter_id"`
	Class      uint   `json:"class"`
	Subject    uint   `json:"subject"` // phy:1, chem:2, math:3, bio:4
	Chapter    string `json:"chapter"`
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

//Signup user
// func Signup(c *gin.Context) {
// 	var user User
// 	c.BindJSON(&user)
// 	var existinguser User
// 	if err := db.Where("username = ?", user.Username).First(&existinguser).Error; err == nil {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{user.Username: "already exists. Try another"})
// 	} else if err := db.Where("emailid = ?", user.Emailid).First(&existinguser).Error; err == nil {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(202, gin.H{user.Emailid: "already exists. Try another"})
// 	} else {
// 		user.Points = 0
// 		user.Role = 0
// 		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
// 		user.Password = string(hashedPwd)
// 		db.Create(&user)
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, user)
// 	}
// }

// //Signin user
// func Signin(c *gin.Context) {
// 	var user User
// 	var existinguser User
// 	c.BindJSON(&user)
// 	if err := db.Where("username = ?", user.Username).First(&existinguser).Error; err != nil {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{user.Username: "does not exist"})
// 	} else {
// 		if err = bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte(user.Password)); err != nil {
// 			c.Header("access-control-allow-origin", "*")
// 			c.JSON(202, gin.H{user.Username: "incorrect password"})
// 		} else {
// 			c.Header("access-control-allow-origin", "*")
// 			c.JSON(200, existinguser)
// 		}
// 	}
// }

// //GetUserPoints for any user
// func GetUserPoints(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user Points
// 	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, user)
// 	}
// }

// // -----------------------------------------------------admin functions--------------------------------------------------------------//

// //GetPeople for admin
// func GetPeople(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		var people []User
// 		if err := db.Find(&people).Error; err != nil {
// 			c.AbortWithStatus(404)
// 			fmt.Println(err)
// 		} else {
// 			c.Header("access-control-allow-origin", "*")
// 			c.JSON(200, people)
// 		}
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// //DeleteUser for admin
// func DeleteUser(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		id := c.Params.ByName("id")
// 		var user User
// 		d := db.Where("id = ?", id).Delete(&user)
// 		fmt.Println(d)
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, gin.H{"id #" + id: "deleted"})
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// //CreateQuiz given genre, quiz_num, num_questions for admin
// func CreateQuiz(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		var quiz Quiz
// 		c.BindJSON(&quiz)
// 		quiz.Num_Questions = 0
// 		var existingquiz Quiz
// 		if e := db.Where("genre = ?", quiz.Genre).Where("quiz_num = ?", quiz.Quiz_Num).First(&existingquiz).Error; e == nil {
// 			c.Header("access-control-allow-origin", "*")
// 			c.JSON(202, existingquiz)
// 		} else {
// 			db.Create(&quiz)
// 			var genre Genre
// 			if err := db.Where("genre_name = ?", quiz.Genre).First(&genre).Error; err != nil {
// 				// create genre
// 				genre.Genre_Name = quiz.Genre
// 				genre.Num_Quizzes = quiz.Quiz_Num
// 				db.Create(&genre)
// 			} else {
// 				// update genre
// 				genre.Num_Quizzes = quiz.Quiz_Num
// 				db.Save(&genre)
// 			}
// 		}
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, quiz)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// //DeleteQuiz given genre, quiz_num for admin
// func DeleteQuiz(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		id := c.Params.ByName("id")
// 		var quiz Quiz
// 		db.Where("id = ?", id).First(&quiz)
// 		genre_name := quiz.Genre
// 		quiz_num := quiz.Quiz_Num
// 		num_questions := quiz.Num_Questions
// 		d := db.Where("genre = ?", genre_name).Where("quiz_num = ?", quiz_num).Delete(&quiz)

// 		var genre Genre
// 		db.Where("genre_name = ?", genre_name).First(&genre)
// 		if genre.Num_Quizzes > 0 {
// 			genre.Num_Quizzes--
// 		}
// 		db.Save(&genre)

// 		if num_questions > 0 {
// 			var question []Question
// 			db.Where("genre = ?", genre_name).Where("quiz_num = ?", quiz_num).Delete(&question)
// 		}

// 		fmt.Println(d)
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, quiz)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// // CreateQuestion for an existing quiz for admin
// func CreateQuestion(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		var question Question
// 		c.BindJSON(&question)
// 		db.Create(&question)

// 		var quiz Quiz
// 		db.Where("Genre = ?", question.Genre).Where("Quiz_Num = ?", question.Quiz_Num).First(&quiz)
// 		quiz.Num_Questions++
// 		db.Save(&quiz)

// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, question)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// // DeleteQuestion given genre, quiz_num, question_num for admin
// func DeleteQuestion(c *gin.Context) {
// 	username := c.Params.ByName("username")
// 	var user User
// 	db.Where("username = ?", username).First(&user)
// 	if user.Role == 1 {
// 		id := c.Params.ByName("id")
// 		var question Question
// 		var q Question
// 		db.Where("id = ?", id).First(&q)
// 		d := db.Where("id = ?", id).Delete(&question)
// 		genre := q.Genre
// 		quiz_num := q.Quiz_Num
// 		var quiz Quiz
// 		db.Where("genre = ?", genre).Where("quiz_num = ?", quiz_num).First(&quiz)
// 		quiz.Num_Questions--
// 		db.Save(&quiz)
// 		fmt.Println(d)
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, q)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(201, gin.H{username: "You're not an admin"})
// 	}
// }

// //GetQuestion for viewing
// func GetQuestion(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	var question Question
// 	db.Where("id = ?", id).First(&question)
// 	c.Header("access-control-allow-origin", "*")
// 	c.JSON(200, question)
// }

// //UpdateQuestion while editing
// func UpdateQuestion(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	var question Question
// 	if err := db.Where("id = ?", id).First(&question).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	}
// 	c.BindJSON(&question)
// 	db.Save(&question)
// 	c.Header("access-control-allow-origin", "*")
// 	c.JSON(200, question)
// }

// // --------------------------------------------------admin functions over-------------------------------------------------------------//

// //GetQuizDetails for a quiz
// func GetQuizDetails(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	var quiz Quiz
// 	db.Where("id = ?", id).First(&quiz)
// 	c.Header("access-control-allow-origin", "*")
// 	c.JSON(200, quiz)
// }

// //GetAllQuiz to get all the quizes
// func GetAllQuiz(c *gin.Context) {
// 	var quiz []Quiz
// 	if err := db.Find(&quiz).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, quiz)
// 	}
// }

// //GetGenres to retrieve all genres
// func GetGenres(c *gin.Context) {
// 	var genre []Genre
// 	if err := db.Find(&genre).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, genre)
// 	}
// }

// // GetNumQuizzes using genre to get number of quizzes in each genre
// func GetNumQuizzes(c *gin.Context) {
// 	genre := c.Params.ByName("genre")
// 	var quizzes Genre
// 	if err := db.Where("Genre_Name = ?", genre).First(&quizzes).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, quizzes)
// 	}
// }

// // GetQuiz for a quiz retirves all the questions for a quiz
// func GetQuiz(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	var quiz Quiz
// 	db.Where("id = ?", id).First(&quiz)
// 	genre := quiz.Genre
// 	quiz_num := quiz.Quiz_Num
// 	var question []Question
// 	if err := db.Where("genre = ?", genre).Where("quiz_num = ?", quiz_num).Find(&question).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		c.JSON(200, question)
// 	}
// }

// //GetLeaderboard across all genres
// func GetLeaderboard(c *gin.Context) {
// 	var user []User
// 	if err := db.Find(&user).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		sort.SliceStable(user, func(i, j int) bool {
// 			return user[i].Points > user[j].Points
// 		})
// 		c.JSON(200, user)
// 	}
// }

// //GetLeaderboardByGenre for a particular genre
// func GetLeaderboardByGenre(c *gin.Context) {
// 	var user []Points
// 	genre := c.Params.ByName("genre")
// 	if err := db.Where("genre = ?", genre).Find(&user).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		sort.SliceStable(user, func(i, j int) bool {
// 			return user[i].Points > user[j].Points
// 		})
// 		c.JSON(200, user)
// 	}
// }

// // GetHistory function
// func GetHistory(c *gin.Context) {
// 	var points []History
// 	username := c.Params.ByName("username")
// 	if err := db.Where("username = ?", username).Find(&points).Error; err != nil {
// 		c.AbortWithStatus(404)
// 		fmt.Println(err)
// 	} else {
// 		c.Header("access-control-allow-origin", "*")
// 		sort.SliceStable(points, func(i, j int) bool {
// 			return points[i].Timestamp > points[j].Timestamp
// 		})
// 		c.JSON(200, points)
// 	}
// }

// // EvaluateQuiz function
// func EvaluateQuiz(c *gin.Context) {
// 	// update total user points, points table, history
// 	id := c.Params.ByName("id") //quiz id
// 	username := c.Params.ByName("username")
// 	p := c.Params.ByName("points")
// 	points, _ := strconv.ParseUint(p, 10, 64)
// 	currentTime := time.Now()
// 	var user User
// 	var quiz Quiz
// 	var history History
// 	var point Points

// 	db.Where("username = ?", username).First(&user)
// 	db.Where("id = ?", id).First(&quiz)

// 	user.Points += points
// 	db.Save(&user)

// 	history.Username = user.Username
// 	history.Genre = quiz.Genre
// 	history.Quiz_Num = quiz.Quiz_Num
// 	history.Score = points
// 	history.Timestamp = currentTime.Format("2006-01-02 15:04:05")
// 	db.Save(&history)

// 	if err := db.Where("genre = ?", quiz.Genre).First(&point).Error; err != nil {
// 		//create user in points table
// 		point.Username = user.Username
// 		point.Genre = quiz.Genre
// 		point.Points = points
// 		db.Create(&point)
// 	} else {
// 		point.Points += points
// 		db.Save(&point)
// 	}
// 	c.Header("access-control-allow-origin", "*")
// 	c.JSON(200, point)
// }
