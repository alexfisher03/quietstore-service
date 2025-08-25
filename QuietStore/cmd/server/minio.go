// cmd/server/minio.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newMinIOS3Client(endpoint, accessKey, secretKey string, useSSL bool) *s3.Client {
	if !strings.HasPrefix(endpoint, "http") {
		if useSSL {
			endpoint = "https://" + endpoint
		} else {
			endpoint = "http://" + endpoint
		}
	}

	u, _ := url.Parse(endpoint)

	var roots *x509.CertPool
	if pem, err := ioutil.ReadFile("/etc/ssl/certs/minio-ca.pem"); err == nil {
		roots = x509.NewCertPool()
		_ = roots.AppendCertsFromPEM(pem)
	}

	tlsCfg := &tls.Config{
		RootCAs:    roots,
		MinVersion: tls.VersionTLS12,
	}

	if useSSL {
		tlsCfg.ServerName = "localhost"
	}

	transport := &http.Transport{
		TLSClientConfig: tlsCfg,
		TLSNextProto:    map[string]func(string, *tls.Conn) http.RoundTripper{},
	}
	httpClient := &http.Client{Transport: transport}

	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      "us-east-1",
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: u.String(), HostnameImmutable: true}, nil
			},
		),
		HTTPClient: httpClient,
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
}
