package curly

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Thing struct {
	Method  string
	URI     string
	Headers map[string]string
	Body    []byte
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
	req, err := http.NewRequest(http.MethodGet, t.URI, http.NoBody)
	if err != nil {
		return err
	}

	for k, v := range t.Headers {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	log.Println("status code", resp.StatusCode)
	// body, _ := ioutil.ReadAll(resp.Body)
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
