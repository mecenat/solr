package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type connection interface {
	request(ctx context.Context, method, path, contentType string, body []byte) (*Response, error)
	formatBasePath() string
	setBasicAuth(username, password string)
}

// Connection represents the connection to the solr server and
// includes information about the address of the server and
// and the client to be used for connecting to it.
type Connection struct {
	httpClient *http.Client
	Host       string
	Core       string
	Username   string
	Password   string
}

// NewConnection ...
func NewConnection(host, core string, client *http.Client) (*Connection, error) {
	if host == "" || core == "" {
		return nil, ErrInvalidConfig
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	return &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}, nil
}

func (c *Connection) formatBasePath() string {
	return formatBasePath(c.Host, c.Core)
}

func (c *Connection) setBasicAuth(username, password string) {
	c.Username = username
	c.Password = password
}

func (c *Connection) request(ctx context.Context, method, url, contentType string, body []byte) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	res, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r Response
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Error != nil {
		return &r, r.Error
	}

	return &r, nil
}

// RetryableConnection implements the retryablehttp library from Hashicorp that allows
// making a http request multiple times with a set time in case of failure due to
// connectivity issues. This for example can be useful if your solr servers are
// being shutdown while a new one gets started, the request can continue
// trying allowing for the server to be replaced without dropping it.
type RetryableConnection struct {
	Host        string
	Core        string
	Username    string
	Password    string
	Timeout     time.Duration
	httpClient  *http.Client
	retryClient *retryablehttp.Client
}

type RetryableConfig struct {
	Timeout      time.Duration
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	RetryMax     int
	NoLog        bool
}

// NewRetryableConnection ...
func NewRetryableConnection(host, core string, client *http.Client, conf *RetryableConfig) (*RetryableConnection, error) {
	if host == "" || core == "" {
		return nil, ErrInvalidConfig
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		conf = &RetryableConfig{
			Timeout:      10 * time.Second,
			RetryWaitMin: 50 * time.Millisecond,
			RetryWaitMax: 2 * time.Second,
			RetryMax:     4,
			NoLog:        true,
		}
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = client
	retryClient.HTTPClient.Timeout = conf.Timeout
	retryClient.RetryWaitMin = conf.RetryWaitMin
	retryClient.RetryWaitMax = conf.RetryWaitMax
	retryClient.RetryMax = conf.RetryMax
	if conf.NoLog {
		retryClient.Logger = log.New(io.Discard, "", log.LstdFlags)
	}

	return &RetryableConnection{
		Host:        host,
		Core:        core,
		Timeout:     conf.Timeout,
		httpClient:  client,
		retryClient: retryClient,
	}, nil
}

// NewDefaultHTTPClient returns an *http.Client configured with sensible defaults
// for connection pooling and keep-alive. Use this instead of http.DefaultClient
// which only keeps 2 idle connections per host, causing excessive reconnections
// under concurrent load.
func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			ForceAttemptHTTP2:   true,
		},
	}
}

func (c *RetryableConnection) formatBasePath() string {
	return formatBasePath(c.Host, c.Core)
}

func (c *RetryableConnection) setBasicAuth(username, password string) {
	c.Username = username
	c.Password = password
}

func (c *RetryableConnection) request(ctx context.Context, method, path, contentType string, body []byte) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := retryablehttp.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	res, err := c.retryClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r Response
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Error != nil {
		return &r, r.Error
	}

	return &r, nil
}
