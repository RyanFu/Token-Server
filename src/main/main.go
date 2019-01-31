package main

import (
	"flag"
	"server"
)

var (
	port     = flag.Int("port", 3000, "Port to Serve")
	appid    = flag.String("appid", "", "WeChat APP ID")
	secret   = flag.String("secret", "", "WeChat Secret")
	host     = flag.String("host", "localhost:6379", "Redis Host")
	password = flag.String("password", "", "Redis Password")
	db       = flag.Int("db", 0, "Redis Database")
)

func main() {
	flag.Parse()
	s := server.Server{
		Port:     *port,
		Appid:    *appid,
		Secret:   *secret,
		Host:     *host,
		Password: *password,
		Database: *db,
	}
	s.Start()
}
