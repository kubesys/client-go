/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server                   string
	ClientCertificateData    string
	ClientKeyData            string
	CertificateAuthorityData string
}

func NewForConfig(kubeConfig string) (*Config, error) {
	f, err := os.Open(kubeConfig)
	if err != nil {
		f.Close()
		return nil, err
	}
	config := new(Config)
	buf := bufio.NewReader(f)
	n := 0
	for {
		b, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			f.Close()
			return nil, errors.New("kubeconfig file error")
		}
		s := strings.Replace(string(b), " ", "", -1)
		if s == "" {
			break
		} else if strings.Contains(s, "server:") {
			config.Server = strings.Replace(strings.Split(s, "server:")[1], "\n", "", -1)
			n++
		} else if strings.Contains(s, "client-certificate-data:") {
			config.ClientCertificateData = strings.Replace(strings.Split(s, "client-certificate-data:")[1], "\n", "", -1)
			n++
		} else if strings.Contains(s, "client-key-data:") {
			config.ClientKeyData = strings.Replace(strings.Split(s, "client-key-data:")[1], "\n", "", -1)
			n++
		} else if strings.Contains(s, "certificate-authority-data:") {
			config.CertificateAuthorityData = strings.Replace(strings.Split(s, "certificate-authority-data:")[1], "\n", "", -1)
			n++
		}
	}
	f.Close()
	if n != 4 {
		return nil, errors.New("kubeconfig file error")
	}
	return config, nil
}

func HTTPClientFor(config *Config) (*http.Client, error) {
	if config == nil {
		return nil, errors.New("kubeconfig nil")
	}
	tlsConfig, err := TLSConfigFor(config)
	if err != nil {
		return nil, err
	}
	transport := http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
		MaxIdleConnsPerHost: 25,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	return &http.Client{Transport: &transport}, nil
}

func TLSConfigFor(config *Config) (*tls.Config, error) {
	certificateAuthorityData, err := base64.StdEncoding.DecodeString(config.CertificateAuthorityData)
	if err != nil {
		return nil, err
	}
	clientCertificateData, err := base64.StdEncoding.DecodeString(config.ClientCertificateData)
	if err != nil {
		return nil, err
	}
	clientKeyData, err := base64.StdEncoding.DecodeString(config.ClientKeyData)
	if err != nil {
		return nil, err
	}
	ca := rootCertPool(certificateAuthorityData)
	cert, err := tls.X509KeyPair(clientCertificateData, clientKeyData)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		RootCAs:            ca,
		Certificates:       []tls.Certificate{cert},
	}, nil
}

// rootCertPool returns nil if caData is empty.  When passed along, this will mean "use system CAs".
// When caData is not empty, it will be the ONLY information used in the CertPool.
func rootCertPool(caData []byte) *x509.CertPool {
	// What we really want is a copy of x509.systemRootsPool, but that isn't exposed.  It's difficult to build (see the go
	// code for a look at the platform specific insanity), so we'll use the fact that RootCAs == nil gives us the system values
	// It doesn't allow trusting either/or, but hopefully that won't be an issue
	if len(caData) == 0 {
		return nil
	}

	// if we have caData, use it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return certPool
}
