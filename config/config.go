package config

import (
	"os"

	"github.com/iotames/miniutils"
	"github.com/joho/godotenv"
)

const ENV_PROD = "prod"
const ENV_DEV = "dev"
const ENV_FILE = ".env"
const DEFAULT_ENV_FILE = "env.default"

const DRIVER_SQLITE3 = "sqlite3"
const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"
const SQLITE_FILENAME = "sqlite3.db"

const DEFAULT_ENV_FILE_CONTENT = `# 此文件由系统自动创建，配置项为默认值。可修改本目录下的 .env 文件，以更新默认值。
# DB_DRIVER support: mysql,sqlite3,postgres
DB_DRIVER = "sqlite3"
DB_HOST = "127.0.0.1"
# DB_PORT like: 3306(mysql); 5432(postgres)
DB_PORT = 3306
DB_NAME = "lemocoder"
# DB_USERNAME like: root, postgres
DB_USERNAME = "root"
DB_PASSWORD = "root"
DB_NODE_ID = 1

# Server
SERVER_IP = "0.0.0.0"
SERVER_PORT = 8888
# 一段时间不活跃，自动断开连接
SERVER_DROP_AFTER = 300
# json, protobuf
MSG_FORMAT = "protobuf"
`

func LoadEnv() {
	initEnvFile()
	err := godotenv.Load(ENV_FILE, DEFAULT_ENV_FILE)
	if err != nil {
		panic("godotenv Error: " + err.Error())
	}
}

func initEnvFile() {
	if !miniutils.IsPathExists(ENV_FILE) {
		f, err := os.Create(ENV_FILE)
		if err != nil {
			panic("Create .env Error: " + err.Error())
		}
		f.Close()
	}
	if !miniutils.IsPathExists(DEFAULT_ENV_FILE) {
		f, err := os.Create(DEFAULT_ENV_FILE)
		if err != nil {
			panic("Create .env Error: " + err.Error())
		}
		f.WriteString(DEFAULT_ENV_FILE_CONTENT)
		f.Close()
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
