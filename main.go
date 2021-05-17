package main

import (
	"github.com/BurntSushi/toml"
	"github.com/laupski/url-shortener/api"
	"github.com/laupski/url-shortener/etcd"
	"github.com/laupski/url-shortener/logs"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Title string `toml:"title"`
	Api   api.Config
	Logs  logs.Config
	Etcd  etcd.Config
}

func main() {
	var c Config
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		logrus.Fatal("Read configuration file error ", err.Error())
	}

	// Set the log configuration from config.toml
	logs.SetLog(c.Logs)

	// Set the etcd configuration from config.toml and open a connection
	conn, err := etcd.NewEtcdClient(c.Etcd)
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Client.Close()

	api.RunApi(c.Api, *conn)
}
