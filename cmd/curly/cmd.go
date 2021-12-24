package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bilal-bhatti/curly/internal/curly"
	"github.com/google/subcommands"
	"gopkg.in/yaml.v3"
)

type requestCmd struct {
}

func (*requestCmd) Name() string { return "run" }

func (*requestCmd) Synopsis() string {
	return "execute HTTP request file"
}

func (*requestCmd) Usage() string {
	return `
execute request

	examples: 
		curly get.httpbin.yml 

`
}

func (a *requestCmd) SetFlags(f *flag.FlagSet) {}

func (a *requestCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// if _, err := os.Stat(a.template); errors.Is(err, os.ErrNotExist) {
	// 	log.Println(a.Usage())
	// 	log.Fatalln(err)
	// }

	var rfs []string
	for _, a := range f.Args() {
		if strings.HasSuffix(a, ".yml") {
			rfs = append(rfs, a)
		}
	}

	c := curly.NewCurly()

	for _, rf := range rfs {
		log.Println("* running", rf)

		cwd := path.Dir(rf)
		cwd, err := filepath.Abs(cwd)

		if err != nil {
			log.Fatalln(err)
		}

		env, err := curly.Env(cwd)
		if err != nil {
			log.Fatalln(err)
		}

		bites, err := ioutil.ReadFile(rf)
		if err != nil {
			log.Fatalln(err)
		}

		var raw interface{}
		var t curly.Thing

		err = yaml.Unmarshal(bites, &raw)
		if err != nil {
			log.Fatalln(err)
		}

		raw = curly.MapI2MapS(raw)

		err = curly.Merge(env.Data, raw)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.NewEncoder(os.Stdout).Encode(env.Data)
		if err != nil {
			log.Fatalln(err)
		}

		bites, err = json.Marshal(env.Data)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(bites, &t)
		if err != nil {
			log.Fatalln(err)
		}

		c.Go(t)
	}

	return subcommands.ExitSuccess
}
