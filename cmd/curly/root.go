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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"moul.io/http2curl/v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curly [flags] <request-file.yml>",
	Short: "Execute an http request from supplied <request-file.yml>",
	Long: `curly is a small wrapper around go http.request to make
working with rest apis easy, as persistent collections of request
files. It can also print out the equivalent cURL command.

examples:

curly <request-file.yml>
curly -c <request-file.yml>
eval "$(curly -c <request-file.yml>)"
`,
	Args:   cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {},
	Run:    run,
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.curly.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("verbose", "v", false, "run with verbose")
	rootCmd.Flags().BoolP("curl", "c", false, "print cURL command")
	rootCmd.Flags().StringP("env", "e", "", "environment settings file")
}

func run(cmd *cobra.Command, args []string) {
	curly.Debug, _ = cmd.Flags().GetBool("verbose")

	var rfs []string
	for _, a := range args {
		if strings.HasSuffix(a, ".yml") {
			rfs = append(rfs, a)
		}
	}

	for _, req_file := range rfs {
		log.Println("* running", req_file)

		var env_path string
		var err error

		if env_f, _ := cmd.Flags().GetString("env"); env_f != "" {
			env_path = env_f
		} else {
			env_path = path.Dir(req_file)
		}

		env_path, err = filepath.Abs(env_path)
		if err != nil {
			log.Fatalln(err)
		}

		env, err := curly.Env(env_path)
		if err != nil {
			log.Fatalln(err)
		}

		bites, err := ioutil.ReadFile(req_file)
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

		err = yaml.NewEncoder(log.Writer()).Encode(env.Data)
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

		if curl, _ := cmd.Flags().GetBool("curl"); curl {
			req, err := t.Request()
			if err != nil {
				log.Fatalln(err)
			}

			curl, err := http2curl.GetCurlCommand(req)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("\n*** cURL command")
			fmt.Println(curl.String())
		} else {
			curly.NewCurly().Go(t)
		}
	}
}
