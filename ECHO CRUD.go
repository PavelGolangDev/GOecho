package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	patient struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Diag  string `json:"diag"`
		Vrach string `json:"vrach"`
	}
)

var (
	patients = map[int]*patient{}
	seq      = 1
)

func createPatient(c echo.Context) error {
	u := &patient{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	patients[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getPatient(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, patients[id])
}

func updatePatient(c echo.Context) error {
	u := new(patient)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	patients[id].Name = u.Name
	return c.JSON(http.StatusOK, patients[id])
}

func deletePatient(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(patients, id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/patients", createPatient)
	e.GET("/patients/:id", getPatient)
	e.PUT("/patients/:id", updatePatient)
	e.DELETE("/patients/:id", deletePatient)

	e.Logger.Fatal(e.Start(":1323"))
}
