package config

import (
	"os"

	"github.com/iotames/miniutils"
	"github.com/joho/godotenv"
)

const ENV_PROD = "prod"
const ENV_DEV = "dev"
const ENV_FILE = ".env"

const DRIVER_SQLITE3 = "sqlite3"
const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"
const SQLITE_FILENAME = "sqlite3.db"

func LoadEnv() {
	if !miniutils.IsPathExists(ENV_FILE) {
		f, err := os.Create(ENV_FILE)
		if err != nil {
			panic("Create .env Error: " + err.Error())
		}
		f.Close()
	}
	err := godotenv.Load(ENV_FILE, "env.default")
	if err != nil {
		panic("godotenv Error: " + err.Error())
	}
}

// var envconfig *EnvConfig
// var once sync.Once

// type EnvConfig struct {
// 	Database Database
// }

// func (e *EnvConfig) Load() {
// 	e.Database = *GetDatabase()
// }

// func GetEnvConfig() EnvConfig {
// 	once.Do(func() {
// 		fmt.Println("-----First---GetEnvConfig---once.Do")
// 		envconfig = &EnvConfig{}
// 		envconfig.Load()
// 	})
// 	return *envconfig
// }
