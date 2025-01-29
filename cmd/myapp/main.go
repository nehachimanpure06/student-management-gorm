package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"student-management-gorm/config"
	"student-management-gorm/pkg/controller"
	"student-management-gorm/pkg/repository/sqldb"
	"student-management-gorm/pkg/service"
	"student-management-gorm/routes"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetOutput(os.Stdout)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)

	log.Info("Welcome to the student management application...")
	log.Info("starting the server")
	defer log.Warn("Exiting the server..")

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	r := gin.Default()

	studentController, courseController := setupDependancies(cfg)
	routes.InitRoutes(r, studentController, courseController)

	server := http.Server{
		Addr:        cfg.HTTPServer.Port,
		Handler:     r,
		ReadTimeout: time.Duration(cfg.HTTPServer.ReadTimeout) * time.Second,
	}
	// err := r.Run(":8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Panic("error occured while starting the server : " + err.Error())
	}
}

func setupDependancies(cfg *config.Config) (*controller.StudentController, *controller.CourseController) {
	var db *gorm.DB

	if cfg.Env == "production" {
		db = newPostgresDBConnection(cfg)
	} else if cfg.Env == "development" {
		db = newMySQLDBConnection(cfg)
	}

	studentRepo := sqldb.NewSQLStudentRepository(db)
	if err := studentRepo.Migrate(); err != nil {
		log.Fatalf("Error during student table migration: %v\n", err)
	} else {
		log.Info("student table migration successful.")
	}

	courseRepo := sqldb.NewSQLCourseRepository(db)
	if err := courseRepo.Migrate(); err != nil {
		log.Fatalf("Error during course table migration: %v\n", err)
	} else {
		log.Info("course table migration successful.")
	}

	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo)

	studentController := controller.NewStudentController(studentService)
	courseController := controller.NewCourseController(courseService)

	return studentController, courseController
}

func newMySQLDBConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.MysqlConfig.Username, cfg.MysqlConfig.Password,
		cfg.MysqlConfig.HostName, cfg.MysqlConfig.Port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error occured while initializing database : ", err.Error())
	}

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.MysqlConfig.Database)
	err = db.Exec(query).Error
	if err != nil {
		log.Fatal("error occured while creating database : ", err.Error())
	}

	err = db.Exec(fmt.Sprintf("USE %s", cfg.MysqlConfig.Database)).Error
	if err != nil {
		log.Fatal("error occured while connecting with database : ", err.Error())
	}

	return db
}

func newPostgresDBConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("user=%s password=%s host=%s sslmode=%s", cfg.PostgreSQLConfig.Username,
		cfg.PostgreSQLConfig.Password, cfg.PostgreSQLConfig.HostName, cfg.PostgreSQLConfig.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error occured while initializing database : ", err.Error())
	}

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.PostgreSQLConfig.Database)
	err = db.Exec(query).Error
	if err != nil {
		log.Fatal("error occured while creating database : ", err.Error())
	}

	err = db.Exec(fmt.Sprintf("USE %s", cfg.PostgreSQLConfig.Database)).Error
	if err != nil {
		log.Fatal("error occured while connecting with database : ", err.Error())
	}

	return db
}
