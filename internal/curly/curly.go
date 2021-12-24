package curly

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"moul.io/http2curl/v2"
)

type Curly struct {
	client *http.Client
	curl   bool
}

func NewCurly(curl bool) *Curly {
	return &Curly{
		client: http.DefaultClient,
		curl:   curl,
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

	if c.curl {
		curl, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return err
		}
		log.Println(curl.String())
		return nil
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

	if c.curl {
		curl, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return err
		}
		log.Println(curl.String())
		return nil
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

	if c.curl {
		curl, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return err
		}
		log.Println(curl.String())
		return nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	return dump(resp)
}

func dump(resp *http.Response) error {
	defer resp.Body.Close()

	log.Println(">", http.StatusText(resp.StatusCode), resp.StatusCode)
	for k := range resp.Header {
		log.Println(">H", k, resp.Header.Get(k))
	}

	if resp.StatusCode != http.StatusNoContent {
		bites, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		count, err := os.Stdout.Write(bites)
		if err != nil {
			return err
		}

		log.Printf("* received %d bytes", count)
	}

	return nil
}
