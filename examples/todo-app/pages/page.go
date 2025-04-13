package index

import (
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Page struct{}

type todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []todo = []todo{
	{
		ID:        1,
		Title:     "Open a GitHub account",
		Completed: true,
	},
	{
		ID:        2,
		Title:     "Build a FullStack App with GV",
		Completed: false,
	},
}

func (Page) List(c echo.Context) error {
	return c.JSON(200, todos)
}

func (Page) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, t := range todos {
		if t.ID == id {
			return c.JSON(200, t)
		}
	}

	return c.JSON(404, "Not found")
}

func (Page) Create(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = len(todos) + 1
	todos = append(todos, t)

	return c.JSON(201, t)
}

func (Page) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	for i, _todo := range todos {
		if _todo.ID == id {
			todos[i] = todo{
				ID:        id,
				Title:     todos[i].Title,
				Completed: t.Completed,
			}
			return c.JSON(200, t)
		}
	}

	return c.JSON(404, "Not found")
}

func (Page) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			return c.JSON(200, "OK")
		}
	}

	return c.JSON(404, "Not found")
}
