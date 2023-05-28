package controllers

import (
	"net/http"
	"strconv"

	"github.com/khasyah-fr/sky-todo/entities"
	"github.com/khasyah-fr/sky-todo/entities/activities"
	"github.com/khasyah-fr/sky-todo/helpers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ActivityController struct {
	db *gorm.DB
}

func NewActivityController(db *gorm.DB) *ActivityController {
	return &ActivityController{db: db}
}

// @Summary Get all activities
// @Description Retrieves all activities
// @Produce json
// @Success 200 {array} entities.GetActivityResponse
// @Router /activities [get]
func (ac *ActivityController) GetAllActivities(c echo.Context) error {
	db := ac.db
	response := entities.Response[[]entities.GetActivityResponse]{}

	activities := []activities.Activity{}
	if err := db.Find(&activities).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	for _, activity := range activities {
		response.Data = append(response.Data, entities.GetActivityResponse{
			ID:        activity.ActivityID,
			Title:     activity.Title,
			Email:     activity.Email,
			CreatedAt: activity.CreatedAt,
			UpdatedAt: activity.UpdatedAt,
		})
	}

	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS

	return c.JSON(http.StatusOK, response)
}

// @Summary Get an activity
// @Description Retrieves a specific activity by ID
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} entities.GetActivityResponse
// @Router /activities/{id} [get]
func (ac *ActivityController) GetActivity(c echo.Context) error {
	db := ac.db
	response := entities.Response[entities.GetActivityResponse]{}
	id, _ := strconv.Atoi(c.Param("id"))

	activity := activities.Activity{}
	condition := activities.Activity{ActivityID: id}
	if err := db.Where(&condition).Find(&activity).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if activity == (activities.Activity{}) {
		response.Status = helpers.NOT_FOUND
		response.Message = "Activity with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	response.Data = entities.GetActivityResponse{
		ID:        activity.ActivityID,
		Title:     activity.Title,
		Email:     activity.Email,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS

	return c.JSON(http.StatusOK, response)
}

// @Summary Create an activity
// @Description Creates a new activity
// @Accept json
// @Produce json
// @Param activity body entities.CreateActivityRequest true "Activity object"
// @Success 201 {object} entities.GetActivityResponse
// @Router /activities [post]
func (ac *ActivityController) CreateActivity(c echo.Context) error {
	db := ac.db
	response := entities.Response[entities.GetActivityResponse]{}
	request := entities.CreateActivityRequest{}

	if err := c.Bind(&request); err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if request.Title == "" {
		response := entities.Response[[]string]{}
		response.Data = make([]string, 0)
		response.Status = helpers.ERROR_BAD_REQUEST
		response.Message = "title cannot be null"
		return c.JSON(http.StatusBadRequest, response)
	}

	activity := activities.Activity{
		Title: request.Title,
		Email: request.Email,
	}

	if err := db.Create(&activity).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.GetActivityResponse{
		ID:        activity.ActivityID,
		Title:     activity.Title,
		Email:     activity.Email,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS

	return c.JSON(http.StatusCreated, response)
}

// @Summary Update an activity
// @Description Updates an existing activity
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param activity body entities.UpdateActivityRequest true "Updated activity object"
// @Success 200 {object} entities.GetActivityResponse
// @Router /activities/{id} [put]
func (ac *ActivityController) UpdateActivity(c echo.Context) error {
	db := ac.db
	response := entities.Response[entities.GetActivityResponse]{}
	id, _ := strconv.Atoi(c.Param("id"))
	request := entities.UpdateActivityRequest{}

	if err := c.Bind(&request); err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if request.Title == "" {
		response := entities.Response[[]string]{}
		response.Data = make([]string, 0)
		response.Status = helpers.ERROR_BAD_REQUEST
		response.Message = "title cannot be null"
		return c.JSON(http.StatusBadRequest, response)
	}

	activity := activities.Activity{}
	condition := activities.Activity{ActivityID: id}
	if err := db.Where(&condition).Find(&activity).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if activity == (activities.Activity{}) {
		response.Status = helpers.NOT_FOUND
		response.Message = "Activity with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	if err := db.Model(&activity).Updates(activities.Activity{Title: request.Title}).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.GetActivityResponse{
		ID:        activity.ActivityID,
		Title:     activity.Title,
		Email:     activity.Email,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS

	return c.JSON(http.StatusOK, response)
}

// @Summary Delete an activity
// @Description Deletes an activity by ID
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200
// @Router /activities/{id} [delete]
func (ac *ActivityController) DeleteActivity(c echo.Context) error {
	db := ac.db
	response := entities.Response[entities.NullStruct]{}
	id, _ := strconv.Atoi(c.Param("id"))

	activity := activities.Activity{}
	condition := activities.Activity{ActivityID: id}
	if err := db.Where(&condition).Find(&activity).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if activity == (activities.Activity{}) {
		response.Status = helpers.NOT_FOUND
		response.Message = "Activity with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	if err := db.Delete(&activity).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.NullStruct{}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS

	return c.JSON(http.StatusOK, response)
}
