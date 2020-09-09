package main

import (
	"bufio"
	"gopkg.in/robfig/cron.v2"
	"net/http"
	"os"
	"time"
	"weatherWeb/config"
	"weatherWeb/models"
	"weatherWeb/weather"
)

var cities []string // array de ciudades disponibles
//Información del tiempo de una ciudad

var env config.Env
var logger config.LogConf

func main() {
	//Selecciòn del entorno
	if len(os.Args) > 1 {
		env, logger = config.GetEnv(os.Args[1])
	} else if len(os.Getenv("ENV")) > 0 {
		env, logger = config.GetEnv(os.Getenv("ENV"))
	} else {
		env, logger = config.GetEnv("debug")
	}
	logger.Info.Println("Starting weather server...")
	//Manejadores de urls
	http.HandleFunc("/", index)
	http.HandleFunc("/showWeather", showWeather)
	http.Handle("/getapikey", getApiKey(env))
	//Start cron schedule
	crn := cron.New()
	//cargamos el listado de ciudades
	loadCitiesList()
	weather.InitWeather(logger, cities)
	_, err := crn.AddFunc("@every 0h10m0s", weather.CronWeatherLoader)
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}
	//crn.Start()
	//Arrancamos el servidor
	logger.Info.Println("Starting server...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	} else {
		logger.Info.Println("Server started")
	}
}

//Manejador de la url base
func index(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title  string
		Cities []string
	}{
		Title:  "Iweather",
		Cities: cities,
	}
	err := env.Templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}
}

//Manejador de la pagina del tiempo
func showWeather(w http.ResponseWriter, r *http.Request) {
	//Extraemos el resultado del form
	var err = r.ParseForm()
	var cityData string
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}
	//Comprobamos si tenemos historico de la ciudad
	if valid, data := weather.FindHistoricalRecord(r.Form["city"][0]); valid {
		cityData = data
	} else { //Si no hay historico valido
		//solicitamos la info al exterior
		data, err := weather.GetWeather(r.Form["city"][0])
		cityData = data
		if err != nil {
			logger.Error.Println(err)
			index(w, r)
			return
		}
		//Guardamos la info en el historico
		var weatherRecord weather.Record
		weatherRecord.Time = time.Now()
		weatherRecord.WeatherInfo = cityData
		weather.Records[r.Form["city"][0]] = weatherRecord
	}
	//Parseamos la info en json a nuestro formato
	weatherInfo := weather.JsonParser(cityData)

	//creamos la estructura de datos a pasar a la plantilla
	data := struct {
		Title       string
		WeatherInfo weather.Info
		City        string
	}{
		Title:       "El tiempo en " + r.Form["city"][0],
		City:        r.Form["city"][0],
		WeatherInfo: weatherInfo,
	}
	//Devolvemos la plantilla
	err = env.Templates.ExecuteTemplate(w, "showWeather", data)
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}
}

//Pruba de BBDD
func getApiKey(env config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			logger.Warn.Println("Error 405 on getApiKey")
			return
		}
		_, err := models.GetGetUserById(env.Db, 1)
		if err != nil {
			config.FatalError(env.DebugMode, logger, err)
		}
	})
}

//Carga de las ciudades disponibles en el txt
func loadCitiesList() {
	file, err := os.Open("cities.txt")

	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		cities = append(cities, scanner.Text())
	}
	err = file.Close()
	if err != nil {
		config.FatalError(env.DebugMode, logger, err)
	}
}
