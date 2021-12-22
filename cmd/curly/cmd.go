package main

import (
	"context"
	"flag"
	"log"

	"github.com/bilal-bhatti/curly/internal/curly"
	"github.com/google/subcommands"
)

type requestCmd struct {
}

func (*requestCmd) Name() string { return "apply" }

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

	env, err := curly.Env()
	if err != nil {
		log.Fatalln(err)
	}

	curly.Tracef("env: %v", env)

	return subcommands.ExitSuccess
}
