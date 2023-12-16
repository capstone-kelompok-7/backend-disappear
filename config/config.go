package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   int
	DBPort       int
	DBHost       string
	DBUser       string
	DBPass       string
	DBName       string
	Secret       string
	CCName       string
	CCAPIKey     string
	CCAPISecret  string
	CCFolder     string
	MongoURL     string
	OpenAiApiKey string
	ClientKey    string
	ServerKey    string
	Redis        Redis
	ResiKey      string
	FirebaseKey  string
}

type Redis struct {
	Addr string
	Pass string
}

func InitConfig() *Config {
	return loadConfig()

}

func loadConfig() *Config {

	var res = new(Config)
	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to fetch .env file")
		}
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
	if value, found := os.LookupEnv("CCNAME"); found {
		res.CCName = value
	}

	if value, found := os.LookupEnv("CCAPIKEY"); found {
		res.CCAPIKey = value
	}
	if value, found := os.LookupEnv("CCAPISECRET"); found {
		res.CCAPISecret = value
	}
	if value, found := os.LookupEnv("CCFOLDER"); found {
		res.CCFolder = value
	}
	if value, found := os.LookupEnv("MONGOURL"); found {
		res.MongoURL = value
	}
	if value, found := os.LookupEnv("OPENAIAPIKEY"); found {
		res.OpenAiApiKey = value
	}
	if value, found := os.LookupEnv("CLIENTKEY"); found {
		res.ClientKey = value
	}
	if value, found := os.LookupEnv("SERVERKEY"); found {
		res.ServerKey = value
	}
	if value, found := os.LookupEnv("REDIS_ADDR"); found {
		res.Redis.Addr = value
	}
	if value, found := os.LookupEnv("REDIS_PASS"); found {
		res.Redis.Pass = value
	}
	if value, found := os.LookupEnv("RESIKEY"); found {
		res.ResiKey = value
	}
	if value, found := os.LookupEnv("FIREBASEKEY"); found {
		res.FirebaseKey = value
	}

	return res
}
