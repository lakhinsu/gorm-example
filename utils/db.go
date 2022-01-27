package utils

import (
	"fmt"
	"strconv"

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

func CreateDBConnection() (*gorm.DB, error) {
	log.Info().Msg(GetDSN())
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	return db, err
}

func AutoMigrateDB() error {
	// Auto migrate database
	db, connErr := CreateDBConnection()
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
