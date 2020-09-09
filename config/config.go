package config

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"weatherWeb/models"
)

type Env struct {
	Db        *sql.DB
	Templates *template.Template
	DebugMode bool
}
type LogConf struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}
type dbConf struct {
	DbPath   string `json:"dbPath"`
	User     string `json:"user"`
	Password string `json:"password"`
	SslMode  string `json:"sslMode"`
}
type Config struct {
	DebugMode    bool   `json:"debugMode"`
	Db           dbConf `json:"db"`
	TemplatePath string `json:"templatePath"`
	LogFile      string `json:"logFile"`
}

//Seleccion del entorno de trabajo
func GetEnv(Option string) (Env, LogConf) {
	env := Env{}
	logger := LogConf{}
	var config Config
	switch Option {
	case "debug":
		config = loadConfigFile("config/debugConfig.json")
	case "prod":
		config = loadConfigFile("config/config.json")
	default:
		log.Fatal("Environment not found")
	}
	db, err := models.NewDB("postgres://" + config.Db.User + ":" + config.Db.Password + "@" + config.Db.DbPath + "?sslmode=" + config.Db.SslMode)
	if err != nil {
		log.Panic(err)
	}
	env.Db = db
	logger.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	logger.Warn = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	env.Templates = template.Must(template.ParseGlob(config.TemplatePath))
	env.DebugMode = config.DebugMode
	return env, logger
}

func FatalError(debugMode bool, logger LogConf, err error) {
	logger.Error.Println(err)
	if !debugMode {
		os.Exit(1)
	}
}
func loadConfigFile(fileName string) Config {
	var config Config
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
