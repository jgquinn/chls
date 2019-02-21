package chls // import "jgquinn.com/chls"

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const Version = "0.1.1-dev"

// ExecuteRequest replays connections on the provided http.Client,
// optionally overriding recorded headers with any provided NVPair.
func ExecuteRequest(client *http.Client, c *Connection, override ...*NVPair) (rep *http.Response, err error) {
	req := c.BuildRequest()
	for _, op := range override {
		req.Header.Set(op.Name, op.Value)
	}

	rep, err = client.Do(req)
	return
}

type NVPair struct {
	Name  string `json:'name'`
	Value string `json:'value'`
}

type Header struct {
	FirstLine string   `json:'firstLine'`
	Headers   []NVPair `json:'headers'`
}

func (h *Header) Map() map[string][]string {
	hm := make(map[string][]string)
	for _, hp := range h.Headers {
		hvp, found := hm[hp.Name]
		if !found {
			hvp = []string{hp.Value}
			hm[hp.Name] = hvp
		} else {
			hm[hp.Name] = append(hm[hp.Name], hp.Value)
		}
	}
	return hm
}

type Body struct {
	Charset string `json:'charset'`
	Text    string
}

type RepBody struct {
	Body
	Decoded bool `json:'decoded'`
}

type Side struct {
	Charset  *string `json:'charset,omitempty'`
	Encoding *string `json:'contentEncoding,omitempty'`
	Header   Header  `json:'header'`
	MIMEType string  `json:'mimeType'`
}

type Request struct {
	Side
	Body *Body `json:'body,omitempty'`
}

type Response struct {
	Side
	Status int      `json:'status'`
	Body   *RepBody `json:'body,omitempty'`
}

type Connection struct {
	LocalURL   *url.URL `json:'-''`
	Scheme     string   `json:'scheme'`
	ActualPort int      `json:'actualPort'`
	Port       int      `json:'port'`
	Host       string   `json:'host'`
	Method     string   `json:'method'`
	Path       string   `json:'path'`
	Query      *string  `json:'query'`
	Request    Request  `json:'request'`
	Response   Response `json:'response'`
}

func (c *Connection) URL() (u *url.URL) {

	if c.LocalURL != nil {
		u = c.LocalURL
		u.Path = c.Path
		return
	}

	u = &url.URL{
		Scheme: c.Scheme,
		Host:   fmt.Sprintf("%s:%d", c.Host, c.ActualPort),
		Path:   c.Path,
	}

	if c.Query != nil {
		u.RawQuery = *c.Query
	}

	return
}

func (c *Connection) BuildRequest() (req *http.Request) {
	req = &http.Request{
		Method: c.Method,
		URL:    c.URL(),
		Header: c.Request.Header.Map(),
	}

	if c.Request.Body != nil {
		req.Body = ioutil.NopCloser(strings.NewReader(c.Request.Body.Text))
	}

	return req
}

// Session represents the outer structure of a Charles JSON session file.
type Session []Connection
