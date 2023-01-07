//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

func TestITCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	var e Expense

	res := request(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), body)
	err := json.NewDecoder(res.Body).Decode(&e)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "buy a new phone", e.Title)
	assert.Equal(t, 39000.0, e.Amount)
	assert.Equal(t, "buy a new phone", e.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, e.Tags)
}

func TestITGetExpenseByID(t *testing.T) {
	c := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, c.ID), nil)
	err := json.NewDecoder(res.Body).Decode(&latest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.Equal(t, c.Title, latest.Title)
	assert.Equal(t, c.Amount, latest.Amount)
	assert.Equal(t, c.Note, latest.Note)
	assert.Equal(t, c.Tags, latest.Tags)
}

func TestITUpdateExpense(t *testing.T) {
	c := seedExpense(t)
	body := bytes.NewBufferString(`{
		"title": "iPhone 14 Pro Max 1TB",
		"amount": 66900,
		"note": "birthday gift from my love", 
		"tags": ["gadget"]
	}`)

	var latest Expense
	res := request(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, c.ID), body)
	err := json.NewDecoder(res.Body).Decode(&latest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.Equal(t, "iPhone 14 Pro Max 1TB", latest.Title)
	assert.Equal(t, 66900.0, latest.Amount)
	assert.Equal(t, "birthday gift from my love", latest.Note)
	assert.Equal(t, []string{"gadget"}, latest.Tags)
}

func TestITGetAllExpense(t *testing.T) {
	var items []Expense
	res := request(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), nil)
	err := json.NewDecoder(res.Body).Decode(&items)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	for _, item := range items {
		assert.NotEqual(t, 0, item.ID)
	}
}

func seedExpense(t *testing.T) Expense {
	var e Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	res := request(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), body)
	err := json.NewDecoder(res.Body).Decode(&e)
	if err != nil {
		t.Fatal("can't create:", err)
	}
	return e
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}
