/*
Copyright © 2021 Bilal Bhatti
*/

package curly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Thing struct {
	Scheme  string                 `yaml:"scheme" json:"scheme"`
	Host    string                 `yaml:"host" json:"host"`
	Method  string                 `yaml:"method" json:"method"`
	Prefix  string                 `yaml:"prefix" json:"prefix"`
	Path    string                 `yaml:"path" json:"path"`
	Headers map[string]string      `yaml:"headers" json:"headers"`
	Body    interface{}            `yaml:"body" json:"body"`
	Query   map[string]interface{} `yaml:"query" json:"query"`
	Form    map[string]interface{} `yaml:"form" json:"form"`
}

func (t *Thing) URL() (*url.URL, error) {
	var uri string

	t.Method = strings.ToUpper(t.Method)
	if t.Method == "" {
		t.Method = http.MethodGet
	}

	if strings.HasPrefix(t.Path, "http") {
		// fully qualified url/path provided
		// ignore all else
		uri = t.Path
	} else {
		var scheme = "http"

		if t.Scheme != "" {
			scheme = t.Scheme
		}

		t.Host = strings.TrimSpace(t.Host)
		t.Path = strings.TrimSpace(t.Path)
		t.Prefix = strings.TrimSpace(t.Prefix)

		if t.Host != "" {
			uri = fmt.Sprintf("%s://%s%s%s", scheme, t.Host, t.Prefix, t.Path)
		}
	}

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

	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (t Thing) Request() (*http.Request, error) {
	endpoint, err := t.URL()
	if err != nil {
		return nil, err
	}

	t.Headers["User-Agent"] = "curly v0.0.1"

	log.Println("*", t.Method, endpoint.String())

	var req *http.Request
	var body io.Reader

	if t.Body != nil {
		log.Println("* setting body")
		buf := &bytes.Buffer{}

		if err := json.NewEncoder(buf).Encode(t.Body); err != nil {
			log.Println(err)
		}

		body = buf
	} else {
		body = http.NoBody
	}

	if t.Form != nil {
		log.Println("* setting form data")
		values := url.Values{}

		for k, vv := range t.Form {
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

		body = strings.NewReader(values.Encode())
	}

	req, err = http.NewRequest(t.Method, endpoint.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range t.Headers {
		log.Println("<H", k, v)
		req.Header.Add(k, v)
	}

	return req, nil
}
