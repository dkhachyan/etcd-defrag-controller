package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConnOpts struct {
	Endpoints   string
	CAfile      string
	Certfile    string
	Keyfile     string
	DialTimeout time.Duration
}

func NewEtcdClient(c ConnOpts) (*clientv3.Client, clientv3.Maintenance, error) {
	cfg := clientv3.Config{
		DialTimeout: c.DialTimeout,
		Endpoints:   endpoinsToList(c.Endpoints),
	}
	if c.CAfile != "" && c.Certfile != "" && c.Keyfile != "" {
		tlsConfig, err := NewTLSConfig(c)
		if err != nil {
			return nil, nil, err
		}
		cfg.TLS = tlsConfig
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	mcli := clientv3.NewMaintenance(cli)

	return cli, mcli, nil
}

func NewTLSConfig(c ConnOpts) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.Certfile, c.Keyfile)
	if err != nil {
		return nil, err
	}

	cacert, err := ioutil.ReadFile(c.Certfile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cacert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	return tlsConfig, nil
}

func endpoinsToList(endpoints string) []string {
	return strings.Split(endpoints, ",")
}