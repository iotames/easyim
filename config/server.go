package config

import (
	"os"
	"strconv"
)

type Server struct {
	IP              string
	Port, DropAfter int
}

func GetServer() Server {
	ip := os.Getenv("SERVER_IP")
	portStr := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Error: Fail To Get SERVER_PORT," + err.Error())
	}
	dropAfterStr := os.Getenv("SERVER_DROP_AFTER")
	dropAfter, err := strconv.Atoi(dropAfterStr)
	if err != nil {
		panic("Error: Fail To Get SERVER_DROP_AFTER," + err.Error())
	}
	return Server{IP: ip, Port: port, DropAfter: dropAfter}
}
