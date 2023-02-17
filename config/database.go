package config

import (
	"fmt"
	"os"
	"strconv"
)

type Database struct {
	Driver, Host, Username, Password, Name string
	Port, NodeID                           int
}

func GetDatabase() Database {
	dbDriver := os.Getenv("DB_DRIVER")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	nodeIdStr := os.Getenv("DB_NODE_ID")
	dbname := os.Getenv("DB_NAME")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Error: Fail To Get DB_PORT," + err.Error())
	}
	nodeID, err := strconv.Atoi(nodeIdStr)
	if err != nil {
		panic("Error: Fail To Get DB_NODE_ID," + err.Error())
	}
	return Database{Driver: dbDriver, Host: host, Username: username, Password: password, Name: dbname, Port: port, NodeID: nodeID}
}

func (d Database) GetAddr() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

func (d Database) GetDSN() string {
	dsnMap := map[string]string{
		DRIVER_MYSQL:    fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.Port, d.Name),
		DRIVER_POSTGRES: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.Username, d.Password, d.Name, d.Port),
	}
	dsn, ok := dsnMap[d.Driver]
	if !ok {
		dsnLen := len(dsnMap)
		ds := make([]string, dsnLen)
		for k := range dsnMap {
			ds = append(ds, k)
		}
		errMsg := fmt.Sprintf("ENV error: DB_DRIVER only Support: %v", ds)
		panic(errMsg)
	}
	return dsn
}
