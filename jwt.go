package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

type jwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

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
		// TODO: create jwt token
		token, err := createJwtToken()
		if err != nil {
			log.Println("error creating jwt token", err)
			return c.String(http.StatusInternalServerError, "wrong")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You are authorized!",
			"token":   token,
		})
	}
	return c.String(http.StatusUnauthorized, "You entered the wrong password!")
}

func createJwtToken() (string, error) {
	claims := jwtClaims{
		"pavelgolang",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("Secret"))
	if err != nil {
		return "", err
	}
	return token, nil
}
func mainJwt(c echo.Context) error {
	return c.String(http.StatusOK, "jwt page.")
}

func main() {
	e := echo.New()
	g := e.Group("/admin")
	cookieGroup := e.Group("/cookie")
	jwtGroup := e.Group("/jwt")

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("SECRET"),
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "pavelgolang" && password == "qwerty" {
			return true, nil
		}
		return false, nil
	}))
	cookieGroup.GET("/main", mainCookie)
	jwtGroup.GET("/main", mainJwt)

	e.GET("/login", login)
	g.GET("/main", mainAdmin)
	e.POST("/patients", createPatient)
	e.GET("/patients/:id", getPatient)
	e.PUT("/patients/:id", updatePatient)
	e.DELETE("/patients/:id", deletePatient)

	e.Logger.Fatal(e.Start(":1323"))
}
