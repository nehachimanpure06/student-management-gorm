package sqldb

import (
	"errors"
	"fmt"

	"student-management-gorm/pkg/model"

	"gorm.io/gorm"
)

type SQLCourseRepository struct {
	DB *gorm.DB
}

type CourseDB struct {
	ID             int    `gorm:"primaryKey;autoIncrement;column:id"`
	Name           string `gorm:"column:name"`
	Description    string `gorm:"column:description"`
	Credits        int    `gorm:"column:credits"`
	Instructor     string `gorm:"column:instructor"`
	Schedule       string `gorm:"column:schedule"`
	Capacity       int    `gorm:"column:capacity"`
	AvailableSeats int    `gorm:"column:availabile_seats"`
}

func (CourseDB) TableName() string {
	return "courses"
}

func NewSQLCourseRepository(db *gorm.DB) *SQLCourseRepository {
	return &SQLCourseRepository{DB: db}
}

func (repo *SQLCourseRepository) toDBModel(course model.Course) CourseDB {
	return CourseDB{
		ID:             course.ID,
		Name:           course.Name,
		Description:    course.Description,
		Credits:        course.Credits,
		Instructor:     course.Instructor,
		Schedule:       course.Schedule,
		Capacity:       course.Capacity,
		AvailableSeats: course.AvailableSeats,
	}
}

func (repo *SQLCourseRepository) toModel(course CourseDB) model.Course {
	return model.Course{
		ID:             course.ID,
		Name:           course.Name,
		Description:    course.Description,
		Credits:        course.Credits,
		Instructor:     course.Instructor,
		Schedule:       course.Schedule,
		Capacity:       course.Capacity,
		AvailableSeats: course.AvailableSeats,
	}
}

func (repo *SQLCourseRepository) Migrate() error {
	if err := repo.DB.AutoMigrate(&CourseDB{}); err != nil {
		return err
	}
	return nil
}

func (repo *SQLCourseRepository) GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	var result []CourseDB
	err := repo.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}

	for _, crs := range result {
		courses = append(courses, repo.toModel(crs))
	}

	return courses, nil
}

func (repo *SQLCourseRepository) AddCourse(course model.Course) (int, error) {
	dbCourse := repo.toDBModel(course)
	err := repo.DB.Create(&dbCourse).Error
	if err != nil {
		return 0, err
	}
	return int(dbCourse.ID), nil
}

func (repo *SQLCourseRepository) UpdateCourse(courseID int, course model.Course) error {
	courseDB := repo.toDBModel(course)
	err := repo.DB.Model(&CourseDB{}).Where("id = ?", courseID).Updates(courseDB).Error
	if err != nil {
		return fmt.Errorf("could not update course: %v", err)
	}
	return nil
}

func (repo *SQLCourseRepository) DeleteCourse(courseID int) error {
	err := repo.DB.Delete(courseID).Error
	if err != nil {
		return fmt.Errorf("could not delete course: %v", err)
	}
	return nil
}

func (repo *SQLCourseRepository) GetCourseByID(courseID int) (model.Course, error) {
	var course CourseDB
	err := repo.DB.First(&course, courseID).Error
	if err != nil {
		return model.Course{}, errors.New("error occurred while getting course data: " + err.Error())
	}
	return repo.toModel(course), nil
}
