package utils

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	models "github.com/lakhinsu/gorm-example/models"
	"github.com/rs/zerolog/log"
)

var user string
var password string
var db string
var host string
var port string
var ssl string
var timezone string
var dbConn *gorm.DB

func init() {
	user = GetEnvVar("POSTGRES_USER")
	password = GetEnvVar("POSTGRES_PASSWORD")
	db = GetEnvVar("POSTGRES_DB")
	host = GetEnvVar("POSTGRES_HOST")
	port = GetEnvVar("POSTGRES_PORT")
	ssl = GetEnvVar("POSTGRES_SSL")
	timezone = GetEnvVar("POSTGRES_TIMEZONE")
}

func GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, db, port, ssl, timezone)
}

func CreateDBConnection() error {
	// Close the existing connection if open
	if dbConn != nil {
		CloseDBConnection(dbConn)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Err(err).Msg("Error occurred while connecting with the database")
	}

	sqlDB, err := db.DB()

	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db
	return err
}

func GetDatabaseConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return dbConn, err
	}
	if err := sqlDB.Ping(); err != nil {
		return dbConn, err
	}
	return dbConn, nil
}

func CloseDBConnection(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.Err(err).Msg("Error occurred while closing a DB connection")
	}
	defer sqlDB.Close()
}

func AutoMigrateDB() error {
	// Auto migrate database
	db, connErr := GetDatabaseConnection()
	if connErr != nil {
		return connErr
	}
	// Add new models here
	err := db.AutoMigrate(&models.User{})
	return err
}

// Pagination helper for GORM
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
