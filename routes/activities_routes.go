package routes

import (
	"github.com/khasyah-fr/sky-todo/config"
	controllers "github.com/khasyah-fr/sky-todo/controllers/activities"

	"github.com/labstack/echo/v4"
)

func ActivityRoute(e *echo.Echo) {
	activityController := controllers.NewActivityController(database.DB.GetConnection())

	e.GET("/activity-groups", activityController.GetAllActivity)
	e.GET("/activity-groups/:id", activityController.GetOneActivity)
	e.POST("/activity-groups", activityController.CreateActivity)
	e.PATCH("/activity-groups/:id", activityController.UpdateActivity)
	e.DELETE("/activity-groups/:id", activityController.DeleteActivity)
}
