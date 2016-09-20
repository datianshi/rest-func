package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Rest struct {
	URL string
}

type connect func(*http.Request, http.RoundTripper) (*http.Response, error)

type ConnectParams struct {
	Request   *http.Request
	Transport http.RoundTripper
	con       connect
}

type httpMethod func(*http.Request) *http.Request

var Get httpMethod = func(request *http.Request) *http.Request {
	request.Method = "GET"
	return request
}

var POST httpMethod = func(request *http.Request) *http.Request {
	request.Method = "POST"
	return request
}

var PUT httpMethod = func(request *http.Request) *http.Request {
	request.Method = "PUT"
	return request
}

var DELETE httpMethod = func(request *http.Request) *http.Request {
	request.Method = "DELETE"
	return request
}

func (r *Rest) Build() *ConnectParams {
	request, _ := http.NewRequest("GET", r.URL, nil)
	return &ConnectParams{
		Request:   request,
		Transport: http.DefaultTransport,
	}
}

func (c *ConnectParams) Connect() (*http.Response, error) {
	debug := "true" == os.Getenv("HTTP_DEBUG")
	if debug {
		dump, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			return nil, errors.New("Dump http is enabled, failed to dump the message")
		}
		fmt.Fprintf(os.Stdout, "Request:\n%s\n\n", dump)
	}
	client := &http.Client{
		Transport: c.Transport,
	}
	response, err := client.Do(c.Request)
	if err != nil {
		return nil, err
	}
	if debug {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, errors.New("Dump http is enabled, failed to dump the message")
		}
		fmt.Fprintf(os.Stdout, "Response:\n%s", dump)
	}
	return response, nil
}

func (c *ConnectParams) WithHttpMethod(method httpMethod) *ConnectParams {
	c.Request = method(c.Request)
	return c
}

func (c *ConnectParams) WithHttpHeader(key string, value string) *ConnectParams {
	c.Request.Header.Add(key, value)
	return c
}

func (c *ConnectParams) WithContentType(value string) *ConnectParams {
	return c.WithHttpHeader("Content-Type", value)
}

func (c *ConnectParams) WithBasicAuth(user string, password string) *ConnectParams {
	auth := fmt.Sprintf("%s:%s", user, password)
	encode := base64.StdEncoding.EncodeToString([]byte(auth))
	return c.WithHttpHeader("Authorization", fmt.Sprintf("Basic %s", encode))
}

func (c *ConnectParams) WithMultipartForm(paramName string, file *os.File) *ConnectParams {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(paramName, file.Name())
	io.Copy(part, file)
	writer.Close()
	c.Request.Body = ioutil.NopCloser(body)
	return c.WithContentType(writer.FormDataContentType()).WithHttpMethod(POST)
}

func (c *ConnectParams) WithHttpBody(body io.ReadCloser) *ConnectParams {
	c.Request.Body = body
	return c
}

func (c *ConnectParams) WithFormValue(values url.Values) *ConnectParams {
	c.Request.Body = ioutil.NopCloser(strings.NewReader(values.Encode()))
	c.WithContentType("application/x-www-form-urlencoded")
	return c.WithHttpMethod(POST)
}

func (c *ConnectParams) SkipSslVerify(skip bool) *ConnectParams {
	if skip {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		c.Transport = transport
	}
	return c
}
