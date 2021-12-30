/*
Copyright © 2021 Bilal Bhatti
*/

package curly

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var Version = "DEV"

type Curly struct {
	client *http.Client
}

func NewCurly() *Curly {
	return &Curly{
		client: http.DefaultClient,
	}
}

type dumper func(resp *http.Response) error

func (c *Curly) Go(t Thing) {
	c.do(t, dump)
}

func (c Curly) do(t Thing, dump dumper) error {
	req, err := t.Request()
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("curly v%s", Version))

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
