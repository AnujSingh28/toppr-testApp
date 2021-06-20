package entity

type Classes struct {
	ID   uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Name string `json:"name" binding:"min=1,max=100" gorm:"type:varchar(100)"`
}

type Subjects struct {
	ID   uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Name string `json:"name" binding:"min=2,max=100" gorm:"type:varchar(100)"`
}

type Chapters struct {
	ID   uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Name string `json:"name" binding:"min=2,max=100" gorm:"type:varchar(100)"`
}

type Questions struct {
	ID         uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Question   string `json:"question" binding:"required" gorm:"type:text"`
	OptionA    string `json:"optiona" binding:"required" gorm:"type:text"`
	OptionB    string `json:"optionb" binding:"required" gorm:"type:text"`
	OptionC    string `json:"optionc" binding:"required" gorm:"type:text"`
	OptionD    string `json:"optiond" binding:"required" gorm:"type:text"`
	Answer     int    `json:"answer" binding:"min=1,max=1" gorm:"type:int"`
	TimesAsked int    `json:"timesasked" gorm:"index;type:int"`
}

type Tests struct {
	ID    uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Score int8   `json:"score" gorm:"type:int"`
}

type ClassSubs struct {
	Class     Classes  `json:"-" binding:"-" gorm:"foreignkey:ClassID"`
	ClassID   uint64   `json:"-" binding:"-"`
	Subject   Subjects `json:"-" binding:"-" gorm:"foreignkey:SubjectID"`
	SubjectID uint64   `json:"-" binding:"-"`
}

type SubChaps struct {
	Subject   Subjects `json:"-" binding:"-" gorm:"foreignkey:SubjectID"`
	SubjectID uint64   `json:"-" binding:"-"`
	Chapter   Chapters `json:"-" binding:"-" gorm:"foreignkey:ChapterID"`
	ChapterID uint64   `json:"-" binding:"-"`
}

type ChapQues struct {
	Chapter    Chapters  `json:"-" binding:"-" gorm:"foreignkey:ChapterID"`
	ChapterID  uint64    `json:"-" binding:"-"`
	Question   Questions `json:"-" binding:"-" gorm:"foreignkey:QuestionID"`
	QuestionID uint64    `json:"-" binding:"-"`
}

type ChapTests struct {
	Chapter   Chapters `json:"-" binding:"-" gorm:"foreignkey:ChapterID"`
	ChapterID uint64   `json:"-" binding:"-"`
	Test      Tests    `json:"-" binding:"-" gorm:"foreignkey:TestID"`
	TestID    uint64   `json:"-" binding:"-"`
}
