package controllers

import (
	"net/http"
	"strconv"

	"github.com/khasyah-fr/sky-todo/entities"
	"github.com/khasyah-fr/sky-todo/entities/todos"
	"github.com/khasyah-fr/sky-todo/helpers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TodoController struct {
	db *gorm.DB
}

func NewTodoController(db *gorm.DB) *TodoController {
	return &TodoController{db: db}
}

// @Summary Get all todos
// @Description Retrieves all todos
// @Produce json
// @Success 200 {array} entities.GetTodoResponse
// @Router /todos [get]
func (tc *TodoController) GetAllTodos(c echo.Context) error {
	db := tc.db
	response := entities.Response[interface{}]{}

	allTodos := []todos.Todo{}
	if activityGroupID, err := strconv.Atoi(c.QueryParam("activity_group_id")); err == nil {
		condition := todos.Todo{ActivityGroupID: activityGroupID}
		if err := db.Where(&condition).Find(&allTodos).Error; err != nil {
			response.Status = helpers.FAILED
			response.Message = helpers.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		if err := db.Find(&allTodos).Error; err != nil {
			response.Status = helpers.FAILED
			response.Message = helpers.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	if len(allTodos) == 0 {
		response.Status = helpers.SUCCESS
		response.Message = helpers.SUCCESS
		response.Data = make([]string, 0)
		return c.JSON(http.StatusOK, response)
	}

	data := make([]entities.GetTodoResponse, 0, len(allTodos))
	for _, todo := range allTodos {
		data = append(data, entities.GetTodoResponse{
			ID:              todo.TodoID,
			ActivityGroupID: todo.ActivityGroupID,
			Title:           todo.Title,
			IsActive:        todo.IsActive,
			Priority:        todo.Priority,
			CreatedAt:       todo.CreatedAt,
			UpdatedAt:       todo.UpdatedAt,
		})
	}

	response.Data = data
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS
	return c.JSON(http.StatusOK, response)
}

// @Summary Get a todo
// @Description Retrieves a specific todo by ID
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} entities.GetTodoResponse
// @Router /todos/{id} [get]
func (tc *TodoController) GetTodo(c echo.Context) error {
	db := tc.db
	response := entities.Response[entities.GetTodoResponse]{}
	id, _ := strconv.Atoi(c.Param("id"))

	todo := todos.Todo{}
	condition := todos.Todo{TodoID: id}
	if err := db.Where(&condition).Find(&todo).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if todo == (todos.Todo{}) {
		response.Status = helpers.NOT_FOUND
		response.Message = "Todo with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	response.Data = entities.GetTodoResponse{
		ID:              todo.TodoID,
		ActivityGroupID: todo.ActivityGroupID,
		Title:           todo.Title,
		IsActive:        todo.IsActive,
		Priority:        todo.Priority,
		CreatedAt:       todo.CreatedAt,
		UpdatedAt:       todo.UpdatedAt,
	}

	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS
	return c.JSON(http.StatusOK, response)
}

// @Summary Create a todo
// @Description Creates a new todo
// @Accept json
// @Produce json
// @Param todo body entities.CreateTodoRequest true "Todo object"
// @Success 201 {object} entities.GetTodoResponse
// @Router /todos [post]
func (tc *TodoController) CreateTodo(c echo.Context) error {
	db := tc.db
	response := entities.Response[entities.GetTodoResponse]{}
	request := entities.CreateTodoRequest{}

	if err := c.Bind(&request); err != nil {
		response := entities.Response[[]string]{}
		response.Data = make([]string, 0)
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

	if request.ActivityGroupID == 0 {
		response := entities.Response[[]string]{}
		response.Data = make([]string, 0)
		response.Status = helpers.ERROR_BAD_REQUEST
		response.Message = "activity_group_id cannot be null"
		return c.JSON(http.StatusBadRequest, response)
	}

	todo := todos.Todo{
		Title:           request.Title,
		ActivityGroupID: request.ActivityGroupID,
		IsActive:        request.IsActive,
		Priority:        request.Priority,
	}

	if !request.IsActive {
		todo.IsActive = true
	}

	if request.Priority == "" {
		todo.Priority = "very-high"
	}

	if err := db.Create(&todo).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.GetTodoResponse{
		ID:              todo.TodoID,
		ActivityGroupID: todo.ActivityGroupID,
		Title:           todo.Title,
		IsActive:        todo.IsActive,
		Priority:        todo.Priority,
		CreatedAt:       todo.CreatedAt,
		UpdatedAt:       todo.UpdatedAt,
	}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS
	return c.JSON(http.StatusCreated, response)
}

// @Summary Update a todo
// @Description Updates an existing todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body entities.UpdateTodoRequest true "Updated todo object"
// @Success 200 {object} entities.GetTodoResponse
// @Router /todos/{id} [put]
func (tc *TodoController) UpdateTodo(c echo.Context) error {
	db := tc.db
	response := entities.Response[entities.GetTodoResponse]{}
	id, _ := strconv.Atoi(c.Param("id"))
	request := entities.UpdateTodoRequest{}

	if err := c.Bind(&request); err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	todo := todos.Todo{}
	condition := todos.Todo{TodoID: id}
	if err := db.Where(&condition).Find(&todo).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if todo == (todos.Todo{}) {
		response := entities.Response[[]string]{}
		response.Data = make([]string, 0)
		response.Status = helpers.NOT_FOUND
		response.Message = "Todo with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	if err := db.Model(&todo).Updates(todos.Todo{
		Title:    request.Title,
		Priority: request.Priority,
		IsActive: request.IsActive,
	}).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.GetTodoResponse{
		ID:              todo.TodoID,
		ActivityGroupID: todo.ActivityGroupID,
		Title:           todo.Title,
		IsActive:        todo.IsActive,
		Priority:        todo.Priority,
		CreatedAt:       todo.CreatedAt,
		UpdatedAt:       todo.UpdatedAt,
	}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS
	return c.JSON(http.StatusOK, response)
}

// @Summary Delete a todo
// @Description Deletes a todo by ID
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200
// @Router /todos/{id} [delete]
func (tc *TodoController) DeleteTodo(c echo.Context) error {
	db := tc.db
	response := entities.Response[entities.NullStruct]{}
	id, _ := strconv.Atoi(c.Param("id"))

	todo := todos.Todo{}
	condition := todos.Todo{TodoID: id}
	if err := db.Where(&condition).Find(&todo).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if todo == (todos.Todo{}) {
		response.Status = helpers.NOT_FOUND
		response.Message = "Todo with ID " + c.Param("id") + " Not Found"
		return c.JSON(http.StatusNotFound, response)
	}

	if err := db.Delete(&todo).Error; err != nil {
		response.Status = helpers.FAILED
		response.Message = helpers.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = entities.NullStruct{}
	response.Status = helpers.SUCCESS
	response.Message = helpers.SUCCESS
	return c.JSON(http.StatusOK, response)
}
