package su3

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type WeatherData struct {
	Timepoint    int `json:"timepoint"`
	Cloudcover   int `json:"cloudcover"`
	Seeing       int `json:"seeing"`
	Transparency int `json:"transparency"`
	Temp2m       int `json:"temp2m"`
}

type WeatherClient struct {
	BaseURL string
	Client  *http.Client
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		BaseURL: "http://www.7timer.info/bin/api.pl",
		Client:  InitClient(),
	}
}

func (wc *WeatherClient) FetchWeather(lon, lat float64, unit, product string) (*WeatherData, error) {
	params := url.Values{}
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("unit", unit)
	params.Add("product", product)
	params.Add("output", "json")

	reqUrl := fmt.Sprintf("%s?%s", wc.BaseURL, params.Encode())

	log.Printf("Fetching weather data from %s", reqUrl)

	resp, err := wc.Client.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	var data struct {
		Weather []WeatherData `json:"dataseries"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Weather) > 0 {
		return &data.Weather[0], nil
	}

	return nil, fmt.Errorf("no weather data found")
}

// go run main.go -lon=3.1415 -lat=42.1234 -unit=metric -product=civil
func GetWeather() {
	var (
		lon, lat      float64
		unit, product string
	)

	flag.Float64Var(&lon, "lon", 2.3522, "Longitude (e.g., 2.3522 for Paris)")
	flag.Float64Var(&lat, "lat", 48.8566, "Latitude (e.g., 48.8566 for Paris)")
	flag.StringVar(&unit, "unit", "metric", "Unit of measurement (e.g., metric or imperial)")
	flag.StringVar(&product, "product", "civil", "Product type (e.g., civil)")

	flag.Parse()

	client := NewWeatherClient()

	weather, err := client.FetchWeather(lon, lat, unit, product)
	if err != nil {
		panic(err)
	}

	jsonData, _ := json.Marshal(weather)
	fmt.Println(string(jsonData))
}
