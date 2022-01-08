/*
Copyright © 2021 Bilal Bhatti
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bilal-bhatti/curly/internal/curly"
	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"moul.io/http2curl/v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curly [flags] <request-file.yml>",
	Short: "Execute an http request from supplied <request-file.yml>",
	Long: `curly is a small wrapper around go http.request to make
working with rest apis easier, as persistent collections of request
files. It can also print out the equivalent cURL command.

examples:

curly <request-file.yml>
curly -c <request-file.yml> -e <env.yml>
eval "$(curly -c <request-file.yml>)"
`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if curly.Verbose {
			log.Printf("curly v%s\n", curly.Version)
		}
	},
	RunE: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&curly.Verbose, "verbose", "v", false, "run with verbose")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("curl", "c", false, "print cURL command")
	rootCmd.Flags().StringP("env", "e", "", "environment settings file")
	// rootCmd.Flags().Bool("$", false, "cache http response")
}

func run(cmd *cobra.Command, args []string) error {
	logger := os.Stderr
	// ouput := os.Stdout

	for _, req_arg := range args {
		if curly.Verbose {
			fmt.Fprintln(logger, "* executing", req_arg)
		}

		var load loader

		if strings.HasSuffix(req_arg, ".yml") || strings.HasSuffix(req_arg, ".yaml") {
			load = from_file{req_arg: req_arg}
		} else if strings.HasPrefix(req_arg, "http") {
			load = from_path{path: req_arg}
		} else {
			// don't know what that is, skip
			continue
		}

		t, err := prepare(cmd, load)
		if err != nil {
			return nil
		}

		if err := doThing(cmd, t); err != nil {
			return err
		}
	}

	return nil
}

func prepare(cmd *cobra.Command, load loader) (*curly.Thing, error) {
	var env_path string
	var err error

	if env_f, _ := cmd.Flags().GetString("env"); env_f != "" {
		env_path = env_f
	} else {
		env_path = load.env_dir()
	}

	env_path, err = filepath.Abs(env_path)
	if err != nil {
		return nil, err
	}

	env, err := curly.Env(env_path)
	if err != nil {
		return nil, err
	}

	raw, err := load.raw()
	if err != nil {
		return nil, err
	}

	err = curly.Merge(env.Data, raw)
	if err != nil {
		return nil, err
	}

	if curly.Verbose {
		if err = yaml.NewEncoder(log.Writer()).Encode(env.Data); err != nil {
			return nil, err
		}
	}

	jtool := jt.Template{Debug: curly.Verbose}
	err = jtool.Apply(env.Data, env.Data)
	if err != nil {
		return nil, err
	}

	bites, err := json.Marshal(env.Data)
	if err != nil {
		return nil, err
	}

	t := &curly.Thing{}
	t.Cwd = load.env_dir()

	err = json.Unmarshal(bites, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func doThing(cmd *cobra.Command, t *curly.Thing) error {
	if curl, _ := cmd.Flags().GetBool("curl"); curl {
		req, err := t.Request()
		if err != nil {
			return err
		}

		curl, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return err
		}

		log.Println("\n*** cURL command")
		fmt.Println(curl.String())
	} else {
		err := curly.NewCurly().Go(t)
		if err != nil {
			return err
		}
	}

	return nil
}

type loader interface {
	raw() (interface{}, error)
	env_dir() string
}

type from_file struct {
	req_arg string
}

func (f from_file) raw() (interface{}, error) {
	bites, err := ioutil.ReadFile(f.req_arg)
	if err != nil {
		return nil, err
	}

	var raw interface{}

	err = yaml.Unmarshal(bites, &raw)
	if err != nil {
		return nil, err
	}

	return curly.MapI2MapS(raw), nil
}

func (f from_file) env_dir() string {
	return path.Dir(f.req_arg)
}

type from_path struct {
	path string
}

func (f from_path) raw() (interface{}, error) {
	return map[string]string{"path": f.path}, nil
}

func (f from_path) env_dir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
