package main

import (
	"net/http"
	"strconv"
	"time"

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
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "admin main page")
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "cookie main page")
}
func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "pavelgolang" && password == "qwerty" {
		cookie := &http.Cookie{}

		cookie.Name = "sessionID"
		cookie.Value = "same_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(cookie)
		return c.String(http.StatusOK, "You are authorized!")
	}
	return c.String(http.StatusUnauthorized, "You entered the wrong password!")
}

func main() {
	e := echo.New()
	g := e.Group("/admin")
	cookieGroup := e.Group("/cookie")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "pavelgolang" && password == "qwerty" {
			return true, nil
		}
		return false, nil
	}))
	cookieGroup.GET("/main", mainCookie)

	e.GET("/login", login)
	g.GET("/main", mainAdmin)
	e.POST("/patients", createPatient)
	e.GET("/patients/:id", getPatient)
	e.PUT("/patients/:id", updatePatient)
	e.DELETE("/patients/:id", deletePatient)

	e.Logger.Fatal(e.Start(":1323"))
}
