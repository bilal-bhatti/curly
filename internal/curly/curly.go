package curly

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Thing struct {
	Scheme  string                 `yaml:"scheme" json:"scheme"`
	Host    string                 `yaml:"host" json:"host"`
	Method  string                 `yaml:"method" json:"method"`
	Path    string                 `yaml:"path" json:"path"`
	Headers map[string]string      `yaml:"headers" json:"headers"`
	Body    interface{}            `yaml:"body" json:"body"`
	Query   map[string]interface{} `yaml:"query" json:"query"`
	Form    map[string]interface{} `yaml:"form" json:"form"`
}

type Curly struct {
	client *http.Client
}

func NewCurly() *Curly {
	return &Curly{
		client: http.DefaultClient,
	}
}

func (c *Curly) Go(t Thing) {
	switch strings.ToUpper(t.Method) {
	case http.MethodGet:
		c.get(t)
	case http.MethodPost:
	case http.MethodPut:
	default:
		c.get(t)
	}
}

func (c Curly) get(t Thing) error {
	var scheme = "http"
	var uri string

	if t.Scheme != "" {
		scheme = t.Scheme
	}

	if t.Host != "" {
		uri = fmt.Sprintf("%s://%s", scheme, t.Host)
	}

	uri = uri + t.Path

	if t.Query != nil {
		values := url.Values{}
		for k, vv := range t.Query {
			switch vvt := vv.(type) {
			case []interface{}:
				for _, v := range vvt {
					values.Add(k, fmt.Sprintf("%v", v))
				}
			case interface{}:
				values.Add(k, fmt.Sprintf("%v", vvt))
			default:
				values.Add(k, fmt.Sprintf("%v", vvt))
			}
		}

		if strings.Contains(uri, "?") {
			uri = uri + "&" + values.Encode()
		} else {
			uri = uri + "?" + values.Encode()
		}
	}

	log.Println("<", uri)
	req, err := http.NewRequest(http.MethodGet, uri, http.NoBody)
	if err != nil {
		return err
	}

	for k, v := range t.Headers {
		log.Println("<", k, v)
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	log.Println("> status code", resp.StatusCode)
	for k := range resp.Header {
		log.Println(">", k, resp.Header.Get(k))
	}

	bites, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	count, err := os.Stdout.Write(bites)
	if err != nil {
		return err
	}

	log.Printf("received %d bytes", count)

	return nil
}
