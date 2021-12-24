package curly

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Curly struct {
	client *http.Client
}

func NewCurly() *Curly {
	return &Curly{
		client: http.DefaultClient,
	}
}

func (c *Curly) Go(t Thing) {
	t.Method = strings.ToUpper(t.Method)

	switch t.Method {
	case http.MethodGet:
		c.get(t, dump)
	case http.MethodPost:
		c.post(t, dump)
	case http.MethodPut:
		c.put(t, dump)
	default:
		c.get(t, dump)
	}
}

type dumper func(resp *http.Response) error

func (c Curly) get(t Thing, dump dumper) error {
	req, err := t.Request()
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return dump(resp)
}

func (c Curly) put(t Thing, dump dumper) error {
	req, err := t.Request()
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return dump(resp)
}

func (c Curly) post(t Thing, dump dumper) error {
	req, err := t.Request()
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	return dump(resp)
}

func dump(resp *http.Response) error {
	defer resp.Body.Close()

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

	log.Printf("* received %d bytes", count)

	return nil
}
