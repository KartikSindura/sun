package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

type Weather struct {
	Location struct {
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				Chanceofrain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	q := "India"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	key := "a3ba414740374a5ea17101926240307"
	// key := goDotEnvVariable("API_KEY")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q="+q+"&days=1&aqi=no&alerts=no", key)

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather api not available.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	// fmt.Print(weather)
	fmt.Printf("%s, %s: %.0fC, %s\n", weather.Location.Region, weather.Location.Country, weather.Current.TempC, weather.Current.Condition.Text)

	for _, hour := range weather.Forecast.Forecastday[0].Hour {

		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		msg := fmt.Sprintf("%s - %.0fC, %.0f, %s\n", date.Format("15:05"), hour.TempC, hour.Chanceofrain, hour.Condition.Text)
		if hour.Chanceofrain < 40 {
			color.White(msg)
		} else {
			color.Blue(msg)
		}

	}
}
