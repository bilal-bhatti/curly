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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

type Thing struct {
	Cwd     string            `yaml:"-" json:"-"`
	Scheme  string            `yaml:"scheme" json:"scheme"`
	Host    string            `yaml:"host" json:"host"`
	Method  string            `yaml:"method" json:"method"`
	Prefix  string            `yaml:"prefix" json:"prefix"`
	Path    string            `yaml:"path" json:"path"`
	Headers map[string]string `yaml:"headers" json:"headers"`
	Body    any               `yaml:"body" json:"body"`
	Query   map[string]any    `yaml:"query" json:"query"`
	Form    map[string]any    `yaml:"form" json:"form"`
}

var epf = regexp.MustCompile(`\$(@){(.+)}`)

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
		if t.Scheme == "" {
			t.Scheme = "http"
		}

		t.Host = strings.TrimSpace(t.Host)
		t.Path = strings.TrimSpace(t.Path)
		t.Prefix = strings.TrimSpace(t.Prefix)

		if t.Host != "" {
			uri = fmt.Sprintf("%s://%s%s%s", t.Scheme, t.Host, t.Prefix, t.Path)
		}
	}

	if t.Query != nil {
		values := t.values(t.Query)

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

type trimmer struct {
	reader io.Reader
}

func (t trimmer) Read(buf []byte) (int, error) {
	n, err := t.reader.Read(buf)
	if err != nil {
		return n, err
	}
	trimmed := bytes.TrimFunc(buf[:n], unicode.IsSpace)
	n = copy(buf, trimmed)
	return n, nil
}

func (t Thing) Request() (*http.Request, error) {
	endpoint, err := t.URL()
	if err != nil {
		return nil, err
	}

	log.Println("*", t.Method, endpoint.String())

	var body io.Reader

	if t.Body != nil {
		if bs, ok := t.Body.(string); ok {
			match := epf.FindStringSubmatch(bs)
			if len(match) > 0 {
				body = t.body_from_file(match)
			} else {
				body = strings.NewReader(bs)
			}
		} else {
			body = t.body_as_json()
		}
	} else if t.Form != nil {
		if Verbose {
			log.Println("* setting form data")
		}

		values := t.values(t.Form)

		body = strings.NewReader(values.Encode())
	} else {
		body = http.NoBody
	}

	req, err := http.NewRequest(t.Method, endpoint.String(), trimmer{reader: body})
	if err != nil {
		return nil, err
	}

	for k, v := range t.Headers {
		if !slices.Contains(redact, k) {
			log.Println("<H", k, v)
		}
		req.Header.Add(k, v)
	}

	return req, nil
}

func (t Thing) body_as_json() io.Reader {
	if Verbose {
		log.Println("* setting json body")
	}

	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(t.Body); err != nil {
		log.Println(err)
	}

	return buf
}

func (t Thing) body_from_file(match []string) io.Reader {
	fp, err := filepath.Abs(path.Join(t.Cwd, match[2]))

	if err != nil {
		log.Println(err)
	}

	if Verbose {
		log.Println("* setting body from file", fp)
	}

	f, err := os.Open(fp)
	if err != nil {
		log.Println(err)
	}

	return (*os.File)(f)
}

func (t Thing) values(data map[string]any) url.Values {
	values := url.Values{}

	for k, vv := range data {
		switch vvt := vv.(type) {
		case []any:
			for _, v := range vvt {
				values.Add(k, fmt.Sprintf("%v", v))
			}
		case any:
			values.Add(k, fmt.Sprintf("%v", vvt))
		default:
			values.Add(k, fmt.Sprintf("%v", vvt))
		}
	}

	return values
}
