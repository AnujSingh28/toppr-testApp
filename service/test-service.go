package service

import (
	"test-app/entity"
	"test-app/model"
)

type TestService interface {
	AllClasses() []entity.Classes
	AllSubject(class entity.Classes) []entity.Subjects
	FetchChapter(entity.Subjects) []entity.Chapters
	FetchQuestions(entity.Chapters) ([]entity.Questions, uint64)
	UpdateTest(testID uint64, score int)
}

type testService struct {
	testModel model.TestModel
}

func New(mod model.TestModel) TestService {
	return &testService{
		testModel: mod,
	}
}

func (service *testService) AllClasses() []entity.Classes {
	return service.testModel.AllClasses()
}

func (service *testService) AllSubject(class entity.Classes) []entity.Subjects {
	return service.testModel.AllSubject(class)
}

func (service *testService) FetchChapter(subject entity.Subjects) []entity.Chapters {
	return service.testModel.FetchChapter(subject)
}

func (service *testService) FetchQuestions(chapter entity.Chapters) ([]entity.Questions, uint64) {
	return service.testModel.FetchQuestions(chapter)
}

func (service *testService) UpdateTest(testID uint64, score int) {
	service.testModel.UpdateTest(testID, score)
}
