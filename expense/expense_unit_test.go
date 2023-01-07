//go:build unit

package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setup() {
	InitDB()
}

func shutdown() {
	CloseDB()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func createExpense(body *bytes.Buffer, item interface{}) error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	CreateExpenseHandler(c)

	return json.NewDecoder(rec.Body).Decode(item)
}

func TestCreateExpense(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
    	"amount": 89,
    	"note": "no discount",
    	"tags": ["beverage"]
	}`)

	c, rec := makeRequest(http.MethodPost, "/expenses", body)

	// Act
	err := CreateExpenseHandler(c)
	var item Expense
	json.NewDecoder(rec.Body).Decode(&item)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, 0, item.ID)
	assert.Equal(t, "apple smoothie", item.Title)
	assert.Equal(t, 89.0, item.Amount)
	assert.Equal(t, "no discount", item.Note)
	assert.Equal(t, []string{"beverage"}, item.Tags)
}

func TestGetExpenseByID(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
		"title": "iPhone 14 Pro Max 1TB",
		"amount": 66900,
		"note": "birthday gift from my love", 
		"tags": ["gadget"]
	}`)
	var item Expense
	createExpense(body, &item)

	c, rec := makeRequest(http.MethodGet, "/expenses", nil)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", item.ID))

	// Act
	err := GetExpenseHandler(c)

	// Assertions
	var getItem Expense
	json.NewDecoder(rec.Body).Decode(&getItem)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, getItem.ID, item.ID)
	assert.Equal(t, getItem.Title, item.Title)
	assert.Equal(t, getItem.Amount, item.Amount)
	assert.Equal(t, getItem.Note, item.Note)
	assert.Equal(t, getItem.Tags, item.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
    	"amount": 89,
    	"note": "no discount",
    	"tags": ["beverage"]
	}`)
	var item Expense
	createExpense(body, &item)

	updateBody := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 baht", 
		"tags": ["food", "beverage"]
	}`)
	c, rec := makeRequest(http.MethodPut, "/expenses", updateBody)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", item.ID))

	// Act
	err := UpdateExpenseHandler(c)

	// Assertions
	var updateItem Expense
	json.NewDecoder(rec.Body).Decode(&updateItem)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, updateItem.ID, item.ID)
	assert.Equal(t, "strawberry smoothie", updateItem.Title)
	assert.Equal(t, 79.0, updateItem.Amount)
	assert.Equal(t, "night market promotion discount 10 baht", updateItem.Note)
	assert.Equal(t, []string{"food", "beverage"}, updateItem.Tags)
}

func TestGetAllExpense(t *testing.T) {
	// Arrange
	c, rec := makeRequest(http.MethodGet, "/expenses", nil)

	// Act
	err := GetExpensesHandler(c)

	// Assertions
	var items []Expense
	json.NewDecoder(rec.Body).Decode(&items)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	for _, item := range items {
		assert.NotEqual(t, 0, item.ID)
	}
}

func makeRequest(method, url string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
