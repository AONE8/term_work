package main

import (
	//"time"

	"database/sql"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/robfig/cron/v3"

	// jwtware "github.com/gofiber/contrib/jwt"
	"encoding/json"
	"fmt"
	"golang-term-work/db"
	"golang-term-work/models"
	"os"
	//"github.com/golang-jwt/jwt/v5"
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
	go RunServ()
	ch := make(chan string, 1)
	go cronStart(ch)
	for {
		msg := <-ch
		fmt.Println(msg)
	}

}

func RunServ() {
	app := fiber.New()

	app.Get("/", IndexPage)
	app.Post("/api/forecast/now", weatherNowHandler)
	app.Post("/api/forecast/history", weatherHistHandler)

	app.Use(recover.New(), logger.New())

	err := app.Listen(models.AuthConfig.WeatherPort)
	if err != nil {
		fmt.Println(err)
	}
}

func cronStart(ch chan string) {
	c := cron.New()
	for _, s := range models.AuthConfig.Schedule {
		t, err := time.Parse("15:04", s)
		if err != nil {
			fmt.Println(err)
			return
		}

		cronExp := fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour())
		c.AddFunc(cronExp, func() { ch <- getWeatherForecast() })
	}

	c.Start()
}

func IndexPage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}

func getWeatherForecast() string {
	for _, forecast := range models.AuthConfig.ForecastList {
		weatherNowReq := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/complete.json?lat=%s&lon=%s",
			forecast.Lat,
			forecast.Lon)

		req, err := http.NewRequest("GET", weatherNowReq, nil)

		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}

		req.Header.Set("User-Agent", "MyTestApp/0.1")
		req.Header.Set("Accept", "application/json")

		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := netClient.Do(req)

		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(err.Error())
			return err.Error()
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}

		var yrResp models.YrResp

		err = json.Unmarshal(body, &yrResp)
		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}

		setWeatherForecastToDB(yrResp, forecast.Name, forecast.Lat, forecast.Lon)
	}

	return "ok"
}

func setWeatherForecastToDB(yrResp models.YrResp, cityName string, lat string, lon string) error {
	connection, err := db.GetConnection(models.AuthConfig.Db.User,
		models.AuthConfig.Db.Password,
		models.AuthConfig.Db.Port,
		models.AuthConfig.Db.Host,
		models.AuthConfig.Db.Name)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer connection.Close()

	query := "INSERT INTO  `WEATHER` (`city`, `latitude`, `longetude`, `temperature`, `wind_speed`, `time`) VALUES (?, ?, ?, ?, ?, ?)"

	res, err := connection.Exec(query, cityName, lat, lon,
		yrResp.Properties.Timeseries[0].Data.Instant.Details.AirTemperature,
		yrResp.Properties.Timeseries[0].Data.Instant.Details.WindSpeed,
		yrResp.Properties.Timeseries[0].Time)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func getWeatherForecastToDB(req models.WeatherHistReq) (dataArr []models.WeatherHist, err error) {
	var (
		rows       *sql.Rows
		connection *sql.DB
	)

	connection, err = db.GetConnection(models.AuthConfig.Db.User,
		models.AuthConfig.Db.Password,
		models.AuthConfig.Db.Port,
		models.AuthConfig.Db.Host,
		models.AuthConfig.Db.Name)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer connection.Close()

	query := fmt.Sprintf("SELECT * FROM `weather` WHERE city = \"%s\" AND time >= \"%s\" AND time <= \"%s\"", req.City, req.From, req.To)
	rows, err = connection.Query(query)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var data models.WeatherHist
		if err = rows.Scan(&data.Id, &data.City, &data.Latitude, &data.Longetude, &data.Temperature, &data.WindSpeed, &data.Time); err != nil {
			fmt.Println("scan", err.Error())
			return
		}

		dataArr = append(dataArr, data)

	}

	return
}

func weatherNowHandler(c *fiber.Ctx) error {

	var (
		weatherReq models.WeatherNowByLoc
		yrResp     models.YrResp
	)
	err := c.BodyParser(&weatherReq)
	if err != nil {
		fmt.Println("Error parsing JSON:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	weatherNowReq := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/complete.json?lat=%s&lon=%s",
		weatherReq.Lat,
		weatherReq.Lon)

	req, err := http.NewRequest("GET", weatherNowReq, nil)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"success": false, "err": err})
	}

	req.Header.Set("User-Agent", "MyTestApp/0.1")
	req.Header.Set("Accept", "application/json")

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Do(req)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"success": false, "err": err})
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.SendStatus(resp.StatusCode)
		return c.JSON(fiber.Map{"success": false, "err": err})
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"success": false, "err": err})
	}

	err = json.Unmarshal(body, &yrResp)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"success": false, "err": err})
	}

	return c.JSON(fiber.Map{"success": true, "req": yrResp})
}

func weatherHistHandler(c *fiber.Ctx) error {

	var request models.WeatherHistReq

	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println("Error parsing JSON:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	data, err := getWeatherForecastToDB(request)

	if err != nil {
		fmt.Println("Error parsing JSON:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}
