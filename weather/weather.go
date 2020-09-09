package weather

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weatherWeb/config"
)

type Info struct {
	Description string
	Temp        float64
	MaxTemp     float64
	MinTemp     float64
	FeelTemp    float64
	Pressure    float64
	Humidity    float64
}

//Historial del tiempo
type Record struct {
	Time        time.Time
	WeatherInfo string
}

var Records map[string]Record
var logger config.LogConf
var cities []string // array de ciudades disponibles

func InitWeather(log config.LogConf, cityList []string) {
	logger = log
	Records = make(map[string]Record)
	cities = cityList
}

//Realiza la peticiòn del tiempo a la API externa
func GetWeather(city string) (string, error) {
	url := "https://community-open-weather-map.p.rapidapi.com/weather?lang=es&units=metric&q=" + city
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("x-rapidapi-host", "community-open-weather-map.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "a955abc71fmsh13d85fad086053dp179b4fjsn1f05b8f85c97")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil || res.StatusCode > 200 {
		if res.StatusCode > 200 {
			return "", errors.New(res.Status + " " + city)
		}
		return "", err
	}
	return string(body), err
}

//Parse de los datos del tiempo extraídos de la API externa
func JsonParser(jsonData string) Info {
	var weatherInfo Info
	weatherInfo.Description = gjson.Get(jsonData, "weather").Array()[0].Map()["description"].String()
	weatherInfo.Humidity = gjson.Get(jsonData, "main.humidity").Float()
	weatherInfo.Pressure = gjson.Get(jsonData, "main.pressure").Float()
	weatherInfo.Temp = gjson.Get(jsonData, "main.temp").Float()
	weatherInfo.FeelTemp = gjson.Get(jsonData, "main.feels_like").Float()
	weatherInfo.MaxTemp = gjson.Get(jsonData, "main.temp_max").Float()
	weatherInfo.MinTemp = gjson.Get(jsonData, "main.temp_min").Float()
	return weatherInfo
}

//Comprobamos el historico
func FindHistoricalRecord(city string) (bool, string) {
	//Comprobamos si tenemos historico de la ciudad
	if _, found := Records[city]; found {
		//Comprobamos la antigüedad del registro historico
		if Records[city].Time.Add(11 * time.Minute).Before(time.Now()) {
			return false, ""
		} else {
			return true, Records[city].WeatherInfo
		}
	} else { //Si no hay historico
		return false, ""
	}
}

//Tarea programada de recarga del historio del tiempo
func CronWeatherLoader() {
	log.Println("Reloading weather records")
	for _, city := range cities {
		//solicitamos la info al exterior
		cityData, err := GetWeather(city)
		if err != nil {
			logger.Error.Println(err)
		} else {
			logger.Info.Println(city + " weather data reloaded")
		}
		//Guardamos la info en el historico
		var weatherRecord Record
		weatherRecord.Time = time.Now()
		weatherRecord.WeatherInfo = cityData
		Records[city] = weatherRecord
	}
	logger.Info.Println("Weather records reloaded")
}
