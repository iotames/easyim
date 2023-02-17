package config

import (
	"os"
	"strconv"
)

type Server struct {
	Port int
}

func GetServer() Server {
	portStr := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Error: Fail To Get WEB_SERVER_PORT," + err.Error())
	}
	return Server{Port: port}
}
