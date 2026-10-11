package controllers

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Education level CRUD
func (ctrl *JobController) CreateEducationLevelHandler(c *gin.Context) {
	var item models.EducationLevel
	// in this create handlers this binds all the attributes for the model. including id and the created_at which shouldnt taken from the client.
	// we can avoid this with binding an input struct
	// SERIAL only advances its sequence when the id comes from the column default (nextval).
	// an explicit id from the client skips nextval, so the sequence falls behind and later
	// creates get an id that already exists, failing with a duplicate key error (500).

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateEducationLevel(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create education level", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllEducationLevelsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllEducationLevels()
	if err != nil {
		// here returning err.Error() is problematic since it reveals the db driver internal states.
		// better option is log the error and return a generic message to the user
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch education levels", "details": err.Error()})
		return
	}
	// here this works for now but gin.H is a map of string to any.
	// there is no any type validation before it sends the response.
	// there are 2 ways to do that one is create a generic response struct in somewhere (response.go) and pass it here
	// or you can create separate struct for that.
	// generic struct would be the better option but that will break the frontend since these responses use specific keys for each models.
	// here is the sample struct for the generic list response

	// internal/api/response.go
	// type ListResponse[T any] struct {
	// 	Count int `json:"count"`
	// 	Items []T `json:"items"`
	// }

	c.JSON(http.StatusOK, gin.H{"count": len(items), "education_levels": items})
}

func (ctrl *JobController) GetEducationLevelByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	// use strconv.ParseUint instead.
	// check UpdateEducationLevelHandler comments for better explanation
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetEducationLevelByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Education level not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch education level", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateEducationLevelHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	// this is problematic. this Sscanf function reads the leading digits and ignores the rest.
	// ex: "12abc" becomes 12.
	// better option i see here is strconv.ParseUint(idStr, 10, 64) this will reject the whole string if it isn't a number.
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		// we can print the Id with the error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	// map[string]interface{} (interface{} == any) accepts any JSON object, and it is passed
	// straight to GORM's Updates. unknown keys cause a 500, and worse, the client can write
	// columns it should never control (id, created_at, deleted_at) without any error (mass assignment).
	// better to bind to a struct with only the editable fields and validation, ex:
	//   var in struct {
	//       Level string `json:"level" binding:"required,max=100"`
	//   }
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	item, err := ctrl.repo.UpdateEducationLevel(id, payload)
	// when the id doesn't exist, UpdateEducationLevel returns gorm.ErrRecordNotFound (from First), but this
	// handler maps every error to 500, so the client gets 500 instead of 404. check
	// errors.Is(err, gorm.ErrRecordNotFound) and return 404, as GetEducationLevelByIDHandler already does.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update education level", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteEducationLevelHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	// consider the other comments above
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteEducationLevel(id); err != nil {
		// consider the other comments above
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete education level", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Education level deleted successfully"})
}
