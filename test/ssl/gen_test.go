package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGen(t *testing.T) {
	certSetup()
}

func TestTLSClient(t *testing.T) {
	caCertPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		fmt.Println(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)
	ce, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	fmt.Println(err)
	tlsConf := &tls.Config{
		RootCAs: caCertPool,
		Certificates: []tls.Certificate{
			ce,
		},
	}
	tr := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: tr}
	res, err := client.Get("https://127.0.0.1:8033/tls")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
