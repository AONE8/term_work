package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-term-work/models"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/golang-jwt/jwt/v5"
)

func init() {
	var (
		byteValues []byte
		err        error
	)

	if byteValues, err = os.ReadFile("config.json"); err != nil {
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
	app.Use(recover.New(), logger.New())

	app.Post("/api/auth", postAuthHandler)
	app.Get("/api/auth", getAuthHandler)

	app.Post("/api/forecast/now", postForecastNow)

	err := app.Listen(models.AuthConfig.GatewayPort)
	if err != nil {
		fmt.Println(err)
	}
}

func postAuthHandler(c *fiber.Ctx) (err error) {
	var data models.LoginReq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal JSON",
		})
	}

	resp, err := client.Post("http://localhost:3001/api/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	return c.Status(resp.StatusCode).Send(responseBody)
}

func getAuthHandler(c *fiber.Ctx) (err error) {

	resp, err := http.Get("http://localhost:3001/api/auth")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	return c.Status(resp.StatusCode).Send(responseBody)
}

func postForecastNow(c *fiber.Ctx) error {
	var data models.WeatherNowByLoc
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal JSON",
		})
	}

	resp, err := client.Post("http://localhost:3001/api/forecast/now", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	return c.JSON(fiber.Map{"success": true, "data": responseBody})
}
