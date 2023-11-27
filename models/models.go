package models

import (
	"time"
)

var AuthConfig Config

type LoginReq struct {
	UserName string `json:"name"`
	Pass     string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type User struct {
	Id       string
	Name     string
	Password string
}

type Config struct {
	Secret       string      `json:"secret"`
	AuthPort     string      `json:"authPort"`
	WeatherPort  string      `json:"weatherPort"`
	GatewayPort  string      `json:"gatewayPort"`
	Db           DB          `json:"db"`
	ForecastList []Forecasts `json:"forecasts"`
	Schedule     []string    `json:"schedule"`
}

type Forecasts struct {
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

type DB struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
}

type WeatherNowByLoc struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type WeatherHist struct {
	Id          int       `json:"id"`
	City        string    `json:"city"`
	Latitude    float64   `json:"lat"`
	Longetude   float64   `json:"lon"`
	Temperature float64   `json:"temperature"`
	WindSpeed   float64   `json:"wind_speed"`
	Time        time.Time `json:"time"`
}

type WeatherHistReq struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"lat"`
	Longetude float64 `json:"lon"`
	From      string  `json:"from"`
	To        string  `json:"to"`
}

type YrResp struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at"`
			Units     struct {
				AirPressureAtSeaLevel    string `json:"air_pressure_at_sea_level"`
				AirTemperature           string `json:"air_temperature"`
				AirTemperatureMax        string `json:"air_temperature_max"`
				AirTemperatureMin        string `json:"air_temperature_min"`
				CloudAreaFraction        string `json:"cloud_area_fraction"`
				CloudAreaFractionHigh    string `json:"cloud_area_fraction_high"`
				CloudAreaFractionLow     string `json:"cloud_area_fraction_low"`
				CloudAreaFractionMedium  string `json:"cloud_area_fraction_medium"`
				DewPointTemperature      string `json:"dew_point_temperature"`
				FogAreaFraction          string `json:"fog_area_fraction"`
				PrecipitationAmount      string `json:"precipitation_amount"`
				RelativeHumidity         string `json:"relative_humidity"`
				UltravioletIndexClearSky string `json:"ultraviolet_index_clear_sky"`
				WindFromDirection        string `json:"wind_from_direction"`
				WindSpeed                string `json:"wind_speed"`
			} `json:"units"`
		} `json:"meta"`
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel    float64 `json:"air_pressure_at_sea_level"`
						AirTemperature           float64 `json:"air_temperature"`
						CloudAreaFraction        float64 `json:"cloud_area_fraction"`
						CloudAreaFractionHigh    float64 `json:"cloud_area_fraction_high"`
						CloudAreaFractionLow     float64 `json:"cloud_area_fraction_low"`
						CloudAreaFractionMedium  float64 `json:"cloud_area_fraction_medium"`
						DewPointTemperature      float64 `json:"dew_point_temperature"`
						FogAreaFraction          float64 `json:"fog_area_fraction"`
						RelativeHumidity         float64 `json:"relative_humidity"`
						UltravioletIndexClearSky float64 `json:"ultraviolet_index_clear_sky"`
						WindFromDirection        float64 `json:"wind_from_direction"`
						WindSpeed                float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
					} `json:"details"`
				} `json:"next_12_hours"`
				Next1Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_1_hours"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						AirTemperatureMax   float64 `json:"air_temperature_max"`
						AirTemperatureMin   float64 `json:"air_temperature_min"`
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_6_hours"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}
