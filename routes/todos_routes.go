package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/khasyah-fr/sky-todo/config"
	"github.com/khasyah-fr/sky-todo/controllers/todos"
)

func ConfigureTodoRoutes(e *echo.Echo) {
	todoController := todos.NewTodoController(config.DB.GetConnection())

	group := e.Group("/todo-items")
	group.GET("", todoController.GetAllTodos)
	group.GET("/:id", todoController.GetTodo)
	group.POST("", todoController.CreateTodo)
	group.PATCH("/:id", todoController.UpdateTodo)
	group.DELETE("/:id", todoController.DeleteTodo)
}
