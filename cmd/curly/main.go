package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/nojima/httpie-go"
	"github.com/nojima/httpie-go/exchange"
	"github.com/nojima/httpie-go/flags"
	"github.com/nojima/httpie-go/input"
	"github.com/nojima/httpie-go/output"
)

/*
Usage: curly [-bdFfhjv] [-a value] [--check-status] [--http1] [--ignore-stdin] [--license] [-o value] [--overwrite] [--pretty value] [-p value] [--timeout value] [--verify value] [--version] [METHOD] URL [ITEM [ITEM ...]]
 -a, --auth=value     colon-separated username and password for authentication
 -b, --body           print only response body. shourtcut for --print=b
     --check-status   Also check the HTTP status code and exit with an error if
                      the status indicates one
 -d, --download       download file
 -F, --follow         follow 30x Location redirects
 -f, --form           data items are serialized as form fields
 -h, --headers        print only the request headers. shortcut for --print=h
     --http1          force HTTP/1.1 protocol
     --ignore-stdin   do not attempt to read stdin
 -j, --json           data items are serialized as JSON (default)
     --license        print license information and exit
 -o, --output=value   output file
     --overwrite      overwrite existing file
     --pretty=value   controls output formatting (all, format, none)
 -p, --print=value    specifies what the output should contain (HBhb)
     --timeout=value  timeout seconds that you allow the whole operation to take
 -v, --verbose        print the request as well as the response. shortcut for
                      --print=HBhb
     --verify=value   verify Host SSL certificate, 'yes' or 'no' ('yes' by
                      default, uppercase is also working)
     --version        print version and exit
*/

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func options() {
	options := &flags.OptionSet{
		InputOptions:    input.Options{},
		OutputOptions:   output.Options{},
		ExchangeOptions: exchange.Options{},
	}
	options.InputOptions.JSON = true
	fmt.Println(options)
}

func Main() error {
	options := &flags.OptionSet{
		InputOptions: input.Options{
			JSON:      true,
			Form:      false,
			ReadStdin: false,
		},
		OutputOptions: output.Options{
			PrintRequestHeader:  true,
			PrintRequestBody:    true,
			PrintResponseHeader: true,
			PrintResponseBody:   true,
			EnableFormat:        true,
			EnableColor:         true,
			Download:            false,
			OutputFile:          "",
			Overwrite:           false,
		},
		ExchangeOptions: exchange.Options{
			Timeout:         time.Duration(1 * time.Minute),
			FollowRedirects: false,
			Auth: exchange.AuthOptions{
				Enabled: false,
			},
			SkipVerify:  false,
			ForceHTTP1:  false,
			CheckStatus: true,
		},
	}

	url, _ := url.Parse("https://httpbin.org/anything")
	in := &input.Input{
		Method: input.Method("PUT"),
		URL:    url,
		Body: input.Body{
			BodyType: input.RawBody,
			Raw:      []byte(`{"foo":"bar"}`),
		},
	}

	// Send request and receive response
	status, err := httpie.Exchange(in, &options.ExchangeOptions, &options.OutputOptions)
	if err != nil {
		return err
	}

	if options.ExchangeOptions.CheckStatus {
		os.Exit(getExitStatus(status))
	}

	return nil
}

func getExitStatus(statusCode int) int {
	if 300 <= statusCode && statusCode < 600 {
		return statusCode / 100
	}
	return 0
}
