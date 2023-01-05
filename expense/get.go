package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	e := Expense{}
	row := stmt.QueryRow(id).Scan(&e.ID, &e.Title, &e.Amount, &e.Note, &e.Tags)
	switch row {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
}
