package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/lakhinsu/gorm-example/models"
	"github.com/lakhinsu/gorm-example/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var user models.User

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindJSON(&user); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}
	user.FillDefaults()

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	// Create a user
	result := db.Create(&user)
	if result.Error != nil && result.RowsAffected != 1 {
		log.Err(result.Error).Str("request_id", request_id).Msg("Error occurred while creating a new user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating a new user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"id":      user.ID,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindJSON(&user); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}
	user.FillDefaults()

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	// Create a user object
	var value models.User

	// Read the user which is to be updated
	result := db.First(&value, "id = ?", user.ID)
	if result.Error != nil {
		log.Err(result.Error).Str("request_id", request_id).Msg("Error occurred while updating the user")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while updating user",
			})
		}
		return
	}

	// Update the desired values using the request payload
	value.FirstName = user.FirstName
	value.LastName = user.LastName

	// Save the updated user
	tx := db.Save(&value)
	if tx.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while updating user",
		})
		return
	}

	// Return the updated user with the response
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"result":  value,
	})

}

func GetUser(c *gin.Context) {
	var userId models.UserID

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindUri(&userId); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	var user models.User
	result := db.First(&user, "id = ?", userId.ID)
	if result.Error != nil {
		log.Err(result.Error).Str("request_id", request_id).Msg("Error occurred while fetching the user")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while fetching user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	request_id := c.GetString("x-request-id")

	earliest := c.DefaultQuery("earliest", "0")
	latest := c.DefaultQuery("latest", fmt.Sprint(time.Now().UnixMilli()))

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	tx := db.Where(fmt.Sprintf("created_at >= '%s' and created_at <= '%s'", earliest, latest)).Order("created_at asc").Scopes(utils.Paginate(c)).Find(&users)

	if tx.RowsAffected == 0 {
		log.Info().Msg("Read users returned with empty results")
	}
	c.JSON(http.StatusOK, gin.H{
		"earliest": earliest,
		"latest":   latest,
		"results":  users,
	})
}

func DeleteUser(c *gin.Context) {
	var userId models.UserID

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindUri(&userId); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	var user models.User
	result := db.First(&user, "id = ?", userId.ID)
	if result.Error != nil {
		log.Err(result.Error).Str("request_id", request_id).Msg("Error occurred while deleting the user")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while deleting user",
			})
		}
		return
	}

	tx := db.Delete(&user)
	if tx.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while deleting user",
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
