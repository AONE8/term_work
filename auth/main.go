package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	// jwtware "github.com/gofiber/contrib/jwt"
	"encoding/json"
	"fmt"
	"golang-term-work/db"
	"golang-term-work/middlewares"
	"golang-term-work/models"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func init() {
	var (
		byteValues []byte
		err        error
	)

	if byteValues, err = os.ReadFile("../config.json"); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = json.Unmarshal(byteValues, &models.AuthConfig); err != nil {
		fmt.Println(err.Error())
		return
	}

}

func main() {
	app := fiber.New()

	app.Post("/api/auth", authHandler)
	app.Get("/api/auth", middlewares.AuthMiddleware(models.AuthConfig.Secret), checkAuthHandler)
	app.Use(recover.New(), logger.New())

	err := app.Listen(models.AuthConfig.AuthPort)
	if err != nil {
		fmt.Println(err)
	}
}

func authHandler(c *fiber.Ctx) error {

	var loginRequest models.LoginReq
	err := c.BodyParser(&loginRequest)
	if err != nil {
		fmt.Println("Error parsing JSON:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	connection, err := db.GetConnection(models.AuthConfig.Db.User,
		models.AuthConfig.Db.Password,
		models.AuthConfig.Db.Port,
		models.AuthConfig.Db.Host,
		models.AuthConfig.Db.Name)

	if err != nil {
		fmt.Println("DB", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer connection.Close()

	query := fmt.Sprintf("SELECT * FROM USERS WHERE NAME='%s'", loginRequest.UserName)

	rows, err := connection.Query(query)

	if err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.Password); err != nil {
			fmt.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  loginRequest.UserName,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(models.AuthConfig.Secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func checkAuthHandler(c *fiber.Ctx) error {

	// user, ok := c.Locals("user").(*jwt.Token)
	// if !ok {
	// 	return c.SendStatus(fiber.StatusUnauthorized)
	// }

	// fmt.Println(user)

	// claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)

	// return c.JSON(fiber.Map{"message": "Authenticated", "user": name})
	return c.SendStatus(fiber.StatusAccepted)
}
