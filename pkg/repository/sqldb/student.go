package sqldb

import (
	"errors"
	"fmt"
	"student-management-gorm/pkg/model"

	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StudentDB struct {
	ID             int    `gorm:"primaryKey;autoIncrement;column:id"`
	FirstName      string `gorm:"column:first_name"`
	LastName       string `gorm:"column:last_name"`
	Email          string `gorm:"unique;column:email"`
	Phone          string `gorm:"column:phone"`
	DateOfBirth    string `gorm:"column:date_of_birth;type:date"`
	EnrollmentDate string `gorm:"column:enrollment_date;type:date"`
	Status         string `gorm:"column:status;type:ENUM('Active', 'Graduated', 'Dropped')"` // Active, Graduated, Dropped
}

type SQLStudentRepository struct {
	DB *gorm.DB
}

func NewSQLStudentRepository(db *gorm.DB) *SQLStudentRepository {
	return &SQLStudentRepository{DB: db}
}

func (StudentDB) TableName() string {
	return "students"
}

func (repo *SQLStudentRepository) Migrate() error {
	err := repo.DB.AutoMigrate(&StudentDB{})
	if err != nil {
		return err
	}
	return nil
}

func (repo *SQLStudentRepository) toDBModel(student model.Student) StudentDB {
	return StudentDB{
		ID:             student.ID,
		FirstName:      student.FirstName,
		LastName:       student.LastName,
		Email:          student.Email,
		Phone:          student.Phone,
		DateOfBirth:    student.DateOfBirth,
		EnrollmentDate: student.EnrollmentDate,
		Status:         student.Status,
	}
}

func (repo *SQLStudentRepository) toModel(student StudentDB) model.Student {
	return model.Student{
		ID:             student.ID,
		FirstName:      student.FirstName,
		LastName:       student.LastName,
		Email:          student.Email,
		Phone:          student.Phone,
		DateOfBirth:    student.DateOfBirth,
		EnrollmentDate: student.EnrollmentDate,
		Status:         student.Status,
	}
}

func (repo *SQLStudentRepository) GetAllStudents() ([]model.Student, error) {
	var result []StudentDB
	var students []model.Student
	err := repo.DB.Find(&result).Error
	if err != nil {
		return nil, errors.New("error occured while getting student data : " + err.Error())
	}

	for _, std := range result {
		students = append(students, repo.toModel(std))
	}
	return students, nil
}

func (repo *SQLStudentRepository) AddStudent(student model.Student) (int, error) {
	dbStudent := repo.toDBModel(student)
	err := repo.DB.Create(&dbStudent).Error
	if err != nil {
		return 0, err
	}
	return int(student.ID), nil
}

func (repo *SQLStudentRepository) UpdateStudent(studentID int, student model.Student) error {
	dbStudent := repo.toDBModel(student)
	err := repo.DB.Model(&StudentDB{}).Where("id = ?", studentID).Updates(dbStudent).Error
	if err != nil {
		return fmt.Errorf("could not update student: %v", err)
	}

	return nil
}

func (repo *SQLStudentRepository) DeleteStudent(studentID int) error {
	err := repo.DB.Where("id = ?", studentID).Delete(&StudentDB{}).Error
	if err != nil {
		return fmt.Errorf("could not delete student: %v", err)
	}

	return nil
}

func (repo *SQLStudentRepository) GetStudentByID(studentId int) (model.Student, error) {
	var student StudentDB

	err := repo.DB.First(&student, studentId).Error
	if err != nil {
		return model.Student{}, errors.New("error occurred while getting student data: " + err.Error())
	}
	return repo.toModel(student), nil
}
