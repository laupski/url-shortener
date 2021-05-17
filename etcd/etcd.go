package etcd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Config struct {
	Address        []string `toml:"address"`
	DialTimeout    string   `toml:"dial_timeout"`
	RequestTimeout string   `toml:"request_timeout"`
}

type Connection struct {
	Context context.Context
	Client  clientv3.Client
}

func NewEtcdClient(c Config) (*Connection, error) {
	rt, err := time.ParseDuration(c.RequestTimeout)
	if err != nil {
		return nil, err
	}
	dt, err := time.ParseDuration(c.DialTimeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), rt)
	defer cancel()
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dt,
		Endpoints:   c.Address,
	})
	return &Connection{
		Context: ctx,
		Client:  *cli,
	}, err
}

func GetRedirect(c Connection, k string) (string, error) {
	fmt.Printf("Getting redirect link for %v\n", k)
	gr, err := c.Client.Get(c.Context, k)
	if err != nil {
		logrus.Fatal(err)
	}

	if len(gr.Kvs) < 1 {
		return "", err
	}

	return string(gr.Kvs[0].Value), err
}

func PutRedirect(c Connection, k string, v string) error {
	fmt.Printf("Putting redirect link for %v\n", v)
	_, err := c.Client.Put(c.Context, k, v)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
