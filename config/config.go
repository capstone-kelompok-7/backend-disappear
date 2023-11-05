package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DBPort     int
	DBHost     string
	DBUser     string
	DBPass     string
	DBName     string
	Secret     string
}

func InitConfig() *Config {
	var res = new(Config)
	res = loadConfig()

	if res == nil {
		log.Fatal("Config : cannot load configuration")
		return nil
	}
	return res

}

func loadConfig() *Config {

	var res = new(Config)
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Config : cannot load config file", err.Error())
		return nil
	}

	if value, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal("Config : invalid server port", err.Error())
			return nil
		}
		res.ServerPort = port
	}

	if value, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal("Config : invalid db port", err.Error())
			return nil
		}
		res.DBPort = port

	}
	if value, found := os.LookupEnv("DBHOST"); found {
		res.DBHost = value
	}
	if value, found := os.LookupEnv("DBUSER"); found {
		res.DBUser = value
	}
	if value, found := os.LookupEnv("DBPASS"); found {
		res.DBPass = value
	}
	if value, found := os.LookupEnv("DBNAME"); found {
		res.DBName = value
	}
	if val, found := os.LookupEnv("SECRET"); found {
		res.Secret = val
	}
	return res
}
