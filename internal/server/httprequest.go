package server

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/database"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return e
}

func Get(e *echo.Echo, db *sql.DB, path string) {
	e.GET(path, func(c echo.Context) error {
		items, err := database.GetAllRSSItems(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch items"})
		}
		return c.JSON(http.StatusOK, items)
	})
}
