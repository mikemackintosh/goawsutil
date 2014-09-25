package goawsutil

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// S3Client is a deminius S3 Client that follows the REST API
type S3Client struct {
	HTTPClient *http.Client
	Signer     *AWSV4Signer
}

// NewS3Client is a S3Client constructor
func NewS3Client(cred *Credentials, client *http.Client) *S3Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &S3Client{
		HTTPClient: client,
		Signer:     NewAWSV4Signer(cred),
	}
}

// Get does a HTTP GET and signs the request
func (c *S3Client) Get(urlStr string, xheaders map[string]string) (*http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	regionName := S3RegionFromURL(u)
	if regionName == "" {
		return nil, fmt.Errorf("URL does not appear to an S3 endpoint: %s", urlStr)

	}

	headers := make(http.Header)
	if xheaders != nil {
		for k, v := range xheaders {
			headers.Set(k, v)
		}
	}
	req := http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       nil,
		Header:     headers,
		Host:       u.Host,
	}

	c.Signer.Sign(&req, []byte{}, "", "s3", time.Now().UTC())

	return c.HTTPClient.Do(&req)
}

// Put does an HTTP PUT request to an S3 Asset
func (c *S3Client) Put(urlStr string, xheaders map[string]string, body []byte) (*http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	regionName := S3RegionFromURL(u)
	if regionName == "" {
		return nil, fmt.Errorf("URL does not appear to an S3 endpoint: %s", urlStr)
	}

	headers := make(http.Header)
	// Add our stuff

	// override
	if xheaders != nil {
		for k, v := range xheaders {
			headers.Set(k, v)
		}
	}

	req := http.Request{
		Method:     "PUT",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Host:       u.Host,
	}

	c.Signer.Sign(&req, body, regionName, "s3", time.Now().UTC())

	return c.HTTPClient.Do(&req)
}